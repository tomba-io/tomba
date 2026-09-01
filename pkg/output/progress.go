package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/tomba-io/tomba/pkg/util"
)

// ProgressBar renders a progress bar to the terminal
type ProgressBar struct {
	Total   int
	Current int
	Width   int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total:   total,
		Current: 0,
		Width:   30,
	}
}

// Increment advances the progress bar by one
func (p *ProgressBar) Increment() {
	p.Current++
	p.Render()
}

// Render displays the current progress bar state
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
	fmt.Printf("\r  Processing: %s %.0f%% (%d/%d)", util.Green(bar), percent*100, p.Current, p.Total)
	if p.Current >= p.Total {
		fmt.Println()
	}
}

// BulkStats holds statistics for bulk operations
type BulkStats struct {
	Total      int
	Found      int
	NotFound   int
	Errors     int
	OutputFile string
	Elapsed    time.Duration
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
		fmt.Printf("%s %d errors encountered\n", util.WarningIcon(), s.Errors)
	}
	if s.Elapsed > 0 {
		fmt.Printf("%s Completed in %s\n", util.InfoIcon(), util.Bold(s.Elapsed.Round(time.Millisecond).String()))
	}
	if s.OutputFile != "" {
		fmt.Printf("%s Results saved to %s\n", util.SuccessIcon(), util.Bold(s.OutputFile))
	}
}
