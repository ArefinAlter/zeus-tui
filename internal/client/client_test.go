package client

import (
	"context"
	"testing"
)

func TestClientMockStatus(t *testing.T) {
	c := New("http://127.0.0.1:18540", true)
	status, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if status.Version != "0.1.0-alpha.1" {
		t.Errorf("expected version 0.1.0-alpha.1, got %s", status.Version)
	}
	if status.ActiveDown < 1 {
		t.Errorf("expected at least 1 active download, got %d", status.ActiveDown)
	}
}

func TestClientDownloadsAndToggle(t *testing.T) {
	c := New("http://127.0.0.1:18540", true)
	downloads, err := c.ListDownloads(context.Background())
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(downloads) == 0 {
		t.Fatalf("expected downloads, got 0")
	}

	firstID := downloads[0].ID
	initialStatus := downloads[0].Status

	err = c.TogglePauseDownload(firstID)
	if err != nil {
		t.Fatalf("toggle pause error: %v", err)
	}

	downloadsAfter, _ := c.ListDownloads(context.Background())
	for _, d := range downloadsAfter {
		if d.ID == firstID {
			if initialStatus == StatusRunning && d.Status != StatusPaused {
				t.Errorf("expected status paused, got %s", d.Status)
			}
		}
	}
}

func TestAgentPromptWorkflow(t *testing.T) {
	c := New("http://127.0.0.1:18540", true)
	action, err := c.SubmitAgentPrompt(context.Background(), "stage upload to s3")
	if err != nil {
		t.Fatalf("submit agent prompt error: %v", err)
	}
	if action.ActionType != "upload" {
		t.Errorf("expected upload action, got %s", action.ActionType)
	}

	pending := c.GetPendingAction()
	if pending == nil || pending.ID != action.ID {
		t.Fatalf("expected pending action %s", action.ID)
	}

	err = c.ApproveAgentAction(action.ID)
	if err != nil {
		t.Fatalf("expected approve to succeed, got %v", err)
	}

	if c.GetPendingAction() != nil {
		t.Errorf("expected pending action to be cleared after approval")
	}
}
