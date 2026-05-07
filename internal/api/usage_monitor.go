package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/usagemonitor"
)

type modelPricesRequest struct {
	Prices map[string]usagemonitor.ModelPrice `json:"prices"`
}

type modelPricesSyncRequest struct {
	Models []string `json:"models"`
}

func (s *Server) handleUsageHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "service": "cpa-manager"})
}

func (s *Server) handleUsageServiceInfo(c *gin.Context) {
	monitor := s.usageMonitor
	if monitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "usage monitor is not configured"})
		return
	}
	info, err := monitor.Info()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (s *Server) handleUsageStatus(c *gin.Context) {
	monitor := s.usageMonitor
	if monitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "usage monitor is not configured"})
		return
	}
	status, err := monitor.Status(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (s *Server) serveUsagePanel(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, usagemonitor.PanelHTML)
}

func (s *Server) handleManagementUsage(c *gin.Context) {
	monitor := s.usageMonitor
	if monitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "usage monitor is not configured"})
		return
	}

	switch c.Request.Method {
	case http.MethodGet:
		if strings.HasSuffix(c.Request.URL.Path, "/export") {
			data, err := monitor.ExportJSONL(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Header("Content-Type", "application/x-ndjson")
			c.Header("Content-Disposition", `attachment; filename="usage-events.jsonl"`)
			_, _ = c.Writer.Write(data)
			return
		}
		payload, err := monitor.UsagePayload(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, payload)
	case http.MethodPost:
		if strings.HasSuffix(c.Request.URL.Path, "/import") {
			s.handleManagementUsageImport(c, monitor)
			return
		}
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func (s *Server) handleManagementUsageImport(c *gin.Context, monitor *usagemonitor.Service) {
	reader := bufio.NewScanner(c.Request.Body)
	reader.Buffer(make([]byte, 64*1024), 10*1024*1024)
	events := make([]usagemonitor.Event, 0)
	failed := 0
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		event, err := usagemonitor.NormalizeRaw([]byte(line))
		if err != nil {
			failed++
			continue
		}
		events = append(events, event)
	}
	if err := reader.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := monitor.ImportJSONL(c.Request.Context(), events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"added":   result.Inserted,
		"skipped": result.Skipped,
		"total":   len(events),
		"failed":  failed,
	})
}

func (s *Server) handleManagementModelPrices(c *gin.Context) {
	monitor := s.usageMonitor
	if monitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "usage monitor is not configured"})
		return
	}

	path := strings.TrimRight(c.Request.URL.Path, "/")
	switch {
	case path == "/v0/management/model-prices" && c.Request.Method == http.MethodGet:
		prices, err := monitor.LoadModelPrices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"prices": prices})
	case path == "/v0/management/model-prices" && c.Request.Method == http.MethodPut:
		var req modelPricesRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Prices == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "prices are required"})
			return
		}
		if err := monitor.SaveModelPrices(c.Request.Context(), req.Prices); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prices, err := monitor.LoadModelPrices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"prices": prices})
	case path == "/v0/management/model-prices/sync" && c.Request.Method == http.MethodPost:
		var req modelPricesSyncRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		remotePrices, skipped, err := usagemonitor.FetchLiteLLMModelPrices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		selectedPrices := usagemonitor.SelectModelPrices(remotePrices, req.Models)
		result, err := monitor.UpsertSyncedModelPrices(c.Request.Context(), selectedPrices)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		prices, err := monitor.LoadModelPrices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"source":   usagemonitor.ModelPriceSyncSource,
			"imported": result.Imported,
			"skipped":  result.Skipped + skipped,
			"prices":   prices,
		})
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}
