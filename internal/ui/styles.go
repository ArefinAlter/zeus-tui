package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors matching the Verne / Zeus desktop theme
	ColorPrimary   = lipgloss.Color("#6366F1") // Indigo
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColorDarkBg    = lipgloss.Color("#0F172A") // Slate 900
	ColorCardBg    = lipgloss.Color("#1E293B") // Slate 800
	ColorBorder    = lipgloss.Color("#334155") // Slate 700
	ColorText      = lipgloss.Color("#F8FAFC") // Slate 50
	ColorMuted     = lipgloss.Color("#94A3B8") // Slate 400
	ColorDim       = lipgloss.Color("#475569") // Slate 600
	ColorSuccess   = lipgloss.Color("#10B981") // Emerald 500
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber 500
	ColorError     = lipgloss.Color("#EF4444") // Rose 500

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText).
			Background(ColorPrimary).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText).
			Background(ColorPrimary).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Background(ColorBorder).
				Padding(0, 2)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorCardBg).
			Padding(1)

	SelectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Background(ColorCardBg).
				Padding(1)

	BadgeSuccess = lipgloss.NewStyle().
			Foreground(ColorDarkBg).
			Background(ColorSuccess).
			Bold(true).
			Padding(0, 1)

	BadgeWarning = lipgloss.NewStyle().
			Foreground(ColorDarkBg).
			Background(ColorWarning).
			Bold(true).
			Padding(0, 1)

	BadgeInfo = lipgloss.NewStyle().
			Foreground(ColorDarkBg).
			Background(ColorSecondary).
			Bold(true).
			Padding(0, 1)

	BadgeMuted = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorDim).
			Padding(0, 1)

	AgentBarFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1)
)

// FormatBytes formats byte counts cleanly (e.g. 1.45 GB).
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatSpeed formats transfer speeds (e.g. 74.50 MB/s).
func FormatSpeed(bps int64) string {
	if bps <= 0 {
		return "0 B/s"
	}
	return FormatBytes(bps) + "/s"
}

// RenderMultiSegmentProgress renders an 8-16 parallel range visualizer across width columns.
func RenderMultiSegmentProgress(width int, progress float64, segmentCount int) string {
	if width < 10 {
		width = 10
	}
	if segmentCount <= 0 {
		segmentCount = 4
	}

	totalFilled := int(progress * float64(width))
	if totalFilled > width {
		totalFilled = width
	}

	var sb strings.Builder
	for i := 0; i < width; i++ {
		isSegmentBoundary := (i > 0) && (i%(width/segmentCount) == 0)
		if isSegmentBoundary {
			sb.WriteString(lipgloss.NewStyle().Foreground(ColorBorder).Render("│"))
			continue
		}

		if i < totalFilled {
			sb.WriteString(lipgloss.NewStyle().Foreground(ColorSecondary).Render("█"))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("░"))
		}
	}
	return sb.String()
}
