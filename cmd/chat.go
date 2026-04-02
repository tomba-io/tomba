package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var (
	chatModel string
	chatFile  string
	chatType  string
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:     "chat",
	Aliases: []string{"ai"},
	Short:   "Open AI chat with Tomba tools — auto-initializes your project on first run.",

	Run:     chatRun,
	Example: chatExample,
}

func init() {
	chatCmd.Flags().StringVar(&chatModel, "model", "gpt-4o", "OpenAI model to use.")
	chatCmd.Flags().StringVar(&chatFile, "file", "", "CSV file for bulk processing via chat.")
	chatCmd.Flags().StringVar(&chatType, "type", "enrich", "Bulk operation type: enrich, verify, finder, search, linkedin, author.")
}

// OpenAI API types
type chatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []toolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type toolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function functionCall `json:"function"`
}

type functionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Tools    []toolDef     `json:"tools"`
}

type toolDef struct {
	Type     string          `json:"type"`
	Function toolFunctionDef `json:"function"`
}

type toolFunctionDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type chatChoice struct {
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// chatRun the actual work chat
func chatRun(cmd *cobra.Command, args []string) {

	// Auto-initialize: check login
	init := start.New(conn)

	// Check for OpenAI API key
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		fmt.Println(util.ErrorIcon(), util.Red("OPENAI_API_KEY environment variable is required."))
		fmt.Println(util.InfoIcon(), "Set it with: export OPENAI_API_KEY=sk-...")
		return
	}

	// Get account info for greeting
	account, err := init.Tomba.Account()
	if err == nil {
		raw, _ := account.Marshal()
		var accData map[string]interface{}
		if json.Unmarshal(raw, &accData) == nil {
			if d, ok := accData["data"].(map[string]interface{}); ok {
				if email, ok := d["email"].(string); ok {
					fmt.Printf("%s Authenticated as %s\n", util.SuccessIcon(), util.Green(email))
				}
			}
		}
	}

	// If --file is provided, run bulk chat mode
	if chatFile != "" {
		chatBulkRun(init, openaiKey)
		return
	}

	fmt.Println()
	fmt.Println(util.Bold("Tomba AI Chat"), util.Gray("— type your question, 'quit' to exit"))
	fmt.Println(util.Gray("Example: Find the VP Sales at Stripe and get their email"))
	fmt.Println()

	// Define Tomba tools for OpenAI
	tools := getTombaTools()

	// Conversation history
	messages := buildSystemMessages()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(util.Bold("You: "))
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" || input == "q" {
			fmt.Println(util.Gray("Goodbye!"))
			break
		}

		messages = append(messages, chatMessage{Role: "user", Content: input})

		// Collect all tool calls silently, show only final result
		var toolChain []string
		fmt.Print(util.Gray("  Thinking..."))

		for {
			resp, err := callOpenAI(openaiKey, chatModel, messages, tools)
			if err != nil {
				fmt.Printf("\r%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
				break
			}

			if resp.Error != nil {
				fmt.Printf("\r%s %s\n", util.ErrorIcon(), util.Red(resp.Error.Message))
				break
			}

			if len(resp.Choices) == 0 {
				fmt.Printf("\r%s No response from AI\n", util.ErrorIcon())
				break
			}

			choice := resp.Choices[0]
			messages = append(messages, choice.Message)

			if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
				// Collect tool names silently
				for _, tc := range choice.Message.ToolCalls {
					toolChain = append(toolChain, tc.Function.Name)
				}

				// Execute tool calls silently
				for _, tc := range choice.Message.ToolCalls {
					result := executeTombaToolCall(init, tc.Function.Name, tc.Function.Arguments)
					messages = append(messages, chatMessage{
						Role:       "tool",
						Content:    result,
						ToolCallID: tc.ID,
					})
				}
				continue
			}

			// Clear "Thinking..." and show workflow + final response
			fmt.Print("\r                    \r")
			if len(toolChain) > 0 {
				fmt.Printf("  %s %s\n", util.Gray("workflow:"), util.Cyan(strings.Join(toolChain, " -> ")))
			}
			if choice.Message.Content != "" {
				fmt.Println()
				fmt.Println(choice.Message.Content)
				fmt.Println()
			}
			break
		}
	}
}

// chatBulkRun processes a CSV file using AI chat for each row
func chatBulkRun(conn *start.Conn, openaiKey string) {
	fmt.Println()

	// Open CSV file
	file, err := os.Open(chatFile)
	if err != nil {
		fmt.Printf("%s Cannot open file: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	defer file.Close()

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
	fmt.Printf("  Rows: %s\n", util.Bold(fmt.Sprintf("%d", len(rows))))
	fmt.Printf("  Operation: %s\n\n", util.Bold(chatType))

	tools := getTombaTools()

	// Build the prompt for AI to process the CSV
	csvSample := strings.Join(headers, ",") + "\n"
	for i, row := range rows {
		if i >= 3 { // show first 3 rows as sample
			break
		}
		csvSample += strings.Join(row, ",") + "\n"
	}

	// Collect all results, then show at the end
	progress := output.NewProgressBar(len(rows))
	stats := &output.BulkStats{
		Total: len(rows),
	}

	var allResults []map[string]interface{}

	for _, row := range rows {
		rowData := buildRowPrompt(headers, row, chatType)
		messages := buildSystemMessages()
		messages = append(messages, chatMessage{
			Role:    "user",
			Content: rowData,
		})

		result := processOneRow(conn, openaiKey, chatModel, messages, tools)
		allResults = append(allResults, result)

		if result != nil {
			if _, hasErr := result["_error"]; hasErr {
				stats.Errors++
			} else {
				stats.Found++
			}
		} else {
			stats.NotFound++
		}
		progress.Increment()
	}

	// Determine output file
	outputFile := conn.Output
	if outputFile == "" {
		outputFile = strings.TrimSuffix(chatFile, ".csv") + "_enriched.csv"
	}
	stats.OutputFile = outputFile

	// Write enriched CSV
	if err := writeChatBulkOutput(outputFile, headers, rows, allResults, chatType); err != nil {
		fmt.Printf("\n%s Error writing output: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	stats.PrintStats()

	// Print markdown summary table
	fmt.Println()
	printMarkdownSummary(headers, rows, allResults, chatType)
}

func buildRowPrompt(headers []string, row []string, opType string) string {
	var sb strings.Builder
	switch opType {
	case "enrich":
		for i, h := range headers {
			lower := strings.ToLower(h)
			if i < len(row) && (lower == "email" || lower == "e-mail" || lower == "mail" || lower == "email_address") {
				sb.WriteString(fmt.Sprintf("Enrich this email and return all available data: %s", row[i]))
				return sb.String()
			}
		}
		sb.WriteString(fmt.Sprintf("Enrich: %s", strings.Join(row, ", ")))
	case "verify":
		for i, h := range headers {
			lower := strings.ToLower(h)
			if i < len(row) && (lower == "email" || lower == "e-mail" || lower == "mail" || lower == "email_address") {
				sb.WriteString(fmt.Sprintf("Verify this email address: %s", row[i]))
				return sb.String()
			}
		}
		sb.WriteString(fmt.Sprintf("Verify: %s", strings.Join(row, ", ")))
	case "finder":
		sb.WriteString("Find the email for: ")
		for i, h := range headers {
			if i < len(row) && row[i] != "" {
				sb.WriteString(fmt.Sprintf("%s=%s ", h, row[i]))
			}
		}
	case "search":
		for i, h := range headers {
			lower := strings.ToLower(h)
			if i < len(row) && (lower == "domain" || lower == "company_domain" || lower == "website") {
				sb.WriteString(fmt.Sprintf("Search all emails for domain: %s", row[i]))
				return sb.String()
			}
		}
		sb.WriteString(fmt.Sprintf("Search: %s", strings.Join(row, ", ")))
	}
	return sb.String()
}

func processOneRow(conn *start.Conn, apiKey, model string, messages []chatMessage, tools []toolDef) map[string]interface{} {
	for attempts := 0; attempts < 10; attempts++ {
		resp, err := callOpenAI(apiKey, model, messages, tools)
		if err != nil {
			return map[string]interface{}{"_error": err.Error()}
		}
		if resp.Error != nil {
			return map[string]interface{}{"_error": resp.Error.Message}
		}
		if len(resp.Choices) == 0 {
			return map[string]interface{}{"_error": "no response"}
		}

		choice := resp.Choices[0]
		messages = append(messages, choice.Message)

		if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
			for _, tc := range choice.Message.ToolCalls {
				result := executeTombaToolCall(conn, tc.Function.Name, tc.Function.Arguments)
				messages = append(messages, chatMessage{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})
				// Parse the tool result as our row result
				var data map[string]interface{}
				if json.Unmarshal([]byte(result), &data) == nil {
					if _, hasErr := data["error"]; !hasErr {
						// Keep going, the AI might call more tools
					}
				}
			}
			continue
		}

		// Conversation finished — extract the last tool result
		// Walk backwards through messages to find the last tool response with real data
		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i].Role == "tool" {
				var data map[string]interface{}
				if json.Unmarshal([]byte(messages[i].Content), &data) == nil {
					if _, hasErr := data["error"]; !hasErr {
						return data
					}
				}
			}
		}
		return nil
	}
	return map[string]interface{}{"_error": "max attempts reached"}
}

func writeChatBulkOutput(filename string, headers []string, rows [][]string, results []map[string]interface{}, opType string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

func printMarkdownSummary(headers []string, rows [][]string, results []map[string]interface{}, opType string) {
	extra := getExtraHeaders(opType)
	allHeaders := append(headers, extra...)

	// Limit to reasonable column count for display
	maxCols := len(allHeaders)
	if maxCols > 12 {
		maxCols = 12
	}
	displayHeaders := allHeaders[:maxCols]

	// Print markdown table header
	fmt.Print("| ")
	for _, h := range displayHeaders {
		fmt.Printf("%s | ", util.Bold(h))
	}
	fmt.Println()

	// Separator
	fmt.Print("| ")
	for range displayHeaders {
		fmt.Print("--- | ")
	}
	fmt.Println()

	// Rows (limit to 20 for display)
	maxRows := len(rows)
	if maxRows > 20 {
		maxRows = 20
	}
	for i := 0; i < maxRows; i++ {
		var result map[string]interface{}
		if i < len(results) {
			result = results[i]
		}
		extraCols := extractExtraCols(result, opType)
		allCols := append(rows[i], extraCols...)

		fmt.Print("| ")
		for j := 0; j < maxCols; j++ {
			val := ""
			if j < len(allCols) {
				val = allCols[j]
			}
			// Truncate long values
			if len(val) > 35 {
				val = val[:32] + "..."
			}
			fmt.Printf("%s | ", val)
		}
		fmt.Println()
	}

	if len(rows) > 20 {
		fmt.Printf("\n  %s ... and %d more rows (see output file)\n", util.Gray(""), len(rows)-20)
	}
}

func buildSystemMessages() []chatMessage {
	return []chatMessage{
		{
			Role: "system",
			Content: `You are Tomba AI assistant. You help users find email addresses, verify emails, enrich contact data, and discover companies using the Tomba API tools.

## How to use tools

When a user asks to find someone's email:
1. Use domain_search to find emails at their company
2. Use email_finder if you know the person's name and domain
3. Use email_verifier to verify found emails
4. Use email_enrichment to get additional data about the person

When a user asks to verify an email, use email_verifier.
When a user asks to enrich, use email_enrichment.
When a user asks about a domain, use domain_search or email_count.

## Response format

Always respond in **Markdown** format. Structure your response with:

- Use headers (##, ###) for sections
- Use **bold** for important values like email addresses, names
- Use tables when showing multiple results
- Use bullet lists for details
- Include ALL available data fields in your response — do not omit any fields returned by the tools

Example response:

## Results

**Email found:** david.singleton@stripe.com
**Verified:** deliverable (score: 95)

### Contact Details
| Field | Value |
|-------|-------|
| Name | David Singleton |
| Position | VP Sales |
| Company | Stripe |
| Country | US |
| LinkedIn | linkedin.com/in/... |
| Twitter | @david |
| Score | 95 |

Always show every field that the API returns. Never summarize or skip fields.`,
		},
	}
}

func callOpenAI(apiKey, model string, messages []chatMessage, tools []toolDef) (*chatResponse, error) {
	reqBody := chatRequest{
		Model:    model,
		Messages: messages,
		Tools:    tools,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResp, nil
}

func executeTombaToolCall(conn *start.Conn, name, arguments string) string {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf(`{"error": "invalid arguments: %s"}`, err.Error())
	}

	switch name {
	case "domain_search":
		domain := getArgStr(args, "domain")
		params := tomba.Params{"domain": domain}
		if page := getArgStr(args, "page"); page != "" {
			params["page"] = page
		}
		if limit := getArgStr(args, "limit"); limit != "" {
			params["limit"] = limit
		}
		result, err := conn.Tomba.DomainSearch(params)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "email_finder":
		domain := getArgStr(args, "domain")
		params := tomba.Params{"domain": domain}
		if fn := getArgStr(args, "first_name"); fn != "" {
			params["first_name"] = fn
		}
		if ln := getArgStr(args, "last_name"); ln != "" {
			params["last_name"] = ln
		}
		if full := getArgStr(args, "full_name"); full != "" {
			params["full_name"] = full
		}
		result, err := conn.Tomba.EmailFinder(params)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "email_enrichment":
		email := getArgStr(args, "email")
		result, err := conn.Tomba.Enrichment(tomba.Params{"email": email})
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "email_verifier":
		email := getArgStr(args, "email")
		result, err := conn.Tomba.EmailVerifier(tomba.Params{"email": email})
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "email_count":
		domain := getArgStr(args, "domain")
		result, err := conn.Tomba.Count(domain)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "email_sources":
		email := getArgStr(args, "email")
		result, err := conn.Tomba.Sources(email)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "linkedin_finder":
		url := getArgStr(args, "url")
		result, err := conn.Tomba.LinkedinFinder(tomba.Params{"url": url})
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "author_finder":
		url := getArgStr(args, "url")
		result, err := conn.Tomba.AuthorFinder(url)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "domain_status":
		domain := getArgStr(args, "domain")
		result, err := conn.Tomba.Status(domain)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	case "phone_finder":
		params := tomba.Params{}
		if email := getArgStr(args, "email"); email != "" {
			params["email"] = email
		}
		if domain := getArgStr(args, "domain"); domain != "" {
			params["domain"] = domain
		}
		result, err := conn.Tomba.PhoneFinder(params)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		raw, _ := result.Marshal()
		return string(raw)

	default:
		return fmt.Sprintf(`{"error": "unknown tool: %s"}`, name)
	}
}

func getArgStr(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	if v, ok := args[key].(float64); ok {
		return fmt.Sprintf("%.0f", v)
	}
	return ""
}

func getTombaTools() []toolDef {
	return []toolDef{
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "domain_search",
				Description: "Search for email addresses from a company domain. Returns a list of emails found for the domain.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"domain": map[string]interface{}{
							"type":        "string",
							"description": "The domain name to search (e.g., stripe.com)",
						},
						"page": map[string]interface{}{
							"type":        "string",
							"description": "Page number for pagination",
						},
						"limit": map[string]interface{}{
							"type":        "string",
							"description": "Max results per page (10, 20, or 50)",
						},
					},
					"required": []string{"domain"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "email_finder",
				Description: "Find the most likely email address for a person given their domain and name.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"domain": map[string]interface{}{
							"type":        "string",
							"description": "Company domain name (e.g., stripe.com)",
						},
						"first_name": map[string]interface{}{
							"type":        "string",
							"description": "Person's first name",
						},
						"last_name": map[string]interface{}{
							"type":        "string",
							"description": "Person's last name",
						},
						"full_name": map[string]interface{}{
							"type":        "string",
							"description": "Person's full name (alternative to first_name + last_name)",
						},
					},
					"required": []string{"domain"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "email_enrichment",
				Description: "Enrich an email address with additional data like name, position, company, social profiles.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address to enrich",
						},
					},
					"required": []string{"email"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "email_verifier",
				Description: "Verify the deliverability of an email address. Returns validation status and score.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address to verify",
						},
					},
					"required": []string{"email"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "email_count",
				Description: "Count the total number of email addresses available for a domain.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"domain": map[string]interface{}{
							"type":        "string",
							"description": "Domain name to count emails for",
						},
					},
					"required": []string{"domain"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "email_sources",
				Description: "Find web sources where an email address has been found.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address to find sources for",
						},
					},
					"required": []string{"email"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "linkedin_finder",
				Description: "Find the email address associated with a LinkedIn profile URL.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "LinkedIn profile URL",
						},
					},
					"required": []string{"url"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "author_finder",
				Description: "Find the email address of an article's author from the article URL.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "Article URL to find the author's email",
						},
					},
					"required": []string{"url"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "domain_status",
				Description: "Check if a domain is a webmail provider or disposable email provider.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"domain": map[string]interface{}{
							"type":        "string",
							"description": "Domain name to check status",
						},
					},
					"required": []string{"domain"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunctionDef{
				Name:        "phone_finder",
				Description: "Find phone numbers associated with an email address or domain.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address to find phone for",
						},
						"domain": map[string]interface{}{
							"type":        "string",
							"description": "Domain to find phone numbers for",
						},
					},
				},
			},
		},
	}
}
