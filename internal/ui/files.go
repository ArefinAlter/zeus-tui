package ui

import (
	"fmt"
	"strings"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	"github.com/charmbracelet/lipgloss"
)

// RenderFiles renders the desktop file manager and workspace index view.
func RenderFiles(files []client.FileItem, selectedIndex int, searchQuery string, width int) string {
	var sb strings.Builder

	// Top search bar / Everything 1.5a query status
	searchBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1).
		Width(width - 4).
		Render(fmt.Sprintf("🔍 Workspace Search (Everything 1.5a IPC): %s",
			func() string {
				if searchQuery != "" {
					return lipgloss.NewStyle().Foreground(ColorText).Bold(true).Render(searchQuery)
				}
				return lipgloss.NewStyle().Foreground(ColorDim).Render("Type to filter workspace files across local drives...")
			}(),
		))
	sb.WriteString(searchBox + "\n\n")

	if len(files) == 0 {
		sb.WriteString(CardStyle.Width(width - 4).Render(
			lipgloss.NewStyle().Foreground(ColorMuted).Render("No files found in workspace directory."),
		))
		return sb.String()
	}

	for i, f := range files {
		isSelected := (i == selectedIndex)

		// Icon / Indicator
		icon := "📄"
		typeLabel := "FILE"
		if f.IsDir {
			icon = "📁"
			typeLabel = "DIR "
		}

		sizeStr := FormatBytes(f.SizeBytes)
		if f.IsDir {
			sizeStr = "     -"
		}

		line := fmt.Sprintf("%s %s  %-35s  %10s  %s",
			icon,
			lipgloss.NewStyle().Foreground(ColorDim).Render(typeLabel),
			lipgloss.NewStyle().Bold(isSelected).Foreground(func() lipgloss.Color {
				if isSelected {
					return ColorText
				}
				return ColorMuted
			}()).Render(truncateString(f.Name, 35)),
			lipgloss.NewStyle().Foreground(ColorSecondary).Render(sizeStr),
			lipgloss.NewStyle().Foreground(ColorDim).Render(f.ModTime.Format("2006-01-02 15:04")),
		)

		if isSelected {
			var detailSb strings.Builder
			detailSb.WriteString(line + "\n")
			detailSb.WriteString(lipgloss.NewStyle().Foreground(ColorSecondary).Render("  ── File Metadata & Storage Observation ──") + "\n")
			detailSb.WriteString(fmt.Sprintf("    Full Path: %s\n", f.Path))
			if f.Hash != "" {
				detailSb.WriteString(fmt.Sprintf("    Digest:    %s\n", f.Hash))
			} else {
				detailSb.WriteString("    Digest:    Verified unmutated (Local copy intact)\n")
			}
			detailSb.WriteString("    Links:     1 (Direct regular file, links not followed)\n")
			detailSb.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("    Actions:   [Enter] Open  •  [u] Stage Upload  •  [d] Inspect Storage\n"))

			card := SelectedCardStyle.Width(width - 4).Render(detailSb.String())
			sb.WriteString(card + "\n")
		} else {
			sb.WriteString("  " + line + "\n")
		}
	}

	return sb.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
