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
	bulkUrlCol       string
	bulkFullNameCol  string
	bulkConcurrency  int
	bulkNoResume     bool
	bulkEnrichMobile   bool
	bulkFull           bool
	bulkPhoneCol       string
	bulkCountryCodeCol string
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
	bulkCmd.Flags().StringVar(&bulkType, "type", "enrich", "Operation type: enrich, verify, finder, search, author, linkedin, phone, phone-validator, company, similar, sources.")
	bulkCmd.Flags().StringVar(&bulkColumn, "column", "", "Column name for email or domain (auto-detected if empty).")
	bulkCmd.Flags().StringVar(&bulkDomainCol, "domain-col", "", "Column name for domain (for finder type).")
	bulkCmd.Flags().StringVar(&bulkFirstCol, "first-col", "", "Column name for first name (for finder type).")
	bulkCmd.Flags().StringVar(&bulkLastCol, "last-col", "", "Column name for last name (for finder type).")
	bulkCmd.Flags().StringVar(&bulkUrlCol, "url-col", "", "Column name for URL (for author/linkedin type).")
	bulkCmd.Flags().StringVar(&bulkFullNameCol, "full-name-col", "", "Column name for full name (for finder type, alternative to first-col + last-col).")
	bulkCmd.Flags().BoolVar(&bulkEnrichMobile, "enrich-mobile", false, "Get the phone number associated with the email address found (for finder/enrich type).")
	bulkCmd.Flags().BoolVar(&bulkFull, "full", false, "Get all phone numbers (for phone type).")
	bulkCmd.Flags().StringVar(&bulkPhoneCol, "phone-col", "", "Column name for phone number (for phone-validator type).")
	bulkCmd.Flags().StringVar(&bulkCountryCodeCol, "country-code-col", "", "Column name for country code (for phone-validator type).")
	bulkCmd.Flags().IntVar(&bulkConcurrency, "concurrency", 0, "Number of concurrent workers (0=auto from plan, max 60).")
	bulkCmd.Flags().BoolVar(&bulkNoResume, "no-resume", false, "Ignore existing output and start fresh.")
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

// bulkSkipped is a sentinel result indicating the row was skipped (empty input column).
var bulkSkipped = map[string]interface{}{"_skipped": true}

type bulkJob struct {
	index int
	row   []string
}

type bulkResult struct {
	index  int
	result map[string]any
}

// StreamingCSVWriter writes CSV rows incrementally as results arrive
type StreamingCSVWriter struct {
	mu     sync.Mutex
	file   *os.File
	writer *csv.Writer
	rows   [][]string
	opType string
}

// NewStreamingCSVWriter opens the output file and writes headers if needed
func NewStreamingCSVWriter(filename string, inputHeaders []string, rows [][]string, opType string, appendMode bool) (*StreamingCSVWriter, error) {
	var file *os.File
	var err error

	if appendMode {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	if !appendMode {
		extraHeaders := getExtraHeaders(opType)
		allHeaders := append(inputHeaders, extraHeaders...)
		if err := writer.Write(allHeaders); err != nil {
			_ = file.Close()
			return nil, err
		}
		writer.Flush()
	}

	return &StreamingCSVWriter{
		file:   file,
		writer: writer,
		rows:   rows,
		opType: opType,
	}, nil
}

// WriteResult writes a single result row to the CSV file (thread-safe)
func (w *StreamingCSVWriter) WriteResult(index int, result map[string]any) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Phone with --full: write one row per phone number
	if w.opType == "phone" && bulkFull && result != nil {
		if _, skipped := result["_skipped"]; skipped {
			count := len(getExtraHeaders(w.opType))
			allCols := append(w.rows[index], make([]string, count)...)
			_ = w.writer.Write(allCols)
			w.writer.Flush()
			return
		}
		rows := extractPhoneRows(result)
		if len(rows) == 0 {
			count := len(getExtraHeaders(w.opType))
			allCols := append(w.rows[index], make([]string, count)...)
			_ = w.writer.Write(allCols)
		} else {
			for _, phoneCols := range rows {
				allCols := append(w.rows[index], phoneCols...)
				_ = w.writer.Write(allCols)
			}
		}
		w.writer.Flush()
		return
	}

	extraCols := extractExtraCols(result, w.opType)
	allCols := append(w.rows[index], extraCols...)
	_ = w.writer.Write(allCols)
	w.writer.Flush()
}

// Close flushes and closes the underlying file
func (w *StreamingCSVWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writer.Flush()
	return w.file.Close()
}

// detectResumableRows checks if the output file exists and returns indices of already-processed rows
func detectResumableRows(outputFile string, inputHeaders []string, rows [][]string, opType string) (completed map[int]bool, resumeFound int, resumeNotFound int, err error) {
	if _, statErr := os.Stat(outputFile); os.IsNotExist(statErr) {
		return nil, 0, 0, nil
	}

	file, err := os.Open(outputFile)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable field count (handle truncated last row)

	allRecords, err := reader.ReadAll()
	if err != nil {
		return nil, 0, 0, err
	}
	if len(allRecords) < 2 {
		return nil, 0, 0, nil
	}

	// Validate headers match
	existingHeaders := allRecords[0]
	expectedHeaders := append(inputHeaders, getExtraHeaders(opType)...)
	if len(existingHeaders) != len(expectedHeaders) {
		return nil, 0, 0, fmt.Errorf("header mismatch: expected %d columns, got %d", len(expectedHeaders), len(existingHeaders))
	}
	for i, h := range expectedHeaders {
		if existingHeaders[i] != h {
			return nil, 0, 0, fmt.Errorf("header mismatch at column %d: expected %q, got %q", i, h, existingHeaders[i])
		}
	}

	// Build a set of completed input row keys
	inputColCount := len(inputHeaders)
	extraColCount := len(getExtraHeaders(opType))
	completedKeys := make(map[string][]string) // key -> extra cols
	for _, record := range allRecords[1:] {
		if len(record) < inputColCount {
			continue
		}
		key := strings.Join(record[:inputColCount], "\x00")
		completedKeys[key] = record[inputColCount:]
	}

	// Match against input rows
	completed = make(map[int]bool)
	found := 0
	notFound := 0
	for i, row := range rows {
		key := strings.Join(row, "\x00")
		if extraCols, ok := completedKeys[key]; ok {
			completed[i] = true
			// Determine if it was found or not by checking if extra cols are all empty
			hasData := false
			if len(extraCols) >= extraColCount {
				for _, col := range extraCols {
					if col != "" {
						hasData = true
						break
					}
				}
			}
			if hasData {
				found++
			} else {
				notFound++
			}
			delete(completedKeys, key) // handle duplicates: only match once
		}
	}

	return completed, found, notFound, nil
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

	// Check for resume
	var completedIndices map[int]bool
	var resumeFound, resumeNotFound int
	if !bulkNoResume {
		completedIndices, resumeFound, resumeNotFound, err = detectResumableRows(outputFile, headers, rows, bulkType)
		if err != nil {
			fmt.Printf("%s Could not parse existing output for resume: %s\n", util.WarningIcon(), util.Yellow(err.Error()))
			fmt.Printf("%s Starting fresh run\n", util.InfoIcon())
			completedIndices = nil
		}
	}

	resumeCount := len(completedIndices)
	remainingCount := len(rows) - resumeCount
	appendMode := resumeCount > 0

	if resumeCount > 0 {
		fmt.Printf("  %s Resuming: %s rows already processed, %s remaining\n\n",
			util.InfoIcon(),
			util.Green(fmt.Sprintf("%d", resumeCount)),
			util.Bold(fmt.Sprintf("%d", remainingCount)))
	}

	if remainingCount <= 0 {
		fmt.Printf("  %s All rows already processed. Use --no-resume to start fresh.\n", util.SuccessIcon())
		return
	}

	// Determine concurrency
	workers := bulkConcurrency
	if workers <= 0 {
		workers = getConcurrency(planName)
	}
	if workers > 60 {
		workers = 60
	}
	if workers > remainingCount {
		workers = remainingCount
	}

	fmt.Printf("  Plan: %s | Workers: %s | Rows: %s\n\n",
		util.Bold(planName), util.Cyan(fmt.Sprintf("%d", workers)), util.Bold(fmt.Sprintf("%d", remainingCount)))

	// Open streaming CSV writer
	streamWriter, err := NewStreamingCSVWriter(outputFile, headers, rows, bulkType, appendMode)
	if err != nil {
		fmt.Printf("%s Error opening output file: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	startTime := time.Now()

	// Process rows concurrently
	progress := output.NewProgressBar(len(rows))
	stats := &output.BulkStats{
		Total:      len(rows),
		Found:      resumeFound,
		NotFound:   resumeNotFound,
		OutputFile: outputFile,
		OpType:     bulkType,
	}
	if bulkType == "verify" {
		stats.ExtraCounts = make(map[string]int)
	}

	// Pre-increment progress for resumed rows
	for i := 0; i < resumeCount; i++ {
		progress.Increment()
	}
	progress.Start()

	jobs := make(chan bulkJob, workers*2)
	resultsCh := make(chan bulkResult, workers*2)

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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

	// Send jobs (skip completed rows)
	go func() {
		for i, row := range rows {
			if completedIndices != nil && completedIndices[i] {
				continue
			}
			jobs <- bulkJob{index: i, row: row}
		}
		close(jobs)
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results and write to file immediately
	var mu sync.Mutex
	for r := range resultsCh {
		streamWriter.WriteResult(r.index, r.result)
		mu.Lock()
		if r.result != nil {
			if _, skipped := r.result["_skipped"]; skipped {
				stats.Skipped++
			} else if errVal, hasError := r.result["_error"]; hasError {
				stats.Errors++
				if errStr, ok := errVal.(string); ok {
					classifyBulkError(errStr, stats)
				}
			} else if bulkResultHasData(r.result, bulkType) {
				stats.Found++
				// Track verify breakdown
				if bulkType == "verify" {
					if d, ok := r.result["data"].(map[string]any); ok {
						if e, ok := d["email"].(map[string]any); ok {
							if result, ok := e["result"].(string); ok {
								stats.ExtraCounts[result]++
							}
						}
					}
				}
			} else {
				stats.NotFound++
			}
		} else {
			stats.NotFound++
		}
		mu.Unlock()
		progress.Increment()
	}

	_ = streamWriter.Close()

	elapsed := time.Since(startTime)
	stats.Elapsed = elapsed
	stats.PrintStatsTable()
	stats.PrintDistributionChart()
}

type columnMapping struct {
	emailIdx       int
	domainIdx      int
	firstIdx       int
	lastIdx        int
	urlIdx         int
	fullNameIdx    int
	phoneIdx       int
	countryCodeIdx int
}

func mapColumns(headers []string, opType string) *columnMapping {
	cm := &columnMapping{
		emailIdx:       -1,
		domainIdx:      -1,
		firstIdx:       -1,
		lastIdx:        -1,
		urlIdx:         -1,
		fullNameIdx:    -1,
		phoneIdx:       -1,
		countryCodeIdx: -1,
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
		case bulkFullNameCol != "" && strings.EqualFold(h, bulkFullNameCol):
			cm.fullNameIdx = i
		case bulkPhoneCol != "" && strings.EqualFold(h, bulkPhoneCol):
			cm.phoneIdx = i
		case bulkCountryCodeCol != "" && strings.EqualFold(h, bulkCountryCodeCol):
			cm.countryCodeIdx = i
		case cm.emailIdx == -1 && (lower == "email" || lower == "e-mail" || lower == "email_address" || lower == "emailaddress" || lower == "mail"):
			cm.emailIdx = i
		case cm.domainIdx == -1 && (lower == "domain" || lower == "company_domain" || lower == "website" || lower == "company_website"):
			cm.domainIdx = i
		case cm.firstIdx == -1 && (lower == "first_name" || lower == "firstname" || lower == "first" || lower == "fname"):
			cm.firstIdx = i
		case cm.lastIdx == -1 && (lower == "last_name" || lower == "lastname" || lower == "last" || lower == "lname"):
			cm.lastIdx = i
		case cm.fullNameIdx == -1 && (lower == "full_name" || lower == "fullname" || lower == "full name"):
			cm.fullNameIdx = i
		case cm.phoneIdx == -1 && (lower == "phone" || lower == "phone_number" || lower == "phonenumber" || lower == "tel" || lower == "telephone"):
			cm.phoneIdx = i
		case cm.countryCodeIdx == -1 && (lower == "country_code" || lower == "countrycode" || lower == "country"):
			cm.countryCodeIdx = i
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
		if cm.fullNameIdx == -1 && (cm.firstIdx == -1 || cm.lastIdx == -1) {
			if cm.firstIdx == -1 {
				cm.firstIdx = promptColumnSelect(headers, "first name")
			}
			if cm.lastIdx == -1 {
				cm.lastIdx = promptColumnSelect(headers, "last name")
			}
		}
		fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
		if cm.fullNameIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for full name\n", util.SuccessIcon(), util.Bold(headers[cm.fullNameIdx]))
		}
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
		if cm.emailIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
		} else if cm.domainIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
		} else if cm.urlIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for URL\n", util.SuccessIcon(), util.Bold(headers[cm.urlIdx]))
		} else {
			cm.emailIdx = promptColumnSelect(headers, "email")
			if cm.emailIdx == -1 {
				fmt.Printf("%s Could not find email, domain, or URL column. Use --column, --domain-col, or --url-col to specify.\n", util.ErrorIcon())
				return nil
			}
			fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
		}
	case "sources":
		if cm.emailIdx == -1 {
			cm.emailIdx = promptColumnSelect(headers, "email")
			if cm.emailIdx == -1 {
				fmt.Printf("%s Could not find email column. Use --column to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
	case "company", "similar":
		if cm.domainIdx == -1 {
			cm.domainIdx = promptColumnSelect(headers, "domain")
			if cm.domainIdx == -1 {
				fmt.Printf("%s Could not find domain column. Use --domain-col to specify.\n", util.ErrorIcon())
				return nil
			}
		}
		fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
	case "phone-validator":
		if cm.phoneIdx == -1 {
			if cm.emailIdx >= 0 {
				cm.phoneIdx = cm.emailIdx
			} else {
				cm.phoneIdx = promptColumnSelect(headers, "phone number")
				if cm.phoneIdx == -1 {
					fmt.Printf("%s Could not find phone column. Use --phone-col to specify.\n", util.ErrorIcon())
					return nil
				}
			}
		}
		fmt.Printf("  %s Using column '%s' for phone number\n", util.SuccessIcon(), util.Bold(headers[cm.phoneIdx]))
		if cm.countryCodeIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for country code\n", util.SuccessIcon(), util.Bold(headers[cm.countryCodeIdx]))
		}
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
			return bulkSkipped
		}
		params := tomba.Params{"email": email}
		if bulkEnrichMobile {
			params["enrich_mobile"] = true
		}
		result, err := conn.Enrichment(params)
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
			return bulkSkipped
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
			return bulkSkipped
		}
		params := tomba.Params{"domain": domain}
		if cm.fullNameIdx >= 0 && cm.fullNameIdx < len(row) {
			params["full_name"] = strings.TrimSpace(row[cm.fullNameIdx])
		} else {
			if cm.firstIdx >= 0 && cm.firstIdx < len(row) {
				params["first_name"] = strings.TrimSpace(row[cm.firstIdx])
			}
			if cm.lastIdx >= 0 && cm.lastIdx < len(row) {
				params["last_name"] = strings.TrimSpace(row[cm.lastIdx])
			}
		}
		if bulkEnrichMobile {
			params["enrich_mobile"] = true
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
			return bulkSkipped
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
			return bulkSkipped
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
			return bulkSkipped
		}
		params := tomba.Params{"url": url}
		if bulkEnrichMobile {
			params["enrich_mobile"] = true
		}
		result, err := conn.LinkedinFinder(params)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "phone":
		params := tomba.Params{}
		if cm.emailIdx >= 0 && cm.emailIdx < len(row) {
			v := strings.TrimSpace(row[cm.emailIdx])
			if v == "" {
				return bulkSkipped
			}
			params["email"] = v
		} else if cm.domainIdx >= 0 && cm.domainIdx < len(row) {
			v := strings.TrimSpace(row[cm.domainIdx])
			if v == "" {
				return bulkSkipped
			}
			params["domain"] = v
		} else if cm.urlIdx >= 0 && cm.urlIdx < len(row) {
			v := strings.TrimSpace(row[cm.urlIdx])
			if v == "" {
				return bulkSkipped
			}
			params["linkedin"] = v
		} else {
			return bulkSkipped
		}
		if bulkFull {
			params["full"] = true
		}
		result, err := conn.Tomba.PhoneFinder(params)
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
			return bulkSkipped
		}
		result, err := conn.Tomba.Sources(email)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "company":
		if cm.domainIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		domain := strings.TrimSpace(row[cm.domainIdx])
		if domain == "" {
			return bulkSkipped
		}
		result, err := conn.CompanyFind(tomba.Params{"domain": domain})
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "similar":
		if cm.domainIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		domain := strings.TrimSpace(row[cm.domainIdx])
		if domain == "" {
			return bulkSkipped
		}
		result, err := conn.SimilarDomains(domain)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		raw, _ := result.Marshal()
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		return data

	case "phone-validator":
		if cm.phoneIdx >= len(row) {
			return map[string]interface{}{"_error": "index out of range"}
		}
		phone := strings.TrimSpace(row[cm.phoneIdx])
		if phone == "" {
			return bulkSkipped
		}
		params := tomba.Params{"phone": phone}
		if cm.countryCodeIdx >= 0 && cm.countryCodeIdx < len(row) {
			cc := strings.TrimSpace(row[cm.countryCodeIdx])
			if cc != "" {
				params["country_code"] = cc
			}
		}
		result, err := conn.Tomba.PhoneValidator(params)
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


func getExtraHeaders(opType string) []string {
	switch opType {
	case "enrich":
		headers := []string{"found_email", "first_name", "last_name", "position", "company", "country", "linkedin", "twitter"}
		if bulkEnrichMobile {
			headers = append(headers, "phone_number")
		}
		return headers
	case "verify":
		return []string{"result", "status", "score", "mx_found", "smtp_check"}
	case "finder":
		headers := []string{"found_email", "first_name", "last_name", "score", "position", "company"}
		if bulkEnrichMobile {
			headers = append(headers, "phone_number")
		}
		return headers
	case "search":
		return []string{"total_emails", "first_email"}
	case "author":
		return []string{"found_email", "first_name", "last_name", "position", "company", "country"}
	case "linkedin":
		headers := []string{"found_email", "first_name", "last_name", "position", "company", "country", "linkedin"}
		if bulkEnrichMobile {
			headers = append(headers, "phone_number")
		}
		return headers
	case "phone":
		return []string{"phone_number", "valid", "country_code", "line_type", "carrier"}
	case "sources":
		return []string{"total_sources", "first_source_url", "first_source_domain"}
	case "company":
		return []string{"company_name", "company_domain", "industry", "country", "size", "founded"}
	case "similar":
		return []string{"total_similar", "first_similar_domain"}
	case "phone-validator":
		return []string{"valid", "local_format", "intl_format", "country_code", "line_type", "carrier"}
	default:
		return []string{}
	}
}

func extractExtraCols(result map[string]interface{}, opType string) []string {
	count := len(getExtraHeaders(opType))

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
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
			getMapStr(d, "country"),
			getMapStr(d, "linkedin"),
			getMapStr(d, "twitter"),
		}
		if bulkEnrichMobile {
			cols = append(cols, getFirstPhone(d))
		}
		return cols
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
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapFloat(d, "score"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
		}
		if bulkEnrichMobile {
			cols = append(cols, getFirstPhone(d))
		}
		return cols
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
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "company"),
			getMapStr(d, "country"),
			getMapStr(d, "linkedin"),
		}
		if bulkEnrichMobile {
			cols = append(cols, getFirstPhone(d))
		}
		return cols
	case "phone":
		phone := extractFirstPhone(result)
		if phone == nil {
			return make([]string, count)
		}
		return []string{
			getMapStr(phone, "intl_format"),
			getMapBool(phone, "valid"),
			getMapStr(phone, "country_code"),
			getMapStr(phone, "line_type"),
			getMapStr(phone, "carrier"),
		}
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
	case "company":
		d := getNestedMap(result, "data")
		org := getNestedMap(d, "organization")
		return []string{
			getMapStr(org, "name"),
			getMapStr(org, "website_url"),
			getMapStr(org, "industries"),
			getMapStr(org, "country"),
			getMapStr(org, "size"),
			getMapStr(org, "founded"),
		}
	case "similar":
		totalSimilar := ""
		firstDomain := ""
		if domains, ok := result["data"].([]interface{}); ok {
			totalSimilar = fmt.Sprintf("%d", len(domains))
			if len(domains) > 0 {
				if s, ok := domains[0].(map[string]interface{}); ok {
					firstDomain = getMapStr(s, "website_url")
				}
			}
		}
		return []string{totalSimilar, firstDomain}
	case "phone-validator":
		d := getNestedMap(result, "data")
		return []string{
			getMapBool(d, "valid"),
			getMapStr(d, "local_format"),
			getMapStr(d, "intl_format"),
			getMapStr(d, "country_code"),
			getMapStr(d, "line_type"),
			getMapStr(d, "carrier"),
		}
	default:
		return []string{}
	}
}

func classifyBulkError(errStr string, stats *output.BulkStats) {
	// SDK error format: "error: <status>, status code: <code>"
	if idx := strings.LastIndex(errStr, "status code: "); idx >= 0 {
		codeStr := strings.TrimSpace(errStr[idx+len("status code: "):])
		if len(codeStr) >= 3 {
			switch codeStr[0] {
			case '4':
				stats.ClientErrors++
				return
			case '5':
				stats.ServerErrors++
				return
			}
		}
	}
	// Non-HTTP errors (timeout, network, etc.) count as client errors
	stats.ClientErrors++
}

func bulkResultHasData(result map[string]interface{}, opType string) bool {
	d := getNestedMap(result, "data")
	switch opType {
	case "enrich":
		return getMapStr(d, "email") != ""
	case "verify":
		e := getNestedMap(d, "email")
		return getMapStr(e, "result") != ""
	case "finder":
		return getMapStr(d, "email") != ""
	case "search":
		m := getNestedMap(result, "meta")
		if v, ok := m["total"].(float64); ok {
			return v > 0
		}
		return false
	case "author":
		return getMapStr(d, "email") != ""
	case "linkedin":
		return getMapStr(d, "email") != ""
	case "phone":
		switch v := result["data"].(type) {
		case map[string]interface{}:
			return getMapStr(v, "intl_format") != ""
		case []interface{}:
			return len(v) > 0
		}
		return false
	case "sources":
		if sources, ok := d["sources"].([]interface{}); ok {
			return len(sources) > 0
		}
		return false
	case "company":
		org := getNestedMap(d, "organization")
		return getMapStr(org, "name") != "" || getMapStr(org, "website_url") != ""
	case "similar":
		if domains, ok := result["data"].([]interface{}); ok {
			return len(domains) > 0
		}
		return false
	case "phone-validator":
		return getMapStr(d, "intl_format") != "" || getMapStr(d, "local_format") != ""
	}
	return true
}

func extractFirstPhone(result map[string]interface{}) map[string]interface{} {
	if result == nil {
		return nil
	}
	switch v := result["data"].(type) {
	case map[string]interface{}:
		return v
	case []interface{}:
		if len(v) > 0 {
			if m, ok := v[0].(map[string]interface{}); ok {
				return m
			}
		}
	}
	return nil
}

func extractPhoneRows(result map[string]interface{}) [][]string {
	if result == nil {
		return nil
	}
	if _, hasError := result["_error"]; hasError {
		count := len(getExtraHeaders("phone"))
		cols := make([]string, count)
		cols[0] = fmt.Sprintf("%v", result["_error"])
		return [][]string{cols}
	}
	var phones []map[string]interface{}
	switch v := result["data"].(type) {
	case map[string]interface{}:
		phones = append(phones, v)
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				phones = append(phones, m)
			}
		}
	}
	var rows [][]string
	for _, p := range phones {
		rows = append(rows, []string{
			getMapStr(p, "intl_format"),
			getMapBool(p, "valid"),
			getMapStr(p, "country_code"),
			getMapStr(p, "line_type"),
			getMapStr(p, "carrier"),
		})
	}
	return rows
}

func getFirstPhone(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	if v := getMapStr(m, "intl_format"); v != "" {
		return v
	}
	if v := getMapStr(m, "phone_number"); v != "" {
		return v
	}
	if phones, ok := m["phone_data"].([]interface{}); ok && len(phones) > 0 {
		if p, ok := phones[0].(map[string]interface{}); ok {
			if v := getMapStr(p, "intl_format"); v != "" {
				return v
			}
			return getMapStr(p, "phone_number")
		}
	}
	return ""
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
