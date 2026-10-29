package client

import "time"

// Status represents the download state.
type Status string

const (
	StatusQueued    Status = "QUEUED"
	StatusRunning   Status = "RUNNING"
	StatusPaused    Status = "PAUSED"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

// Segment represents a parallel HTTP range slice.
type Segment struct {
	Index          int     `json:"index"`
	StartByte      int64   `json:"start_byte"`
	EndByte        int64   `json:"end_byte"`
	CompletedBytes int64   `json:"completed_bytes"`
	SpeedBps       int64   `json:"speed_bps"`
	WorkerID       int     `json:"worker_id"`
	Progress       float64 `json:"progress"`
}

// DownloadJob represents an active or queued download job in zeusd.
type DownloadJob struct {
	ID             string    `json:"id"`
	URL            string    `json:"url"`
	Filename       string    `json:"filename"`
	Destination    string    `json:"destination"`
	TotalBytes     int64     `json:"total_bytes"`
	CompletedBytes int64     `json:"completed_bytes"`
	SpeedBps       int64     `json:"speed_bps"`
	ETASeconds     int       `json:"eta_seconds"`
	Status         Status    `json:"status"`
	Segments       []Segment `json:"segments"`
	HashAlgorithm  string    `json:"hash_algorithm"`
	HashDigest     string    `json:"hash_digest"`
	CreatedAt      time.Time `json:"created_at"`
}

// UploadPlan represents a staged upload awaiting or undergoing transfer.
type UploadPlan struct {
	ID          string    `json:"id"`
	SourcePath  string    `json:"source_path"`
	Destination string    `json:"destination"`
	SourceHash  string    `json:"source_hash"`
	TotalBytes  int64     `json:"total_bytes"`
	SentBytes   int64     `json:"sent_bytes"`
	SpeedBps    int64     `json:"speed_bps"`
	Reviewed    bool      `json:"reviewed"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// FileItem represents an item in the Zeus file manager / index.
type FileItem struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	SizeBytes int64     `json:"size_bytes"`
	IsDir     bool      `json:"is_dir"`
	ModTime   time.Time `json:"mod_time"`
	Hash      string    `json:"hash,omitempty"`
}

// AgentAction represents a proposed action from the natural language agent.
type AgentAction struct {
	ID               string                 `json:"id"`
	ActionType       string                 `json:"action_type"` // "download", "upload", "clone", "organize"
	Prompt           string                 `json:"prompt"`
	Summary          string                 `json:"summary"`
	Source           string                 `json:"source,omitempty"`
	Destination      string                 `json:"destination,omitempty"`
	RequiresApproval bool                   `json:"requires_approval"`
	Approved         bool                   `json:"approved"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
}

// DaemonStatus represents the health and high-level statistics of zeusd.
type DaemonStatus struct {
	Connected      bool   `json:"connected"`
	Version        string `json:"version"`
	DaemonURL      string `json:"daemon_url"`
	UptimeSeconds  int64  `json:"uptime_seconds"`
	ActiveDown     int    `json:"active_downloads"`
	ActiveUp       int    `json:"active_uploads"`
	TotalSpeedBps  int64  `json:"total_speed_bps"`
	StorageFreeGB  float64`json:"storage_free_gb"`
	StorageTotalGB float64`json:"storage_total_gb"`
}
