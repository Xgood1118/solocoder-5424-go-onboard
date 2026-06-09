package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/assets")
	{
		g.POST("", createAsset)
		g.GET("", listAssets)
		g.GET("/:id", getAsset)
		g.PUT("/:id", updateAsset)
		g.POST("/:id/assign", assignAsset)
		g.POST("/:id/return", returnAsset)
		g.GET("/employee/:employee_id", getEmployeeAssets)
	}
}

func createAsset(c *gin.Context) {
	var a model.Asset
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := CreateAsset(&a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func listAssets(c *gin.Context) {
	assetType := c.Query("type")
	holderID := c.Query("holder_id")
	list := ListAssets(model.AssetType(assetType), holderID)
	c.JSON(http.StatusOK, list)
}

func getAsset(c *gin.Context) {
	id := c.Param("id")
	a, ok := GetAsset(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func updateAsset(c *gin.Context) {
	id := c.Param("id")
	var a model.Asset
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := UpdateAsset(id, &a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type assignReq struct {
	HolderID string `json:"holder_id" binding:"required"`
}

func assignAsset(c *gin.Context) {
	id := c.Param("id")
	var req assignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := AssignAsset(id, req.HolderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func returnAsset(c *gin.Context) {
	id := c.Param("id")
	if err := ReturnAsset(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "returned"})
}

func getEmployeeAssets(c *gin.Context) {
	empID := c.Param("employee_id")
	list := GetEmployeeAssets(empID)
	c.JSON(http.StatusOK, list)
}
