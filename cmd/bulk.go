package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
	colMap *columnMapping
}

// NewStreamingCSVWriter opens the output file and writes headers if needed
func NewStreamingCSVWriter(filename string, inputHeaders []string, rows [][]string, opType string, appendMode bool, colMap *columnMapping) (*StreamingCSVWriter, error) {
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
		colMap: colMap,
	}, nil
}

// WriteResult writes a single result row to the CSV file (thread-safe)
func (w *StreamingCSVWriter) WriteResult(index int, result map[string]any) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Expand multi-row results: phone --full, search, and similar
	if result != nil {
		if _, skipped := result["_skipped"]; !skipped {
			var expandedRows [][]string
			if w.opType == "phone" && bulkFull {
				expandedRows = extractPhoneRows(result)
			} else if w.opType == "search" {
				expandedRows = extractSearchRows(result)
			} else if w.opType == "similar" {
				inputDomain := ""
				if w.colMap != nil && w.colMap.domainIdx >= 0 && w.colMap.domainIdx < len(w.rows[index]) {
					inputDomain = strings.TrimSpace(w.rows[index][w.colMap.domainIdx])
				}
				expandedRows = extractSimilarRows(result, inputDomain)
			}
			if expandedRows != nil {
				for _, cols := range expandedRows {
					allCols := append(w.rows[index], cols...)
					_ = w.writer.Write(allCols)
				}
				if len(expandedRows) == 0 {
					count := len(getExtraHeaders(w.opType))
					allCols := append(w.rows[index], make([]string, count)...)
					_ = w.writer.Write(allCols)
				}
				w.writer.Flush()
				return
			}
		}
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
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	var records [][]string
	lineNumber := 1
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			if parseErr, ok := readErr.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
				if len(record) > 0 {
					records = append(records, record)
				}
				lineNumber++
				continue
			}
			fmt.Printf("%s Cannot parse CSV line %d: %s\n", util.ErrorIcon(), lineNumber, util.Red(readErr.Error()))
			return
		}
		records = append(records, record)
		lineNumber++
	}

	if len(records) < 2 {
		fmt.Printf("%s CSV file must have a header row and at least one data row\n", util.ErrorIcon())
		return
	}

	headers := records[0]
	headerLen := len(headers)
	rows := records[1:]

	// Pad short rows to match header count
	for i, row := range rows {
		if len(row) < headerLen {
			padded := make([]string, headerLen)
			copy(padded, row)
			rows[i] = padded
		}
	}

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
	streamWriter, err := NewStreamingCSVWriter(outputFile, headers, rows, bulkType, appendMode, colMap)
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
		case bulkColumn != "" && cm.emailIdx == -1 && strings.EqualFold(h, bulkColumn):
			cm.emailIdx = i
		case bulkDomainCol != "" && cm.domainIdx == -1 && strings.EqualFold(h, bulkDomainCol):
			cm.domainIdx = i
		case bulkFirstCol != "" && cm.firstIdx == -1 && strings.EqualFold(h, bulkFirstCol):
			cm.firstIdx = i
		case bulkLastCol != "" && cm.lastIdx == -1 && strings.EqualFold(h, bulkLastCol):
			cm.lastIdx = i
		case bulkUrlCol != "" && cm.urlIdx == -1 && strings.EqualFold(h, bulkUrlCol):
			cm.urlIdx = i
		case bulkFullNameCol != "" && cm.fullNameIdx == -1 && strings.EqualFold(h, bulkFullNameCol):
			cm.fullNameIdx = i
		case bulkPhoneCol != "" && cm.phoneIdx == -1 && strings.EqualFold(h, bulkPhoneCol):
			cm.phoneIdx = i
		case bulkCountryCodeCol != "" && cm.countryCodeIdx == -1 && strings.EqualFold(h, bulkCountryCodeCol):
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
		// When user explicitly specifies columns, only clear auto-detected ones that conflict
		hasExplicit := bulkColumn != "" || bulkDomainCol != "" || bulkUrlCol != ""
		if hasExplicit {
			// Clear auto-detected columns that weren't explicitly requested
			if bulkColumn == "" {
				cm.emailIdx = -1
			}
			if bulkDomainCol == "" {
				cm.domainIdx = -1
			}
			if bulkUrlCol == "" {
				cm.urlIdx = -1
			}
		}
		found := false
		if cm.emailIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for email\n", util.SuccessIcon(), util.Bold(headers[cm.emailIdx]))
			found = true
		}
		if cm.domainIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for domain\n", util.SuccessIcon(), util.Bold(headers[cm.domainIdx]))
			found = true
		}
		if cm.urlIdx >= 0 {
			fmt.Printf("  %s Using column '%s' for URL\n", util.SuccessIcon(), util.Bold(headers[cm.urlIdx]))
			found = true
		}
		if !found {
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
		// Find email: try mapped column, then scan row
		if cm.emailIdx >= 0 && cm.emailIdx < len(row) {
			if v := strings.TrimSpace(row[cm.emailIdx]); v != "" && strings.Contains(v, "@") {
				params["email"] = v
			}
		}
		if _, ok := params["email"]; !ok && cm.emailIdx >= 0 {
			for _, cell := range row {
				if v := strings.TrimSpace(cell); v != "" && strings.Contains(v, "@") && strings.Contains(v, ".") {
					params["email"] = v
					break
				}
			}
		}
		// Find domain: try mapped column, then scan row
		if cm.domainIdx >= 0 && cm.domainIdx < len(row) {
			if v := strings.TrimSpace(row[cm.domainIdx]); v != "" && strings.Contains(v, ".") && !strings.Contains(v, " ") && !strings.Contains(v, "@") {
				params["domain"] = v
			}
		}
		// Find linkedin: always scan entire row for linkedin URL (handles shifted columns)
		if cm.urlIdx >= 0 {
			for _, cell := range row {
				if v := strings.TrimSpace(cell); v != "" && strings.Contains(v, "linkedin.com/") {
					params["linkedin"] = v
					break
				}
			}
		}
		if len(params) == 0 {
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
		return []string{"Email address", "Domain name", "Organization", "Confidence score",
			"Verification Status", "Type", "Pattern", "First Name", "Last Name", "Position", "Department",
			"Twitter URL", "LinkedIn URL", "Country",
			"Phone", "Phone2", "Phone3", "Phone4", "Phone5", "Phone6", "Phone7", "Phone8", "Phone9", "Phone10"}
	case "verify":
		return []string{"Email address", "Status", "Regexp", "Gibberish", "Disposable",
			"Webmail", "Mx records", "Smtp server", "Smtp check", "Accept all", "Block", "Score",
			"Phone", "Phone2", "Phone3", "Phone4", "Phone5", "Phone6", "Phone7", "Phone8", "Phone9", "Phone10"}
	case "finder":
		return []string{"Email address", "First Name", "Last Name", "Confidence score", "Verification status",
			"Phone", "Phone2", "Phone3", "Phone4", "Phone5", "Phone6", "Phone7", "Phone8", "Phone9", "Phone10"}
	case "search":
		return []string{"Email address", "Domain name", "Organization", "Confidence score",
			"Verification Status", "Type", "Pattern", "First Name", "Last Name", "Position", "Department",
			"Twitter URL", "LinkedIn URL",
			"Phone", "Phone2", "Phone3", "Phone4", "Phone5", "Phone6", "Phone7", "Phone8", "Phone9", "Phone10"}
	case "author":
		return []string{"url", "Email address", "Domain name", "Organization",
			"Confidence score", "Verification Status", "First Name", "Last Name", "Position",
			"Twitter URL", "LinkedIn URL", "Phone", "Title", "Description"}
	case "linkedin":
		return []string{"Email address", "Domain name", "Organization", "Confidence score",
			"Verification Status", "Type", "Pattern", "First Name", "Last Name", "Position", "Department",
			"Twitter URL",
			"Phone", "Phone2", "Phone3", "Phone4", "Phone5", "Phone6", "Phone7", "Phone8", "Phone9", "Phone10"}
	case "phone":
		if bulkFull {
			return []string{"Input Data", "IsValid",
				"Phone 1", "Phone 2", "Phone 3", "Phone 4", "Phone 5",
				"Phone 6", "Phone 7", "Phone 8", "Phone 9", "Phone 10",
				"Phone Local format", "Phone International format", "Phone Line type",
				"Carrier", "Country Code", "Timezone"}
		}
		return []string{"Input Data", "IsValid", "Phone Local format", "Phone International format",
			"Phone Line type", "Carrier", "Country Code", "Timezone"}
	case "sources":
		return []string{"total_sources", "first_source_url", "first_source_domain"}
	case "company":
		return []string{"Website", "Organization name", "Industries", "Employee count",
			"Twitter", "Facebook", "LinkedIn", "Country", "City", "Description",
			"Whois Registrar", "Whois Created date", "Whois Referral url",
			"phone", "phone2", "phone3", "phone4", "phone5", "phone6", "phone7", "phone8", "phone9", "phone10"}
	case "similar":
		return []string{"Domain name", "Organization", "Industries", "Similar To"}
	case "phone-validator":
		return []string{"Input", "IsValid", "Phone", "Phone Local format",
			"Phone International format", "Phone Line type", "Carrier", "Country Code", "Timezone"}
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
		return make([]string, count)
	}

	switch opType {
	case "enrich":
		d := getNestedMap(result, "data")
		v := getNestedMap(d, "verification")
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "domain"),
			getMapStr(d, "company"),
			getMapFloat(d, "score"),
			getMapStr(v, "status"),
			getMapStr(d, "type"),
			getMapStr(d, "pattern"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "department"),
			getMapStr(d, "twitter"),
			getMapStr(d, "linkedin"),
			getMapStr(d, "country"),
		}
		cols = append(cols, extractPhones(d, 10)...)
		return cols
	case "verify":
		d := getNestedMap(result, "data")
		e := getNestedMap(d, "email")
		cols := []string{
			getMapStr(e, "email"),
			getMapStr(e, "result"),
			getMapBool(e, "regexp"),
			getMapBool(e, "gibberish"),
			getMapBool(e, "disposable"),
			getMapBool(e, "webmail"),
			getMapBool(e, "mx_records"),
			getMapBool(e, "smtp_server"),
			getMapBool(e, "smtp_check"),
			getMapBool(e, "accept_all"),
			getMapBool(e, "block"),
			getMapFloat(e, "score"),
		}
		cols = append(cols, extractPhones(d, 10)...)
		return cols
	case "finder":
		d := getNestedMap(result, "data")
		v := getNestedMap(d, "verification")
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapFloat(d, "score"),
			getMapStr(v, "status"),
		}
		cols = append(cols, extractPhones(d, 10)...)
		return cols
	case "search":
		d := getNestedMap(result, "data")
		if emails, ok := d["emails"].([]interface{}); ok && len(emails) > 0 {
			if em, ok := emails[0].(map[string]interface{}); ok {
				return extractSearchEmailCols(em, d)
			}
		}
		return make([]string, count)
	case "author":
		d := getNestedMap(result, "data")
		v := getNestedMap(d, "verification")
		info := getNestedMap(d, "info")
		return []string{
			getMapStr(d, "url"),
			getMapStr(d, "email"),
			getMapStr(d, "domain"),
			getMapStr(d, "company"),
			getMapFloat(d, "score"),
			getMapStr(v, "status"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "twitter"),
			getMapStr(d, "linkedin"),
			getFirstPhone(d),
			getMapStr(info, "title"),
			getMapStr(info, "description"),
		}
	case "linkedin":
		d := getNestedMap(result, "data")
		v := getNestedMap(d, "verification")
		cols := []string{
			getMapStr(d, "email"),
			getMapStr(d, "domain"),
			getMapStr(d, "company"),
			getMapFloat(d, "score"),
			getMapStr(v, "status"),
			getMapStr(d, "type"),
			getMapStr(d, "pattern"),
			getMapStr(d, "first_name"),
			getMapStr(d, "last_name"),
			getMapStr(d, "position"),
			getMapStr(d, "department"),
			getMapStr(d, "twitter"),
		}
		cols = append(cols, extractPhones(d, 10)...)
		return cols
	case "phone":
		phone := extractFirstPhone(result)
		if phone == nil {
			return make([]string, count)
		}
		return []string{
			getMapStr(phone, "input"),
			getMapBool(phone, "valid"),
			getMapStr(phone, "local_format"),
			getMapStr(phone, "intl_format"),
			getMapStr(phone, "line_type"),
			getMapStr(phone, "carrier"),
			getMapStr(phone, "country_code"),
			getFirstTimezone(phone),
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
		loc := getNestedMap(org, "location")
		social := getNestedMap(org, "social_links")
		whois := getNestedMap(org, "whois")
		cols := []string{
			getMapStr(org, "website_url"),
			getMapStr(org, "organization"),
			getMapStr(org, "industries"),
			getMapStr(org, "size"),
			socialURL("https://twitter.com/", getMapStr(social, "twitter_url")),
			socialURL("https://facebook.com/", getMapStr(social, "facebook_url")),
			socialURL("https://linkedin.com/company/", getMapStr(social, "linkedin_url")),
			getMapStr(loc, "country"),
			getMapStr(loc, "city"),
			getMapStr(org, "description"),
			getMapStr(whois, "registrar_name"),
			getMapStr(whois, "created_date"),
			getMapStr(whois, "referral_url"),
		}
		cols = append(cols, extractPhones(org, 10)...)
		return cols
	case "similar":
		// Single-row fallback (multi-row expansion in extractSimilarRows handles the normal case)
		if domains, ok := result["data"].([]interface{}); ok && len(domains) > 0 {
			if s, ok := domains[0].(map[string]interface{}); ok {
				return []string{
					getMapStr(s, "website_url"),
					getMapStr(s, "name"),
					getMapStr(s, "industries"),
					"",
				}
			}
		}
		return make([]string, count)
	case "phone-validator":
		d := getNestedMap(result, "data")
		return []string{
			getMapStr(d, "input"),
			getMapBool(d, "valid"),
			getMapStr(d, "e164_format"),
			getMapStr(d, "local_format"),
			getMapStr(d, "intl_format"),
			getMapStr(d, "line_type"),
			getMapStr(d, "carrier"),
			getMapStr(d, "country_code"),
			getFirstTimezone(d),
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

func extractSearchEmailCols(em map[string]interface{}, parentData map[string]interface{}) []string {
	org := getNestedMap(parentData, "organization")
	v := getNestedMap(em, "verification")
	cols := []string{
		getMapStr(em, "email"),
		getMapStr(org, "website_url"),
		getMapStr(org, "organization"),
		getMapFloat(em, "score"),
		getMapStr(v, "status"),
		getMapStr(em, "type"),
		getMapStr(org, "pattern"),
		getMapStr(em, "first_name"),
		getMapStr(em, "last_name"),
		getMapStr(em, "position"),
		getMapStr(em, "department"),
		getMapStr(em, "twitter"),
		getMapStr(em, "linkedin"),
	}
	cols = append(cols, extractPhones(em, 10)...)
	return cols
}

func extractSearchRows(result map[string]interface{}) [][]string {
	if result == nil {
		return nil
	}
	if _, hasError := result["_error"]; hasError {
		count := len(getExtraHeaders("search"))
		return [][]string{make([]string, count)}
	}
	d := getNestedMap(result, "data")
	emails, ok := d["emails"].([]interface{})
	if !ok || len(emails) == 0 {
		return nil
	}
	var rows [][]string
	for _, item := range emails {
		if em, ok := item.(map[string]interface{}); ok {
			rows = append(rows, extractSearchEmailCols(em, d))
		}
	}
	return rows
}

func extractPhoneRows(result map[string]interface{}) [][]string {
	if result == nil {
		return nil
	}
	if _, hasError := result["_error"]; hasError {
		count := len(getExtraHeaders("phone"))
		return [][]string{make([]string, count)}
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
	if len(phones) == 0 {
		return nil
	}

	// Get input data from the first phone
	input := getMapStr(phones[0], "input")
	isValid := getMapBool(phones[0], "valid")

	// Phone 1-10 columns
	phoneNums := make([]string, 10)
	for i := 0; i < 10 && i < len(phones); i++ {
		phoneNums[i] = getMapStr(phones[i], "e164_format")
	}

	// Detail columns from the first phone
	first := phones[0]
	row := []string{input, isValid}
	row = append(row, phoneNums...)
	row = append(row,
		getMapStr(first, "local_format"),
		getMapStr(first, "intl_format"),
		getMapStr(first, "line_type"),
		getMapStr(first, "carrier"),
		getMapStr(first, "country_code"),
		getFirstTimezone(first),
	)
	return [][]string{row}
}

func extractSimilarRows(result map[string]interface{}, inputDomain string) [][]string {
	if result == nil {
		return nil
	}
	if _, hasError := result["_error"]; hasError {
		count := len(getExtraHeaders("similar"))
		return [][]string{make([]string, count)}
	}
	domains, ok := result["data"].([]interface{})
	if !ok || len(domains) == 0 {
		return nil
	}
	var rows [][]string
	for _, item := range domains {
		if s, ok := item.(map[string]interface{}); ok {
			rows = append(rows, []string{
				getMapStr(s, "website_url"),
				getMapStr(s, "name"),
				getMapStr(s, "industries"),
				inputDomain,
			})
		}
	}
	return rows
}

func extractPhones(m map[string]interface{}, maxCount int) []string {
	phones := make([]string, maxCount)
	if m == nil {
		return phones
	}
	if pd, ok := m["phone_data"].([]interface{}); ok {
		for i := 0; i < maxCount && i < len(pd); i++ {
			if p, ok := pd[i].(map[string]interface{}); ok {
				if v := getMapStr(p, "e164_format"); v != "" {
					phones[i] = v
				} else {
					phones[i] = getMapStr(p, "intl_format")
				}
			}
		}
	}
	return phones
}

func getFirstTimezone(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	if tzs, ok := m["timezones"].([]interface{}); ok && len(tzs) > 0 {
		if tz, ok := tzs[0].(string); ok {
			return tz
		}
	}
	return getMapStr(m, "timezone")
}

func socialURL(prefix, handle string) string {
	if handle == "" {
		return ""
	}
	return prefix + handle
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
