package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ArminDashti/susiana-api/internal/domain"
	"github.com/ArminDashti/susiana-api/internal/repository"
	"github.com/ArminDashti/susiana-api/internal/service"
	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	svc *service.MarketService
}

func NewMarketHandler(svc *service.MarketService) *MarketHandler {
	return &MarketHandler{svc: svc}
}

func (h *MarketHandler) Register(r *gin.Engine) {
	r.GET("/health", h.Health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/stocks", h.ListStocks)
		v1.GET("/stocks/:isin", h.GetStock)
		v1.GET("/etfs", h.ListETFs)
		v1.GET("/etfs/:isin", h.GetETF)
		v1.POST("/sync", h.Sync)
		v1.GET("/sync/runs", h.ListSyncRuns)
	}
}

func (h *MarketHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *MarketHandler) ListStocks(c *gin.Context) {
	params := parseListParams(c)
	items, total, err := h.svc.ListStocks(c.Request.Context(), params)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"page":  params.NormalizedPage(),
		"limit": params.NormalizedLimit(),
		"total": total,
	})
}

func (h *MarketHandler) GetStock(c *gin.Context) {
	isin := c.Param("isin")
	item, err := h.svc.GetStock(c.Request.Context(), isin)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(c, http.StatusNotFound, "stock not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *MarketHandler) ListETFs(c *gin.Context) {
	params := parseListParams(c)
	items, total, err := h.svc.ListETFs(c.Request.Context(), params)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"page":  params.NormalizedPage(),
		"limit": params.NormalizedLimit(),
		"total": total,
	})
}

func (h *MarketHandler) GetETF(c *gin.Context) {
	isin := c.Param("isin")
	item, err := h.svc.GetETF(c.Request.Context(), isin)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(c, http.StatusNotFound, "etf not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *MarketHandler) Sync(c *gin.Context) {
	typeParam := c.DefaultQuery("type", string(domain.InstrumentTypeAll))
	instrumentType := domain.InstrumentType(typeParam)

	result, err := h.svc.Sync(c.Request.Context(), instrumentType)
	if err != nil {
		status := http.StatusBadRequest
		if result != nil && result.Run != nil && result.Run.Status == domain.SyncStatusFailed {
			status = http.StatusBadGateway
		}
		if result != nil {
			c.JSON(status, gin.H{"error": err.Error(), "run": result.Run, "fetched_count": result.FetchedCount})
			return
		}
		writeError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MarketHandler) ListSyncRuns(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	runs, err := h.svc.ListSyncRuns(c.Request.Context(), limit)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": runs})
}

func parseListParams(c *gin.Context) domain.ListParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	return domain.ListParams{
		Page:   page,
		Limit:  limit,
		Ticker: c.Query("ticker"),
	}
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
