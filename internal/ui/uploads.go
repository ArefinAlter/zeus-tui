package ui

import (
	"fmt"
	"strings"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	"github.com/charmbracelet/lipgloss"
)

// RenderUploads renders the staged upload plans view.
func RenderUploads(plans []client.UploadPlan, selectedIndex int, width int) string {
	if len(plans) == 0 {
		return CardStyle.Width(width - 4).Render(
			lipgloss.NewStyle().Foreground(ColorMuted).Render("No active upload plans. Press Ctrl+K to ask agent to stage an upload to Cloud/S3/GDrive."),
		)
	}

	var sb strings.Builder
	for i, plan := range plans {
		isSelected := (i == selectedIndex)
		progress := 0.0
		if plan.TotalBytes > 0 {
			progress = float64(plan.SentBytes) / float64(plan.TotalBytes)
		}

		// Status Badge
		var badge string
		switch {
		case !plan.Reviewed:
			badge = BadgeWarning.Render("REVIEW REQUIRED")
		case plan.Status == client.StatusRunning:
			badge = BadgeInfo.Render("UPLOADING")
		case plan.Status == client.StatusCompleted:
			badge = BadgeSuccess.Render("VERIFIED COMPLETE")
		case plan.Status == client.StatusQueued:
			badge = BadgeSuccess.Render("STAGED & APPROVED")
		default:
			badge = BadgeMuted.Render(string(plan.Status))
		}

		header := fmt.Sprintf("%s  %s  →  %s",
			badge,
			lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render(plan.SourcePath),
			lipgloss.NewStyle().Foreground(ColorSecondary).Render(plan.Destination),
		)

		stats := fmt.Sprintf("Progress: %s / %s  •  Speed: %s  •  Source Hash: %s",
			FormatBytes(plan.SentBytes),
			FormatBytes(plan.TotalBytes),
			lipgloss.NewStyle().Foreground(ColorSuccess).Render(FormatSpeed(plan.SpeedBps)),
			lipgloss.NewStyle().Foreground(ColorMuted).Render(plan.SourceHash),
		)

		barWidth := width - 10
		if barWidth < 20 {
			barWidth = 20
		}
		progressBar := RenderMultiSegmentProgress(barWidth, progress, 4)

		cardContent := fmt.Sprintf("%s\n%s\n%s (%.1f%%)",
			header,
			progressBar,
			stats,
			progress*100,
		)

		if isSelected {
			var detailSb strings.Builder
			detailSb.WriteString("\n" + lipgloss.NewStyle().Foreground(ColorSecondary).Render("── Staged Upload Inspection & Verification Receipt ──") + "\n")
			detailSb.WriteString(fmt.Sprintf("  Plan ID:        %s\n", plan.ID))
			detailSb.WriteString(fmt.Sprintf("  Source File:    %s\n", plan.SourcePath))
			detailSb.WriteString(fmt.Sprintf("  Destination:    %s\n", plan.Destination))
			detailSb.WriteString(fmt.Sprintf("  Content Digest: %s (Verified in staging)\n", plan.SourceHash))
			detailSb.WriteString(fmt.Sprintf("  Byte Limit:     Strictly bounded to %s (No link-following / symlink traversal)\n", FormatBytes(plan.TotalBytes)))
			if !plan.Reviewed {
				detailSb.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render("  [!] Human Approval Required: Press [Space] or [Enter] to approve transfer.\n"))
			} else {
				detailSb.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Render("  [✓] Plan approved and cryptographically locked against source mutation.\n"))
			}
			cardContent += detailSb.String()
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
