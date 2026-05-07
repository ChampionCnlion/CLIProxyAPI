package usagemonitor

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	coreusage "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/usage"
)

func TestStoreInsertAndBuildUsagePayload(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "usage.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	event, err := NormalizeRaw([]byte(`{
		"timestamp":"2026-05-06T10:00:00Z",
		"request_id":"req-1",
		"provider":"codex",
		"model":"gpt-test",
		"endpoint":"POST /v1/chat/completions",
		"source":"user@example.com",
		"auth_index":"codex:1",
		"tokens":{"input_tokens":10,"output_tokens":20,"reasoning_tokens":5},
		"latency_ms":123,
		"failed":false
	}`))
	if err != nil {
		t.Fatalf("NormalizeRaw() error = %v", err)
	}

	result, err := store.InsertEvents(context.Background(), []Event{event, event})
	if err != nil {
		t.Fatalf("InsertEvents() error = %v", err)
	}
	if result.Inserted != 1 || result.Skipped != 1 {
		t.Fatalf("InsertEvents() = %+v, want inserted=1 skipped=1", result)
	}

	events, err := store.RecentEvents(context.Background(), 10)
	if err != nil {
		t.Fatalf("RecentEvents() error = %v", err)
	}
	payload := BuildPayload(events)
	if payload.TotalRequests != 1 || payload.SuccessCount != 1 || payload.TotalTokens != 35 {
		t.Fatalf("payload = %+v, want one successful 35-token request", payload)
	}
	api := payload.APIs["POST /v1/chat/completions"]
	if api == nil || api.Models["gpt-test"] == nil || len(api.Models["gpt-test"].Details) != 1 {
		t.Fatalf("payload missing model details: %+v", payload.APIs)
	}
	detail := api.Models["gpt-test"].Details[0]
	if detail.RequestID != "req-1" || detail.Provider != "codex" || detail.Method != "POST" || detail.Path != "/v1/chat/completions" {
		t.Fatalf("detail metadata = %+v, want request/provider/method/path", detail)
	}
}

func TestServiceHandleUsagePersistsThroughBatchWriter(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		UsageStatisticsEnabled: true,
		UsageDBPath:            filepath.Join(dir, "usage.sqlite"),
		UsageQueryLimit:        100,
	}
	service, err := Configure(cfg, filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	service.HandleUsage(context.Background(), coreusage.Record{
		Provider:    "codex",
		Model:       "gpt-test",
		AuthIndex:   "codex:1",
		AuthType:    "oauth",
		Source:      "tester@example.com",
		RequestedAt: time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC),
		Latency:     250 * time.Millisecond,
		Detail: coreusage.Detail{
			InputTokens:  7,
			OutputTokens: 9,
			TotalTokens:  16,
		},
	})

	deadline := time.Now().Add(3 * time.Second)
	for {
		events, err := service.RecentEvents(context.Background())
		if err != nil {
			t.Fatalf("RecentEvents() error = %v", err)
		}
		if len(events) == 1 {
			if events[0].TotalTokens != 16 || events[0].Model != "gpt-test" {
				t.Fatalf("event = %+v, want gpt-test with 16 tokens", events[0])
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for batch writer, events=%d", len(events))
		}
		time.Sleep(50 * time.Millisecond)
	}
}
