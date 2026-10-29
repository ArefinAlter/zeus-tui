package ui

import (
	"fmt"
	"strings"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	"github.com/charmbracelet/lipgloss"
)

// RenderAgentBar renders the natural language prompt overlay or bar at bottom/top.
func RenderAgentBar(active bool, promptText string, width int) string {
	boxWidth := width - 4
	if boxWidth < 20 {
		boxWidth = 20
	}

	if !active {
		return lipgloss.NewStyle().
			Foreground(ColorDim).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Width(boxWidth).
			Render("  💬 Agent Prompt (Press Ctrl+K or / to type natural language command...)")
	}

	promptDisplay := promptText
	if promptDisplay == "" {
		promptDisplay = lipgloss.NewStyle().Foreground(ColorDim).Render("e.g. 'download https://... to D:\\Models' or 'stage upload of benchmarks to S3'")
	} else {
		promptDisplay = lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render(promptDisplay) + " █"
	}

	content := fmt.Sprintf("  🤖 %s %s\n  %s",
		lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render("Zeus Natural Language Agent"),
		lipgloss.NewStyle().Foreground(ColorMuted).Render("(ACP Subprocess Dispatcher)"),
		promptDisplay,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Background(ColorCardBg).
		Width(boxWidth).
		Render(content)
}

// RenderPendingAction renders the proposed action review card for human approval.
func RenderPendingAction(action *client.AgentAction, width int) string {
	if action == nil {
		return ""
	}

	boxWidth := width - 4
	if boxWidth < 20 {
		boxWidth = 20
	}

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorWarning).Render("⚡ Proposed Agent Action (Requires Human Review)") + "\n\n")
	sb.WriteString(fmt.Sprintf("  Action:      %s\n", strings.ToUpper(action.ActionType)))
	sb.WriteString(fmt.Sprintf("  Summary:     %s\n", action.Summary))
	if action.Source != "" {
		sb.WriteString(fmt.Sprintf("  Source:      %s\n", action.Source))
	}
	if action.Destination != "" {
		sb.WriteString(fmt.Sprintf("  Destination: %s\n", action.Destination))
	}

	sb.WriteString("\n")
	if action.RequiresApproval && !action.Approved {
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorWarning).Render("  Press [Enter] to Approve and Execute  •  [Esc] to Cancel"))
	} else {
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Render("  ✓ Action Approved — Dispatched to zeusd engine"))
	}

	return CardStyle.
		BorderForeground(ColorWarning).
		Width(boxWidth).
		Render(sb.String())
}
