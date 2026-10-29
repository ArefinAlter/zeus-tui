package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	TabDownloads Tab = iota
	TabUploads
	TabFiles
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model is the main Bubble Tea model for zeus-tui.
type Model struct {
	client        *client.Client
	currentTab    Tab
	width         int
	height        int
	status        *client.DaemonStatus
	downloads     []client.DownloadJob
	uploads       []client.UploadPlan
	files         []client.FileItem
	pendingAction *client.AgentAction

	selectedDown   int
	selectedUp     int
	selectedFile   int
	agentActive    bool
	agentPrompt    string
	searchQuery    string
	statusMessage  string
	statusMsgTimer int
}

// NewModel instantiates the TUI model.
func NewModel(c *client.Client) Model {
	return Model{
		client:      c,
		currentTab:  TabDownloads,
		width:       100,
		height:      30,
		agentActive: false,
	}
}

// Init sets up initial commands (refresh & tick).
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.refreshDataCmd(),
		tick(),
	)
}

func (m Model) refreshDataCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		status, _ := m.client.GetStatus(ctx)
		downloads, _ := m.client.ListDownloads(ctx)
		uploads, _ := m.client.ListUploadPlans(ctx)
		files, _ := m.client.ListFiles(ctx)
		pending := m.client.GetPendingAction()

		return dataRefreshMsg{
			status:    status,
			downloads: downloads,
			uploads:   uploads,
			files:     files,
			pending:   pending,
		}
	}
}

type dataRefreshMsg struct {
	status    *client.DaemonStatus
	downloads []client.DownloadJob
	uploads   []client.UploadPlan
	files     []client.FileItem
	pending   *client.AgentAction
}

// Update handles incoming messages and keypresses.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.statusMsgTimer > 0 {
			m.statusMsgTimer--
			if m.statusMsgTimer == 0 {
				m.statusMessage = ""
			}
		}
		return m, tea.Batch(m.refreshDataCmd(), tick())

	case dataRefreshMsg:
		m.status = msg.status
		m.downloads = msg.downloads
		m.uploads = msg.uploads
		m.files = msg.files
		m.pendingAction = msg.pending
		return m, nil

	case tea.KeyMsg:
		// If Agent input prompt is focused
		if m.agentActive {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEsc:
				m.agentActive = false
				m.agentPrompt = ""
				return m, nil
			case tea.KeyEnter:
				if strings.TrimSpace(m.agentPrompt) != "" {
					act, err := m.client.SubmitAgentPrompt(context.Background(), m.agentPrompt)
					if err == nil {
						m.pendingAction = act
						m.statusMessage = "Agent proposal created. Review below."
					} else {
						m.statusMessage = "Agent error: " + err.Error()
					}
					m.agentActive = false
					m.agentPrompt = ""
					m.statusMsgTimer = 6
				}
				return m, nil
			case tea.KeyBackspace:
				if len(m.agentPrompt) > 0 {
					m.agentPrompt = m.agentPrompt[:len(m.agentPrompt)-1]
				}
				return m, nil
			case tea.KeyRunes, tea.KeySpace:
				m.agentPrompt += msg.String()
				return m, nil
			default:
				return m, nil
			}
		}

		// Standard Navigation Mode
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.currentTab = (m.currentTab + 1) % 3
			return m, nil

		case "shift+tab":
			m.currentTab = (m.currentTab + 2) % 3
			return m, nil

		case "1":
			m.currentTab = TabDownloads
			return m, nil
		case "2":
			m.currentTab = TabUploads
			return m, nil
		case "3":
			m.currentTab = TabFiles
			return m, nil

		case "ctrl+k", "/":
			m.agentActive = true
			return m, nil

		case "esc":
			if m.pendingAction != nil {
				m.client.ClearPendingAction()
				m.pendingAction = nil
				m.statusMessage = "Action proposal cleared."
				m.statusMsgTimer = 4
			}
			return m, nil

		case "up", "k":
			switch m.currentTab {
			case TabDownloads:
				if m.selectedDown > 0 {
					m.selectedDown--
				}
			case TabUploads:
				if m.selectedUp > 0 {
					m.selectedUp--
				}
			case TabFiles:
				if m.selectedFile > 0 {
					m.selectedFile--
				}
			}
			return m, nil

		case "down", "j":
			switch m.currentTab {
			case TabDownloads:
				if m.selectedDown < len(m.downloads)-1 {
					m.selectedDown++
				}
			case TabUploads:
				if m.selectedUp < len(m.uploads)-1 {
					m.selectedUp++
				}
			case TabFiles:
				if m.selectedFile < len(m.files)-1 {
					m.selectedFile++
				}
			}
			return m, nil

		case "enter":
			if m.pendingAction != nil {
				_ = m.client.ApproveAgentAction(m.pendingAction.ID)
				m.pendingAction = nil
				m.statusMessage = "Action approved and dispatched to daemon!"
				m.statusMsgTimer = 6
				return m, m.refreshDataCmd()
			}
			return m, nil

		case " ":
			switch m.currentTab {
			case TabDownloads:
				if len(m.downloads) > 0 && m.selectedDown < len(m.downloads) {
					job := m.downloads[m.selectedDown]
					_ = m.client.TogglePauseDownload(job.ID)
					return m, m.refreshDataCmd()
				}
			case TabUploads:
				if len(m.uploads) > 0 && m.selectedUp < len(m.uploads) {
					m.statusMessage = fmt.Sprintf("Upload %s approved.", m.uploads[m.selectedUp].ID)
					m.statusMsgTimer = 5
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// View renders the terminal screen.
func (m Model) View() string {
	var sb strings.Builder

	// Header banner
	title := TitleStyle.Render(" ZEUS ")
	subTitle := lipgloss.NewStyle().Foreground(ColorMuted).Render(" Universal Transfer & Storage Engine  •  v0.1.0-alpha.1")

	// Connection status indicator
	connBadge := BadgeSuccess.Render("● DAEMON CONNECTED")
	speedSummary := ""
	if m.status != nil {
		if !m.status.Connected {
			connBadge = BadgeWarning.Render("○ STANDALONE / SIMULATED")
		}
		speedSummary = fmt.Sprintf("  •  Total: %s  •  Storage: %.1f GB free",
			lipgloss.NewStyle().Foreground(ColorSuccess).Render(FormatSpeed(m.status.TotalSpeedBps)),
			m.status.StorageFreeGB,
		)
	}

	header := fmt.Sprintf("%s%s  [%s%s]\n", title, subTitle, connBadge, speedSummary)
	sb.WriteString(header + "\n")

	// Navigation Tabs
	var tabs []string
	labels := []string{"1. Downloads", "2. Upload Plans", "3. File Workspace"}
	for i, label := range labels {
		if Tab(i) == m.currentTab {
			tabs = append(tabs, TabActiveStyle.Render(label))
		} else {
			tabs = append(tabs, TabInactiveStyle.Render(label))
		}
	}
	sb.WriteString(strings.Join(tabs, " ") + "\n\n")

	// Main Tab Content
	switch m.currentTab {
	case TabDownloads:
		sb.WriteString(RenderDownloads(m.downloads, m.selectedDown, m.width))
	case TabUploads:
		sb.WriteString(RenderUploads(m.uploads, m.selectedUp, m.width))
	case TabFiles:
		sb.WriteString(RenderFiles(m.files, m.selectedFile, m.searchQuery, m.width))
	}

	// Pending action box (if awaiting human approval)
	if m.pendingAction != nil {
		sb.WriteString("\n" + RenderPendingAction(m.pendingAction, m.width) + "\n")
	}

	// Status toast message
	if m.statusMessage != "" {
		sb.WriteString("\n" + lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("▶ "+m.statusMessage) + "\n")
	}

	// Bottom Agent Bar
	sb.WriteString("\n" + RenderAgentBar(m.agentActive, m.agentPrompt, m.width) + "\n")

	// Bottom Keybindings Help Bar
	help := lipgloss.NewStyle().Foreground(ColorDim).Render(
		"[Tab/1-3] Switch View  •  [j/k] Navigate  •  [Space] Pause/Approve  •  [Ctrl+K] Agent Bar  •  [q] Quit",
	)
	sb.WriteString(help)

	return sb.String()
}
