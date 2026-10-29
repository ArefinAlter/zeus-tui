package ui

import (
	"fmt"
	"strings"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	"github.com/charmbracelet/lipgloss"
)

// RenderDownloads renders the downloads view with parallel multi-segment progress.
func RenderDownloads(jobs []client.DownloadJob, selectedIndex int, width int) string {
	if len(jobs) == 0 {
		return CardStyle.Width(width - 4).Render(
			lipgloss.NewStyle().Foreground(ColorMuted).Render("No active downloads. Press Ctrl+K to ask agent to download a URL."),
		)
	}

	var sb strings.Builder
	for i, job := range jobs {
		isSelected := (i == selectedIndex)
		progress := 0.0
		if job.TotalBytes > 0 {
			progress = float64(job.CompletedBytes) / float64(job.TotalBytes)
		}

		// Status Badge
		var badge string
		switch job.Status {
		case client.StatusRunning:
			badge = BadgeInfo.Render("DOWNLOADING")
		case client.StatusPaused:
			badge = BadgeWarning.Render("PAUSED")
		case client.StatusCompleted:
			badge = BadgeSuccess.Render("VERIFIED")
		default:
			badge = BadgeMuted.Render(string(job.Status))
		}

		header := fmt.Sprintf("%s  %s  %s",
			badge,
			lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render(job.Filename),
			lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("[%s / %s]", FormatBytes(job.CompletedBytes), FormatBytes(job.TotalBytes))),
		)

		stats := fmt.Sprintf("Speed: %s  •  ETA: %ds  •  Hash: %s (%s)",
			lipgloss.NewStyle().Foreground(ColorSuccess).Render(FormatSpeed(job.SpeedBps)),
			job.ETASeconds,
			job.HashAlgorithm,
			job.Destination,
		)

		barWidth := width - 10
		if barWidth < 20 {
			barWidth = 20
		}

		segmentsCount := len(job.Segments)
		if segmentsCount == 0 {
			segmentsCount = 4
		}
		progressBar := RenderMultiSegmentProgress(barWidth, progress, segmentsCount)

		cardContent := fmt.Sprintf("%s\n%s\n%s (%.1f%%)",
			header,
			progressBar,
			stats,
			progress*100,
		)

		// If selected, show parallel chunk worker breakdown
		if isSelected && len(job.Segments) > 0 {
			var segSb strings.Builder
			segSb.WriteString("\n" + lipgloss.NewStyle().Foreground(ColorSecondary).Render("── Parallel Range Chunks (Work-Stealing Adaptive Transfer) ──") + "\n")
			for _, seg := range job.Segments {
				segProgress := 0.0
				segTotal := seg.EndByte - seg.StartByte
				if segTotal > 0 {
					segProgress = float64(seg.CompletedBytes) / float64(segTotal)
				}
				segBar := RenderMultiSegmentProgress(24, segProgress, 1)
				segSb.WriteString(fmt.Sprintf("  Chunk #%d [Worker %d]: %s %s [%s/%s]\n",
					seg.Index+1,
					seg.WorkerID,
					segBar,
					FormatSpeed(seg.SpeedBps),
					FormatBytes(seg.CompletedBytes),
					FormatBytes(segTotal),
				))
			}
			cardContent += segSb.String()
		}

		style := CardStyle
		if isSelected {
			style = SelectedCardStyle
		}
		sb.WriteString(style.Width(width - 4).Render(cardContent))
		sb.WriteString("\n")
	}

	return sb.String()
}
