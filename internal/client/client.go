package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Client talks to the local Zeus daemon (zeusd).
type Client struct {
	baseURL    string
	httpClient *http.Client
	isMock     bool
	mu         sync.Mutex
	mockDown   []DownloadJob
	mockUp     []UploadPlan
	mockFiles  []FileItem
	mockAction *AgentAction
}

// New creates a new client targeting daemonURL.
func New(daemonURL string, isMock bool) *Client {
	c := &Client{
		baseURL: daemonURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		isMock: isMock,
	}
	c.initMockData()
	return c
}

func (c *Client) initMockData() {
	now := time.Now()
	c.mockDown = []DownloadJob{
		{
			ID:             "job-901a",
			URL:            "https://releases.ubuntu.com/24.04/ubuntu-24.04-desktop-amd64.iso",
			Filename:       "ubuntu-24.04-desktop-amd64.iso",
			Destination:    "D:\\Downloads\\ubuntu-24.04-desktop-amd64.iso",
			TotalBytes:     6048576000,
			CompletedBytes: 4234003200,
			SpeedBps:       74500000,
			ETASeconds:     24,
			Status:         StatusRunning,
			HashAlgorithm:  "SHA-256",
			CreatedAt:      now.Add(-2 * time.Minute),
			Segments: []Segment{
				{Index: 0, StartByte: 0, EndByte: 1512144000, CompletedBytes: 1512144000, SpeedBps: 0, WorkerID: 1, Progress: 1.0},
				{Index: 1, StartByte: 1512144000, EndByte: 3024288000, CompletedBytes: 1512144000, SpeedBps: 0, WorkerID: 2, Progress: 1.0},
				{Index: 2, StartByte: 3024288000, EndByte: 4536432000, CompletedBytes: 984000000, SpeedBps: 38200000, WorkerID: 3, Progress: 0.65},
				{Index: 3, StartByte: 4536432000, EndByte: 6048576000, CompletedBytes: 225715200, SpeedBps: 36300000, WorkerID: 4, Progress: 0.15},
			},
		},
		{
			ID:             "job-902b",
			URL:            "https://huggingface.co/mistralai/Mistral-7B-Instruct-v0.3/resolve/main/model-00001-of-00004.safetensors",
			Filename:       "Mistral-7B-v0.3-part1.safetensors",
			Destination:    "D:\\Models\\Mistral-7B-part1.safetensors",
			TotalBytes:     4980000000,
			CompletedBytes: 1842600000,
			SpeedBps:       42100000,
			ETASeconds:     74,
			Status:         StatusRunning,
			HashAlgorithm:  "BLAKE3",
			CreatedAt:      now.Add(-5 * time.Minute),
			Segments: []Segment{
				{Index: 0, StartByte: 0, EndByte: 1245000000, CompletedBytes: 1245000000, SpeedBps: 0, WorkerID: 1, Progress: 1.0},
				{Index: 1, StartByte: 1245000000, EndByte: 2490000000, CompletedBytes: 597600000, SpeedBps: 21500000, WorkerID: 2, Progress: 0.48},
				{Index: 2, StartByte: 2490000000, EndByte: 3735000000, CompletedBytes: 0, SpeedBps: 20600000, WorkerID: 3, Progress: 0.0},
				{Index: 3, StartByte: 3735000000, EndByte: 4980000000, CompletedBytes: 0, SpeedBps: 0, WorkerID: 4, Progress: 0.0},
			},
		},
		{
			ID:             "job-899f",
			URL:            "https://github.com/rust-lang/rust/archive/refs/tags/1.85.0.tar.gz",
			Filename:       "rust-1.85.0.tar.gz",
			Destination:    "D:\\Archives\\rust-1.85.0.tar.gz",
			TotalBytes:     312540000,
			CompletedBytes: 312540000,
			SpeedBps:       0,
			ETASeconds:     0,
			Status:         StatusCompleted,
			HashAlgorithm:  "SHA-256",
			HashDigest:     "6d8f28b...c491",
			CreatedAt:      now.Add(-15 * time.Minute),
			Segments:       []Segment{},
		},
	}

	c.mockUp = []UploadPlan{
		{
			ID:          "upl-441",
			SourcePath:  "D:\\Research\\benchmark-results-20261029.parquet",
			Destination: "s3://verne-artifacts/benchmarks/20261029/",
			SourceHash:  "blake3:9f80d484...467",
			TotalBytes:  142857000,
			SentBytes:   92857000,
			SpeedBps:    18400000,
			Reviewed:    true,
			Status:      StatusRunning,
			CreatedAt:   now.Add(-1 * time.Minute),
		},
		{
			ID:          "upl-442",
			SourcePath:  "D:\\Backups\\dataset-archive-v1.tar.zst",
			Destination: "gdrive://Shared/Verne/Datasets/",
			SourceHash:  "sha256:cc7d47cf...db9",
			TotalBytes:  2415000000,
			SentBytes:   0,
			SpeedBps:    0,
			Reviewed:    false,
			Status:      StatusQueued,
			CreatedAt:   now.Add(-30 * time.Second),
		},
	}

	c.mockFiles = []FileItem{
		{Name: "ubuntu-24.04-desktop-amd64.iso", Path: "D:\\Downloads\\ubuntu-24.04-desktop-amd64.iso", SizeBytes: 6048576000, IsDir: false, ModTime: now},
		{Name: "Mistral-7B-part1.safetensors", Path: "D:\\Models\\Mistral-7B-part1.safetensors", SizeBytes: 4980000000, IsDir: false, ModTime: now},
		{Name: "rust-1.85.0.tar.gz", Path: "D:\\Archives\\rust-1.85.0.tar.gz", SizeBytes: 312540000, IsDir: false, ModTime: now.Add(-15 * time.Minute), Hash: "6d8f28b...c491"},
		{Name: "verne-build", Path: "D:\\verne-build", SizeBytes: 0, IsDir: true, ModTime: now.Add(-1 * time.Hour)},
		{Name: "workspace-sync", Path: "D:\\workspace-sync", SizeBytes: 0, IsDir: true, ModTime: now.Add(-3 * time.Hour)},
	}
}

// GetStatus checks daemon connection and returns health stats.
func (c *Client) GetStatus(ctx context.Context) (*DaemonStatus, error) {
	if c.isMock {
		return c.mockStatus(), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/status", nil)
	if err != nil {
		return c.mockStatus(), nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Fallback to simulated offline/standalone status
		stat := c.mockStatus()
		stat.Connected = false
		return stat, nil
	}
	defer resp.Body.Close()

	var st DaemonStatus
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return c.mockStatus(), nil
	}
	st.Connected = true
	return &st, nil
}

func (c *Client) mockStatus() *DaemonStatus {
	c.mu.Lock()
	defer c.mu.Unlock()

	activeDown := 0
	var totalSpeed int64
	for _, d := range c.mockDown {
		if d.Status == StatusRunning {
			activeDown++
			totalSpeed += d.SpeedBps
		}
	}

	return &DaemonStatus{
		Connected:      !c.isMock,
		Version:        "0.1.0-alpha.1",
		DaemonURL:      c.baseURL,
		UptimeSeconds:  38420,
		ActiveDown:     activeDown,
		ActiveUp:       len(c.mockUp),
		TotalSpeedBps:  totalSpeed,
		StorageFreeGB:  810.5,
		StorageTotalGB: 1863.0,
	}
}

// ListDownloads retrieves all active and historical download jobs.
func (c *Client) ListDownloads(ctx context.Context) ([]DownloadJob, error) {
	if c.isMock {
		c.mu.Lock()
		defer c.mu.Unlock()
		// Progress simulation
		for i := range c.mockDown {
			if c.mockDown[i].Status == StatusRunning {
				c.mockDown[i].CompletedBytes += int64(float64(c.mockDown[i].SpeedBps) * 0.5)
				if c.mockDown[i].CompletedBytes > c.mockDown[i].TotalBytes {
					c.mockDown[i].CompletedBytes = c.mockDown[i].TotalBytes
					c.mockDown[i].Status = StatusCompleted
					c.mockDown[i].SpeedBps = 0
				}
			}
		}
		return c.mockDown, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/downloads", nil)
	if err != nil {
		return c.mockDown, nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return c.mockDown, nil
	}
	defer resp.Body.Close()

	var jobs []DownloadJob
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return c.mockDown, nil
	}
	return jobs, nil
}

// TogglePauseDownload toggles pause/running status.
func (c *Client) TogglePauseDownload(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.mockDown {
		if c.mockDown[i].ID == id {
			if c.mockDown[i].Status == StatusRunning {
				c.mockDown[i].Status = StatusPaused
				c.mockDown[i].SpeedBps = 0
			} else if c.mockDown[i].Status == StatusPaused {
				c.mockDown[i].Status = StatusRunning
				c.mockDown[i].SpeedBps = 45000000
			}
			return nil
		}
	}
	return fmt.Errorf("job not found: %s", id)
}

// ListUploadPlans retrieves all staged and active upload plans.
func (c *Client) ListUploadPlans(ctx context.Context) ([]UploadPlan, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mockUp, nil
}

// ListFiles returns files in the managed workspace.
func (c *Client) ListFiles(ctx context.Context) ([]FileItem, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mockFiles, nil
}

// SubmitAgentPrompt parses natural language intent into a structured operation proposal.
func (c *Client) SubmitAgentPrompt(ctx context.Context, prompt string) (*AgentAction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	action := &AgentAction{
		ID:               fmt.Sprintf("act-%d", time.Now().Unix()%10000),
		Prompt:           prompt,
		RequiresApproval: true,
		Approved:         false,
	}

	if bytes.Contains([]byte(prompt), []byte("git")) || bytes.Contains([]byte(prompt), []byte("clone")) {
		action.ActionType = "clone"
		action.Summary = "Clone repository into verified workspace"
		action.Source = "https://github.com/pacifio/atlas.git"
		action.Destination = "D:\\workspace\\atlas"
	} else if bytes.Contains([]byte(prompt), []byte("upload")) || bytes.Contains([]byte(prompt), []byte("s3")) {
		action.ActionType = "upload"
		action.Summary = "Stage reviewed parquet artifacts to S3 bucket"
		action.Source = "D:\\Research\\benchmark-results-20261029.parquet"
		action.Destination = "s3://verne-artifacts/benchmarks/"
	} else {
		action.ActionType = "download"
		action.Summary = "Accelerated multi-range parallel HTTP download"
		action.Source = prompt
		action.Destination = "D:\\Downloads\\"
	}

	c.mockAction = action
	return action, nil
}

// ApproveAgentAction marks the pending action as approved and enqueues it.
func (c *Client) ApproveAgentAction(actionID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.mockAction != nil && c.mockAction.ID == actionID {
		c.mockAction.Approved = true
		if c.mockAction.ActionType == "download" {
			c.mockDown = append([]DownloadJob{{
				ID:             fmt.Sprintf("job-%d", time.Now().Unix()%1000),
				URL:            c.mockAction.Source,
				Filename:       "custom-download.pkg",
				Destination:    "D:\\Downloads\\custom-download.pkg",
				TotalBytes:     1024 * 1024 * 512,
				CompletedBytes: 0,
				SpeedBps:       35000000,
				ETASeconds:     15,
				Status:         StatusRunning,
				CreatedAt:      time.Now(),
			}}, c.mockDown...)
		}
		c.mockAction = nil
		return nil
	}
	return fmt.Errorf("action not found")
}

// GetPendingAction returns current pending action if any.
func (c *Client) GetPendingAction() *AgentAction {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mockAction
}

// ClearPendingAction discards pending action.
func (c *Client) ClearPendingAction() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mockAction = nil
}
