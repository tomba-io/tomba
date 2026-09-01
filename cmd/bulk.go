package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var (
	bulkFile        string
	bulkType        string
	bulkColumn      string
	bulkDomainCol   string
	bulkFirstCol    string
	bulkLastCol     string
	bulkUrlCol      string
	bulkConcurrency int
)

// bulkCmd represents the bulk command
var bulkCmd = &cobra.Command{
	Use:     "bulk",
	Aliases: []string{"b"},
	Short:   "Process a CSV file in bulk — auto-maps columns or use custom mapping.",

	Run:     bulkRun,
	Example: bulkExample,
}

func init() {
	bulkCmd.Flags().StringVar(&bulkFile, "file", "", "Input CSV file path (required).")
	bulkCmd.Flags().StringVar(&bulkType, "type", "enrich", "Operation type: enrich, verify, finder, search, author, linkedin, phone, sources.")
	bulkCmd.Flags().StringVar(&bulkColumn, "column", "", "Column name for email or domain (auto-detected if empty).")
	bulkCmd.Flags().StringVar(&bulkDomainCol, "domain-col", "", "Column name for domain (for finder type).")
	bulkCmd.Flags().StringVar(&bulkFirstCol, "first-col", "", "Column name for first name (for finder type).")
	bulkCmd.Flags().StringVar(&bulkLastCol, "last-col", "", "Column name for last name (for finder type).")
	bulkCmd.Flags().StringVar(&bulkUrlCol, "url-col", "", "Column name for URL (for author/linkedin type).")
	bulkCmd.Flags().IntVar(&bulkConcurrency, "concurrency", 0, "Number of concurrent workers (0=auto from plan, max 60).")
	_ = bulkCmd.MarkFlagRequired("file")
}

// getConcurrency determines the number of concurrent workers based on the plan name.
func getConcurrency(planName string) int {
	switch strings.ToLower(strings.TrimSpace(planName)) {
	case "free":
		return 1
	case "basic":
		return 3
	case "growth":
		return 5
	case "pro":
		return 8
	case "pay-as-you-go 20k", "payg 20k":
		return 10
	default:
		// Pro Plus, Enterprise, Scale, PAYG 50k+ = unlimited, cap at 60
		return 60
	}
}

type bulkJob struct {
	index int
	row   []string
}

type bulkResult struct {
	index  int
	result map[string]any
}

func bulkRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)

	// Show authenticated user and detect plan for concurrency
	var planName string
	account, err := init.Account()
	if err == nil {
		raw, _ := account.Marshal()
		var accData map[string]any
		if json.Unmarshal(raw, &accData) == nil {
			if d, ok := accData["data"].(map[string]any); ok {
				if email, ok := d["email"].(string); ok {
					fmt.Printf("%s Authenticated as %s\n", util.SuccessIcon(), util.Green(email))
				}
				if pricing, ok := d["pricing"].(map[string]any); ok {
					if name, ok := pricing["name"].(string); ok {
						planName = name
					}
				}
			}
		}
	}

	// Open CSV file
	file, err := os.Open(bulkFile)
	if err != nil {
		fmt.Printf("%s Cannot open file: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("%s Cannot parse CSV: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	if len(records) < 2 {
		fmt.Printf("%s CSV file must have a header row and at least one data row\n", util.ErrorIcon())
		return
	}

	headers := records[0]
	rows := records[1:]

	fmt.Printf("  Columns: %s\n", util.Cyan(strings.Join(headers, ", ")))
	fmt.Printf("  Rows: %s\n\n", util.Bold(fmt.Sprintf("%d", len(rows))))

	// Auto-map or use specified columns
	colMap := mapColumns(headers, bulkType)
	if colMap == nil {
		return
	}

	// Determine output file
	outputFile := init.Output
	if outputFile == "" {
		outputFile = strings.TrimSuffix(bulkFile, ".csv") + "_enriched.csv"
	}

	// Determine concurrency
	workers := bulkConcurrency
	if workers <= 0 {
		workers = getConcurrency(planName)
	}
	if workers > 60 {
		workers = 60
	}
	if workers > len(rows) {
		workers = len(rows)
	}

	fmt.Printf("  Plan: %s | Workers: %s | Rows: %s\n\n",
		util.Bold(planName), util.Cyan(fmt.Sprintf("%d", workers)), util.Bold(fmt.Sprintf("%d", len(rows))))

	startTime := time.Now()

	// Process rows concurrently
	progress := output.NewProgressBar(len(rows))
	stats := &output.BulkStats{
		Total:      len(rows),
		OutputFile: outputFile,
	}

	results := make([]map[string]any, len(rows))
	jobs := make(chan bulkJob, workers*2)
	resultsCh := make(chan bulkResult, workers*2)

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Per-worker rate limiter: spread requests across workers
			interval := time.Second / time.Duration(workers)
			if interval < 10*time.Millisecond {
				interval = 10 * time.Millisecond
			}
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for job := range jobs {
				<-ticker.C
				result := processBulkRow(init, job.row, headers, colMap, bulkType)
				resultsCh <- bulkResult{index: job.index, result: result}
			}
		}()
	}

	// Send jobs
	go func() {
		for i, row := range rows {
			jobs <- bulkJob{index: i, row: row}
		}
		close(jobs)
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results
	var mu sync.Mutex
	for r := range resultsCh {
		results[r.index] = r.result
		mu.Lock()
		if r.result != nil {
			if _, hasError := r.result["_error"]; hasError {
				stats.Errors++
			} else {
				stats.Found++
			}
		} else {
			stats.NotFound++
		}
		mu.Unlock()
		progress.Increment()
	}

	// Write output CSV
	if err := writeBulkOutput(outputFile, headers, rows, results, bulkType); err != nil {
		fmt.Printf("\n%s Error writing output: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	elapsed := time.Since(startTime)
	stats.Elapsed = elapsed
	stats.PrintStats()
}

type columnMapping struct {
	emailIdx  int
	domainIdx int
	firstIdx  int
	lastIdx   int
	urlIdx    int
}

func mapColumns(headers []string, opType string) *columnMapping {
	cm := &columnMapping{
		emailIdx:  -1,
		domainIdx: -1,
		firstIdx:  -1,
		lastIdx:   -1,
		urlIdx:    -1,
	}

	// Auto-detect columns
	for i, h := range headers {
		lower := strings.ToLower(strings.TrimSpace(h))
		switch {
		case bulkColumn != "" && strings.EqualFold(h, bulkColumn):
			cm.emailIdx = i
		case bulkDomainCol != "" && strings.EqualFold(h, bulkDomainCol):
			cm.domainIdx = i
		case bulkFirstCol != "" && strings.EqualFold(h, bulkFirstCol):
			cm.firstIdx = i
		case bulkLastCol != "" && strings.EqualFold(h, bulkLastCol):
			cm.lastIdx = i
		case bulkUrlCol != "" && strings.EqualFold(h, bulkUrlCol):
			cm.urlIdx = i
		case cm.emailIdx == -1 && (lower == "email" || lower == "e-mail" || lower == "email_address" || lower == "emailaddress" || lower == "mail"):
			cm.emailIdx = i
		case cm.domainIdx == -1 && (lower == "domain" || lower == "company_domain" || lower == "website" || lower == "company_website"):
			cm.domainIdx = i
		case cm.firstIdx == -1 && (lower == "first_name" || lower == "firstname" || lower == "first" || lower == "fname"):
			cm.firstIdx = i
		case cm.lastIdx == -1 && (lower == "last_name" || lower == "lastname" || lower == "last" || lower == "lname"):
			cm.lastIdx = i
		case cm.urlIdx == -1 && (lower == "url" || lower == "link" || lower == "linkedin" || lower == "linkedin_url" || lower == "article_url" || lower == "profile_url"):
			cm.urlIdx = i
		}
	}

	// Validate required columns for operation type
	switch opType {
	case "enrich", "verify":
		if cm.emailIdx == -1 {
			cm.emailIdx = promptColumnSelect(headers, "email")
			if cm.emailIdx == -1 {
				fmt.Printf("%s Could not find email column. Use --column to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
	case "finder":
		if cm.domainIdx == -1 {
			cm.domainIdx = promptColumnSelect(headers, "domain")
			if cm.domainIdx == -1 {
				fmt.Printf("%s Could not find domain column. Use --domain-col to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		if cm.firstIdx == -1 || cm.lastIdx == -1 {
			if cm.firstIdx == -1 {
				cm.firstIdx = promptColumnSelect(headers, "first name")
			}
			if cm.lastIdx == -1 {
				cm.lastIdx = promptColumnSelect(headers, "last name")
			}
		}
		fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
		if cm.firstIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for first name\n", util.SuccessIcon(), util.Bold(headers[cm.firstIdx]))
		}
		if cm.lastIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for last name\n", util.SuccessIcon(), util.Bold(headers[cm.lastIdx]))
		}
	case "search":
		if cm.domainIdx == -1 {
			cm.domainIdx = promptColumnSelect(headers, "domain")
			if cm.domainIdx == -1 {
				fmt.Printf("%s Could not find domain column. Use --domain-col to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
	case "author", "linkedin":
		if cm.urlIdx == -1 {
			cm.urlIdx = promptColumnSelect(headers, "URL")
			if cm.urlIdx == -1 {
				fmt.Printf("%s Could not find URL column. Use --url-col to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for URL\n", util.SuccessIcon(), util.Bold(headers[cm.urlIdx]))
	case "phone":
		if cm.emailIdx == -1 {
			cm.emailIdx = promptColumnSelect(headers, "email")
			if cm.emailIdx == -1 {
				fmt.Printf("%s Could not find email column. Use --column to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
	case "sources":
		if cm.emailIdx == -1 {
			cm.emailIdx = promptColumnSelect(headers, "email")
			if cm.emailIdx == -1 {
				fmt.Printf("%s Could not find email column. Use --column to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
	}

	fmt.Println()
	return cm
}

func promptColumnSelect(headers []string, fieldName string) int {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select the column for %s", fieldName),
		Items: headers,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return -1
	}
	return idx
}

func processBulkRow(conn *start.Conn, row []string, headers []string, cm *columnMapping, opType string) map[string]interface{} {
	switch opType {
	case "enrich":
		if cm.emailIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		email := strings.TrimSpace(row[cm.emailIdx])
		if email == "" {
			return nil
		}
		result, err := conn.Enrichment(tomba.Params{"email": email})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "verify":
		if cm.emailIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		email := strings.TrimSpace(row[cm.emailIdx])
		if email == "" {
			return nil
		}
		result, err := conn.EmailVerifier(tomba.Params{"email": email})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "finder":
		if cm.domainIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		domain := strings.TrimSpace(row[cm.domainIdx])
		if domain == "" {
			return nil
		}
		params := tomba.Params{"domain": domain}
		if cm.firstIdx >= 0 && cm.firstIdx < len(row) {
			params["first_name"] = strings.TrimSpace(row[cm.firstIdx])
		}
		if cm.lastIdx >= 0 && cm.lastIdx < len(row) {
			params["last_name"] = strings.TrimSpace(row[cm.lastIdx])
		}
		result, err := conn.EmailFinder(params)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "search":
		if cm.domainIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		domain := strings.TrimSpace(row[cm.domainIdx])
		if domain == "" {
			return nil
		}
		result, err := conn.DomainSearch(tomba.Params{"domain": domain})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "author":
		if cm.urlIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		url := strings.TrimSpace(row[cm.urlIdx])
		if url == "" {
			return nil
		}
		result, err := conn.AuthorFinder(tomba.Params{"url": url})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "linkedin":
		if cm.urlIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		url := strings.TrimSpace(row[cm.urlIdx])
		if url == "" {
			return nil
		}
		result, err := conn.LinkedinFinder(tomba.Params{"url": url})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "phone":
		if cm.emailIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		email := strings.TrimSpace(row[cm.emailIdx])
		if email == "" {
			return nil
		}
		result, err := conn.Tomba.PhoneFinder(tomba.Params{"email": email})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "sources":
		if cm.emailIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		email := strings.TrimSpace(row[cm.emailIdx])
		if email == "" {
			return nil
		}
		result, err := conn.Tomba.Sources(email)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data
	}

	return nil
}

func writeBulkOutput(filename string, headers []string, rows [][]string, results []map[string]interface{}, opType string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Build enriched headers
	extraHeaders := getExtraHeaders(opType)
	allHeaders := append(headers, extraHeaders...)
	if err := writer.Write(allHeaders); err != nil {
		return err
	}

	for i, row := range rows {
		var result map[string]interface{}
		if i < len(results) {
			result = results[i]
		}
		extraCols := extractExtraCols(result, opType)
		allCols := append(row, extraCols...)
		if err := writer.Write(allCols); err != nil {
			return err
		}
	}

	return nil
}

func getExtraHeaders(opType string) []string {
	switch opType {
	case "enrich":
		return []string{"found_email", "first_name", "last_name", "position", "company", "country", "linkedin", "twitter"}
	case "verify":
		return []string{"result", "status", "score", "mx_found", "smtp_check"}
	case "finder":
		return []string{"found_email", "score", "position", "company"}
	case "search":
		return []string{"total_emails", "first_email"}
	case "author":
		return []string{"found_email", "first_name", "last_name", "position", "company", "country"}
	case "linkedin":
		return []string{"found_email", "first_name", "last_name", "position", "company", "country", "linkedin"}
	case "phone":
		return []string{"phone_number", "phone_type", "country_code", "country_name"}
	case "sources":
		return []string{"total_sources", "first_source_url", "first_source_domain"}
	default:
		return []string{}
	}
}

func extractExtraCols(result map[string]interface{}, opType string) []string {
	extraCount := map[string]int{
		"enrich":   8,
		"verify":   5,
		"finder":   4,
		"search":   2,
		"author":   6,
		"linkedin": 7,
		"phone":    4,
		"sources":  3,
	}
	count := extraCount[opType]

	if result == nil {
		return make([]string, count)
	}

	if _, hasError := result["_error"]; hasError {
		errMsg := fmt.Sprintf("%v", result["_error"])
		cols := make([]string, count)
		cols[0] = errMsg
		return cols
	}

	switch opType {
	case "enrich":
		d := getNestedMap(result, "data")
		return []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
			getMapStr(d, "country"),
			getMapStr(d, "linkedin"),
			getMapStr(d, "twitter"),
		}
	case "verify":
		d := getNestedMap(result, "data")
		e := getNestedMap(d, "email")
		return []string{
			getMapStr(e, "result"),
			getMapStr(e, "status"),
			getMapFloat(e, "score"),
			getMapBool(e, "mx_found"),
			getMapBool(e, "smtp_check"),
		}
	case "finder":
		d := getNestedMap(result, "data")
		return []string{
			getMapStr(d, "email"),
			getMapFloat(d, "score"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
		}
	case "search":
		m := getNestedMap(result, "meta")
		d := getNestedMap(result, "data")
		firstEmail := ""
		if emails, ok := d["emails"].([]interface{}); ok && len(emails) > 0 {
			if em, ok := emails[0].(map[string]interface{}); ok {
				firstEmail = getMapStr(em, "email")
			}
		}
		return []string{getMapFloat(m, "total"), firstEmail}
	case "author":
		d := getNestedMap(result, "data")
		return []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
			getMapStr(d, "country"),
		}
	case "linkedin":
		d := getNestedMap(result, "data")
		return []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
			getMapStr(d, "country"),
			getMapStr(d, "linkedin"),
		}
	case "phone":
		d := getNestedMap(result, "data")
		phone := ""
		phoneType := ""
		countryCode := ""
		countryName := ""
		if phones, ok := d["phones"].([]interface{}); ok && len(phones) > 0 {
			if p, ok := phones[0].(map[string]interface{}); ok {
				phone = getMapStr(p, "phone")
				phoneType = getMapStr(p, "type")
				countryCode = getMapStr(p, "country_code")
				countryName = getMapStr(p, "country_name")
			}
		}
		return []string{phone, phoneType, countryCode, countryName}
	case "sources":
		d := getNestedMap(result, "data")
		totalSources := ""
		firstURL := ""
		firstDomain := ""
		if sources, ok := d["sources"].([]interface{}); ok {
			totalSources = fmt.Sprintf("%d", len(sources))
			if len(sources) > 0 {
				if s, ok := sources[0].(map[string]interface{}); ok {
					firstURL = getMapStr(s, "url")
					firstDomain = getMapStr(s, "domain")
				}
			}
		}
		return []string{totalSources, firstURL, firstDomain}
	default:
		return []string{}
	}
}

func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func getMapStr(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getMapFloat(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(float64); ok {
		return fmt.Sprintf("%.0f", v)
	}
	return ""
}

func getMapBool(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(bool); ok {
		if v {
			return "true"
		}
		return "false"
	}
	return ""
}
