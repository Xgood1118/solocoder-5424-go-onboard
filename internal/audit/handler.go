package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/audit")
	{
		g.GET("", listLogs)
	}
}

func listLogs(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID := c.Query("target_id")
	logs := ListLogs(targetType, targetID)
	c.JSON(http.StatusOK, logs)
}
