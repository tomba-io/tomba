package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/tomba-io/tomba/pkg/util"
)

// ProgressBar renders a progress bar to the terminal
type ProgressBar struct {
	Total       int
	Current     int
	Width       int
	StartTime   time.Time
	StartOffset int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total: total,
		Width: 30,
	}
}

// Start marks the beginning of processing (call after resume pre-increments)
func (p *ProgressBar) Start() {
	p.StartOffset = p.Current
	p.StartTime = time.Now()
}

// Increment advances the progress bar by one
func (p *ProgressBar) Increment() {
	p.Current++
	p.Render()
}

// Render displays the current progress bar state with ETA
func (p *ProgressBar) Render() {
	if p.Total == 0 {
		return
	}
	percent := float64(p.Current) / float64(p.Total)
	filled := int(percent * float64(p.Width))
	if filled > p.Width {
		filled = p.Width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)

	etaStr := ""
	processed := p.Current - p.StartOffset
	if !p.StartTime.IsZero() && processed > 0 && p.Current < p.Total {
		elapsed := time.Since(p.StartTime)
		avgPerRow := elapsed / time.Duration(processed)
		remaining := p.Total - p.Current
		eta := avgPerRow * time.Duration(remaining)
		etaStr = fmt.Sprintf(" | ETA: %s", formatDuration(eta))
	}

	fmt.Printf("\r  Processing: %s %.0f%% (%d/%d)%s   ", util.Green(bar), percent*100, p.Current, p.Total, etaStr)
	if p.Current >= p.Total {
		fmt.Println()
	}
}

// formatDuration formats a duration into a human-readable short string
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		m := int(d.Minutes())
		s := int(d.Seconds()) - m*60
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	h := int(d.Hours())
	m := int(d.Minutes()) - h*60
	return fmt.Sprintf("%dh%02dm", h, m)
}

// BulkStats holds statistics for bulk operations
type BulkStats struct {
	Total        int
	Found        int
	NotFound     int
	Errors       int
	ClientErrors int // 4xx errors (invalid input)
	ServerErrors int // 5xx errors (Tomba API)
	OutputFile   string
	Elapsed      time.Duration
	OpType       string
	ExtraCounts  map[string]int
}

// PrintStats displays the final bulk operation statistics
func (s *BulkStats) PrintStats() {
	hitRate := float64(0)
	if s.Total > 0 {
		hitRate = float64(s.Found) / float64(s.Total) * 100
	}
	fmt.Println()
	fmt.Printf("%s %d/%d emails found (%.1f%% hit rate)\n", util.SuccessIcon(), s.Found, s.Total, hitRate)
	if s.Errors > 0 {
		fmt.Printf("%s %d errors (%d client, %d server)\n", util.WarningIcon(), s.Errors, s.ClientErrors, s.ServerErrors)
	}
	if s.Elapsed > 0 {
		fmt.Printf("%s Completed in %s\n", util.InfoIcon(), util.Bold(s.Elapsed.Round(time.Millisecond).String()))
	}
	if s.OutputFile != "" {
		fmt.Printf("%s Results saved to %s\n", util.SuccessIcon(), util.Bold(s.OutputFile))
	}
}

// PrintStatsTable displays a formatted ASCII table with bulk operation statistics
func (s *BulkStats) PrintStatsTable() {
	hitRate := float64(0)
	if s.Total > 0 {
		hitRate = float64(s.Found) / float64(s.Total) * 100
	}

	type tableRow struct {
		label string
		value string
	}

	rows := []tableRow{
		{"Total", fmt.Sprintf("%d", s.Total)},
		{"Found", fmt.Sprintf("%d", s.Found)},
		{"Not Found", fmt.Sprintf("%d", s.NotFound)},
		{"Client Errors", fmt.Sprintf("%d", s.ClientErrors)},
		{"Server Errors", fmt.Sprintf("%d", s.ServerErrors)},
		{"Hit Rate", fmt.Sprintf("%.1f%%", hitRate)},
		{"Duration", s.Elapsed.Round(time.Millisecond).String()},
	}

	labelW := 14
	valueW := 12
	sep := strings.Repeat("─", labelW+2)
	sepV := strings.Repeat("─", valueW+2)

	var sb strings.Builder
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  ┌%s┬%s┐\n", sep, sepV)
	fmt.Fprintf(&sb, "  │ %-*s │ %-*s │\n", labelW, "Metric", valueW, "Value")
	fmt.Fprintf(&sb, "  ├%s┼%s┤\n", sep, sepV)

	for i, r := range rows {
		colorFn := util.White
		switch i {
		case 1:
			colorFn = util.Green
		case 2:
			colorFn = util.Yellow
		case 3, 4:
			colorFn = util.Red
		}
		padded := fmt.Sprintf("%*s", valueW, r.value)
		if i >= 1 && i <= 4 {
			padded = fmt.Sprintf("%*s", valueW, r.value)
			padded = colorFn(padded)
		}
		fmt.Fprintf(&sb, "  │ %-*s │ %s │\n", labelW, r.label, padded)
	}

	fmt.Fprintf(&sb, "  └%s┴%s┘\n", sep, sepV)
	fmt.Print(sb.String())

	if s.OutputFile != "" {
		fmt.Printf("  %s Results saved to %s\n\n", util.SuccessIcon(), util.Bold(s.OutputFile))
	}
}

// PrintDistributionChart renders a horizontal bar chart showing result distribution
func (s *BulkStats) PrintDistributionChart() {
	if s.Total == 0 {
		return
	}

	maxBarWidth := 30

	type chartRow struct {
		label string
		count int
		color func(string) string
	}

	fmt.Printf("\n  %s\n\n", util.Bold("Result Distribution"))

	mainRows := []chartRow{
		{"Found       ", s.Found, util.Green},
		{"NotFound    ", s.NotFound, util.Yellow},
		{"ClientErrors", s.ClientErrors, util.Red},
		{"ServerErrors", s.ServerErrors, util.Red},
	}

	for _, r := range mainRows {
		pct := float64(r.count) / float64(s.Total) * 100
		barLen := int(float64(r.count) / float64(s.Total) * float64(maxBarWidth))
		if r.count > 0 && barLen == 0 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen) + strings.Repeat("░", maxBarWidth-barLen)
		fmt.Printf("  %s %s %4d (%5.1f%%)\n", r.label, r.color(bar), r.count, pct)
	}

	if s.OpType == "verify" && len(s.ExtraCounts) > 0 {
		fmt.Printf("\n  %s\n\n", util.Bold("Verification Breakdown"))

		verifyRows := []chartRow{
			{"Deliverable  ", s.ExtraCounts["deliverable"], util.Green},
			{"Undeliverable", s.ExtraCounts["undeliverable"], util.Red},
			{"Risky        ", s.ExtraCounts["risky"], util.Yellow},
			{"Unknown      ", s.ExtraCounts["unknown"], util.Cyan},
		}

		for _, r := range verifyRows {
			pct := float64(r.count) / float64(s.Total) * 100
			barLen := int(float64(r.count) / float64(s.Total) * float64(maxBarWidth))
			if r.count > 0 && barLen == 0 {
				barLen = 1
			}
			bar := strings.Repeat("█", barLen) + strings.Repeat("░", maxBarWidth-barLen)
			fmt.Printf("  %s %s %4d (%5.1f%%)\n", r.label, r.color(bar), r.count, pct)
		}
	}

	fmt.Println()
}
