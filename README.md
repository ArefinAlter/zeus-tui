# zeus-tui

> **Terminal Client for Zeus Transfer & Storage Engine**  
> Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Part of the **Verne Project Suite** (`ArefinAlter/zeus`).

---

## Overview

`zeus-tui` provides a terminal user interface for monitoring and interacting with `zeusd` (the Zeus desktop transfer and storage daemon):

- **Parallel Range Visualizer**: Live multi-segment visual progress across 8–16 adaptive range workers with per-worker speed and chunk boundaries.
- **Desktop File Manager & Workspace Index**: Fast inspection of files, digests, link counts, and storage observations.
- **Staged Upload Management**: Human review workflow for local-file uploads to S3, Google Drive, and cloud destinations, enforcing bounded byte limits and verified staging receipts.
- **Natural Language Agent Bar (`Ctrl+K`)**: Prompt bar proposing transfer and organization actions via ACP subprocess adapters, requiring explicit human approval before execution.

---

## Keybindings

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Cycle between tabs (Downloads, Upload Plans, Files) |
| `1`, `2`, `3` | Directly switch to view 1, 2, or 3 |
| `j` / `k` (or `↓` / `↑`) | Navigate items in current list |
| `Space` | Toggle pause/resume (downloads) or approve (upload/agent plan) |
| `Enter` | Confirm and dispatch approved agent action |
| `Ctrl+K` or `/` | Focus Natural Language Agent prompt bar |
| `Esc` | Cancel / clear prompt or pending action |
| `q` or `Ctrl+C` | Quit |

---

## Building and Running

### Prerequisites
- Go 1.24+

### Run Standalone Simulation Mode
```powershell
go run ./cmd/zeus-tui --demo
```

### Connect to Running Zeus Daemon
```powershell
go run ./cmd/zeus-tui --daemon-url "http://127.0.0.1:18540"
```

### Compile Binary
```powershell
go build -o zeus-tui.exe ./cmd/zeus-tui
```

---

## License
MIT / Apache-2.0 (Verne Project Suite)
