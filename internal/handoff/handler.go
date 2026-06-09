package handoff

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/handoffs")
	{
		g.POST("", createHandoff)
		g.GET("/:id", getHandoff)
		g.POST("/:id/items", addItem)
		g.POST("/:id/items/:item_id/complete", completeItem)
		g.DELETE("/:id/items/:item_id", removeItem)
		g.GET("/:id/completed", checkAllCompleted)
	}
}

type createHandoffReq struct {
	Items []model.HandoffItem `json:"items"`
}

func createHandoff(c *gin.Context) {
	var req createHandoffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h, err := CreateHandoff(req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h)
}

func getHandoff(c *gin.Context) {
	id := c.Param("id")
	h, ok := GetHandoff(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, h)
}

func addItem(c *gin.Context) {
	id := c.Param("id")
	var item model.HandoffItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h, err := AddHandoffItem(id, &item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h)
}

func completeItem(c *gin.Context) {
	handoffID := c.Param("id")
	itemID := c.Param("item_id")
	h, err := CompleteHandoffItem(handoffID, itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h)
}

func removeItem(c *gin.Context) {
	handoffID := c.Param("id")
	itemID := c.Param("item_id")
	h, err := RemoveHandoffItem(handoffID, itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h)
}

func checkAllCompleted(c *gin.Context) {
	id := c.Param("id")
	completed, err := IsAllCompleted(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"all_completed": completed})
}
