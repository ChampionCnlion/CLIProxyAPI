package usagemonitor

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	internallogging "github.com/router-for-me/CLIProxyAPI/v6/internal/logging"
	coreusage "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/usage"
	log "github.com/sirupsen/logrus"
)

const (
	serviceID         = "cpa-manager"
	defaultQueryLimit = 50000
	defaultBatchSize  = 100
	flushInterval     = 500 * time.Millisecond
)

type CollectorStatus struct {
	Collector      string `json:"collector"`
	Upstream       string `json:"upstream"`
	Mode           string `json:"mode"`
	Transport      string `json:"transport"`
	Queue          string `json:"queue"`
	LastConsumedAt int64  `json:"lastConsumedAt"`
	LastInsertedAt int64  `json:"lastInsertedAt"`
	TotalInserted  int64  `json:"totalInserted"`
	TotalSkipped   int64  `json:"totalSkipped"`
	DeadLetters    int64  `json:"deadLetters"`
	LastError      string `json:"lastError,omitempty"`
}

type StatusResponse struct {
	Service     string          `json:"service"`
	DBPath      string          `json:"dbPath"`
	Events      int64           `json:"events"`
	DeadLetters int64           `json:"deadLetters"`
	Collector   CollectorStatus `json:"collector"`
}

type InfoResponse struct {
	Service   string `json:"service"`
	Mode      string `json:"mode"`
	StartedAt int64  `json:"startedAt"`
}

type Service struct {
	mu         sync.RWMutex
	store      *Store
	dbPath     string
	queryLimit int
	enabled    bool
	startedAt  int64
	status     CollectorStatus
	events     chan Event
	startOnce  sync.Once
}

var (
	global     = newService()
	registerMu sync.Once
)

func newService() *Service {
	return &Service{
		queryLimit: defaultQueryLimit,
		startedAt:  time.Now().UnixMilli(),
		events:     make(chan Event, 4096),
		status: CollectorStatus{
			Collector: "stopped",
			Upstream:  "in-process",
			Mode:      "embedded",
			Transport: "direct",
			Queue:     "in-process",
		},
	}
}

func Configure(cfg *config.Config, configFilePath string) (*Service, error) {
	registerMu.Do(func() {
		coreusage.RegisterPlugin(global)
	})
	global.startOnce.Do(func() {
		go global.runWriter()
	})
	if err := global.configure(cfg, configFilePath); err != nil {
		return global, err
	}
	return global, nil
}

func Current() *Service {
	return global
}

func (s *Service) configure(cfg *config.Config, configFilePath string) error {
	if s == nil {
		return errors.New("usage monitor is nil")
	}
	path := ResolveDBPath(cfg, configFilePath)
	queryLimit := defaultQueryLimit
	enabled := false
	if cfg != nil {
		enabled = cfg.UsageStatisticsEnabled
		if cfg.UsageQueryLimit > 0 {
			queryLimit = cfg.UsageQueryLimit
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store != nil && s.dbPath == path {
		s.enabled = enabled
		s.queryLimit = queryLimit
		s.status.Collector = collectorState(enabled)
		s.status.LastError = ""
		return nil
	}

	if s.store != nil {
		if err := s.store.Close(); err != nil {
			log.WithError(err).Warn("failed to close previous usage sqlite store")
		}
		s.store = nil
	}

	store, err := Open(path)
	if err != nil {
		s.dbPath = path
		s.enabled = false
		s.queryLimit = queryLimit
		s.status.Collector = "error"
		s.status.LastError = "open sqlite: " + err.Error()
		return err
	}

	s.store = store
	s.dbPath = path
	s.enabled = enabled
	s.queryLimit = queryLimit
	s.status.Collector = collectorState(enabled)
	s.status.Upstream = "in-process"
	s.status.Mode = "embedded"
	s.status.Transport = "direct"
	s.status.Queue = "in-process"
	s.status.LastError = ""
	return nil
}

func ResolveDBPath(cfg *config.Config, configFilePath string) string {
	raw := ""
	if cfg != nil {
		raw = strings.TrimSpace(cfg.UsageDBPath)
	}
	base := configBaseDir(configFilePath)
	if raw == "" {
		return filepath.Join(base, "usage.sqlite")
	}
	if filepath.IsAbs(raw) {
		return raw
	}
	return filepath.Join(base, raw)
}

func configBaseDir(configFilePath string) string {
	configFilePath = strings.TrimSpace(configFilePath)
	if configFilePath == "" {
		if wd, err := os.Getwd(); err == nil {
			return wd
		}
		return "."
	}
	info, err := os.Stat(configFilePath)
	if err == nil && info.IsDir() {
		return configFilePath
	}
	return filepath.Dir(configFilePath)
}

func collectorState(enabled bool) string {
	if enabled {
		return "running"
	}
	return "stopped"
}

func (s *Service) Close() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.store == nil {
		return nil
	}
	err := s.store.Close()
	s.store = nil
	s.status.Collector = "stopped"
	return err
}

func (s *Service) HandleUsage(ctx context.Context, record coreusage.Record) {
	if s == nil {
		return
	}

	payload, err := marshalRecord(ctx, record)
	if err != nil {
		s.markError("marshal", err)
		return
	}

	event, err := NormalizeRaw(payload)
	if err != nil {
		s.addDeadLetter(ctx, string(payload), err)
		return
	}

	s.mu.Lock()
	if !s.enabled || s.store == nil {
		s.mu.Unlock()
		return
	}
	s.status.LastConsumedAt = time.Now().UnixMilli()
	s.mu.Unlock()

	s.events <- event
}

func (s *Service) runWriter() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	batch := make([]Event, 0, defaultBatchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		events := append([]Event(nil), batch...)
		batch = batch[:0]
		s.insertBatch(context.Background(), events)
	}

	for {
		select {
		case event := <-s.events:
			batch = append(batch, event)
			if len(batch) >= defaultBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (s *Service) insertBatch(ctx context.Context, events []Event) {
	if len(events) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.enabled || s.store == nil {
		return
	}
	result, err := s.store.InsertEvents(ctx, events)
	if err != nil {
		s.status.Collector = "error"
		s.status.LastError = "insert: " + err.Error()
		return
	}
	if result.Inserted > 0 || result.Skipped > 0 {
		s.status.LastInsertedAt = time.Now().UnixMilli()
		s.status.TotalInserted += int64(result.Inserted)
		s.status.TotalSkipped += int64(result.Skipped)
		s.status.LastError = ""
		s.status.Collector = collectorState(s.enabled)
	}
}

func (s *Service) addDeadLetter(ctx context.Context, payload string, err error) {
	s.mu.RLock()
	if s.store != nil {
		_ = s.store.AddDeadLetter(ctx, payload, err)
	}
	s.mu.RUnlock()
	s.mu.Lock()
	s.status.DeadLetters++
	s.status.LastError = "normalize: " + err.Error()
	s.mu.Unlock()
}

func (s *Service) markError(stage string, err error) {
	if s == nil || err == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Collector = "error"
	s.status.LastError = stage + ": " + err.Error()
}

func (s *Service) Info() (InfoResponse, error) {
	if s == nil {
		return InfoResponse{}, errors.New("usage monitor is not configured")
	}
	s.mu.RLock()
	startedAt := s.startedAt
	s.mu.RUnlock()
	return InfoResponse{
		Service:   serviceID,
		Mode:      "embedded",
		StartedAt: startedAt,
	}, nil
}

func (s *Service) Status(ctx context.Context) (StatusResponse, error) {
	if s == nil {
		return StatusResponse{}, errors.New("usage monitor is not configured")
	}
	s.mu.RLock()
	store := s.store
	dbPath := s.dbPath
	status := s.status
	s.mu.RUnlock()

	var events, deadLetters int64
	var err error
	if store != nil {
		events, deadLetters, err = store.Counts(ctx)
		if err != nil {
			return StatusResponse{}, err
		}
		status.DeadLetters = deadLetters
	}

	return StatusResponse{
		Service:     serviceID,
		DBPath:      dbPath,
		Events:      events,
		DeadLetters: deadLetters,
		Collector:   status,
	}, nil
}

func (s *Service) UsagePayload(ctx context.Context) (Payload, error) {
	events, err := s.RecentEvents(ctx)
	if err != nil {
		return Payload{}, err
	}
	return BuildPayload(events), nil
}

func (s *Service) RecentEvents(ctx context.Context) ([]Event, error) {
	s.mu.RLock()
	store := s.store
	limit := s.queryLimit
	s.mu.RUnlock()
	if store == nil {
		return nil, errors.New("usage sqlite store is not configured")
	}
	return store.RecentEvents(ctx, limit)
}

func (s *Service) ExportJSONL(ctx context.Context) ([]byte, error) {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		return nil, errors.New("usage sqlite store is not configured")
	}
	return store.ExportJSONL(ctx)
}

func (s *Service) ImportJSONL(ctx context.Context, events []Event) (InsertResult, error) {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		return InsertResult{}, errors.New("usage sqlite store is not configured")
	}
	return store.InsertEvents(ctx, events)
}

func (s *Service) LoadModelPrices(ctx context.Context) (map[string]ModelPrice, error) {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		return nil, errors.New("usage sqlite store is not configured")
	}
	return store.LoadModelPrices(ctx)
}

func (s *Service) SaveModelPrices(ctx context.Context, prices map[string]ModelPrice) error {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		return errors.New("usage sqlite store is not configured")
	}
	return store.SaveModelPrices(ctx, prices)
}

func (s *Service) UpsertSyncedModelPrices(ctx context.Context, prices map[string]ModelPrice) (ModelPriceSyncResult, error) {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		return ModelPriceSyncResult{}, errors.New("usage sqlite store is not configured")
	}
	return store.UpsertSyncedModelPrices(ctx, prices)
}

func marshalRecord(ctx context.Context, record coreusage.Record) ([]byte, error) {
	timestamp := record.RequestedAt
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	modelName := strings.TrimSpace(record.Model)
	if modelName == "" {
		modelName = "unknown"
	}
	aliasName := strings.TrimSpace(record.Alias)
	if aliasName == "" {
		aliasName = modelName
	}
	provider := strings.TrimSpace(record.Provider)
	if provider == "" {
		provider = "unknown"
	}
	authType := strings.TrimSpace(record.AuthType)
	if authType == "" {
		authType = "unknown"
	}

	tokens := tokenStats{
		InputTokens:     record.Detail.InputTokens,
		OutputTokens:    record.Detail.OutputTokens,
		ReasoningTokens: record.Detail.ReasoningTokens,
		CachedTokens:    record.Detail.CachedTokens,
		TotalTokens:     record.Detail.TotalTokens,
	}
	if tokens.TotalTokens == 0 {
		tokens.TotalTokens = tokens.InputTokens + tokens.OutputTokens + tokens.ReasoningTokens
	}
	if tokens.TotalTokens == 0 {
		tokens.TotalTokens = tokens.InputTokens + tokens.OutputTokens + tokens.ReasoningTokens + tokens.CachedTokens
	}

	failed := record.Failed
	if !failed {
		failed = !resolveSuccess(ctx)
	}

	return json.Marshal(queuedUsageDetail{
		Timestamp: timestamp,
		LatencyMs: record.Latency.Milliseconds(),
		Source:    record.Source,
		AuthIndex: record.AuthIndex,
		Tokens:    tokens,
		Failed:    failed,
		Provider:  provider,
		Model:     modelName,
		Alias:     aliasName,
		Endpoint:  strings.TrimSpace(internallogging.GetEndpoint(ctx)),
		AuthType:  authType,
		APIKey:    strings.TrimSpace(record.APIKey),
		RequestID: strings.TrimSpace(internallogging.GetRequestID(ctx)),
	})
}

type queuedUsageDetail struct {
	Timestamp time.Time  `json:"timestamp"`
	LatencyMs int64      `json:"latency_ms"`
	Source    string     `json:"source"`
	AuthIndex string     `json:"auth_index"`
	Tokens    tokenStats `json:"tokens"`
	Failed    bool       `json:"failed"`
	Provider  string     `json:"provider"`
	Model     string     `json:"model"`
	Alias     string     `json:"alias"`
	Endpoint  string     `json:"endpoint"`
	AuthType  string     `json:"auth_type"`
	APIKey    string     `json:"api_key"`
	RequestID string     `json:"request_id"`
}

type tokenStats struct {
	InputTokens     int64 `json:"input_tokens"`
	OutputTokens    int64 `json:"output_tokens"`
	ReasoningTokens int64 `json:"reasoning_tokens"`
	CachedTokens    int64 `json:"cached_tokens"`
	TotalTokens     int64 `json:"total_tokens"`
}

func resolveSuccess(ctx context.Context) bool {
	status := internallogging.GetResponseStatus(ctx)
	if status == 0 {
		return true
	}
	return status < httpStatusBadRequest
}

const httpStatusBadRequest = 400
