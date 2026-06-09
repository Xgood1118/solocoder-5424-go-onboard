package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/stats")
	{
		g.GET("/monthly", listMonthlyReports)
		g.GET("/monthly/:month", getMonthlyReport)
		g.POST("/monthly/generate", generateReport)
		g.GET("/current", getCurrentStats)
	}
}

func listMonthlyReports(c *gin.Context) {
	list := ListMonthlyReports()
	c.JSON(http.StatusOK, list)
}

func getMonthlyReport(c *gin.Context) {
	month := c.Param("month")
	r, ok := GetMonthlyReport(month)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, r)
}

type generateReq struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func generateReport(c *gin.Context) {
	var req generateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Year == 0 || req.Month == 0 {
		now := time.Now()
		lastMonth := now.AddDate(0, -1, 0)
		req.Year = lastMonth.Year()
		req.Month = int(lastMonth.Month())
	}
	report := GenerateMonthlyReport(req.Year, time.Month(req.Month))
	c.JSON(http.StatusOK, report)
}

func getCurrentStats(c *gin.Context) {
	stats := GetCurrentMonthStats()
	c.JSON(http.StatusOK, stats)
}
