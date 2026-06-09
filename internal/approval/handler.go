package approval

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/approval")
	{
		g.POST("/flows", createFlow)
		g.GET("/flows", listFlows)
		g.GET("/flows/:id", getFlow)
		g.PUT("/flows/:id", updateFlow)

		g.POST("/instances", startInstance)
		g.GET("/instances", listInstances)
		g.GET("/instances/:id", getInstance)
		g.POST("/instances/:id/approve", approveInstance)
		g.POST("/instances/:id/reject", rejectInstance)
	}
}

func createFlow(c *gin.Context) {
	var flow model.ApprovalFlowConfig
	if err := c.ShouldBindJSON(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := CreateFlow(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flow)
}

func listFlows(c *gin.Context) {
	flows := ListFlows()
	c.JSON(http.StatusOK, flows)
}

func getFlow(c *gin.Context) {
	id := c.Param("id")
	flow, ok := GetFlow(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}
	c.JSON(http.StatusOK, flow)
}

func updateFlow(c *gin.Context) {
	id := c.Param("id")
	var flow model.ApprovalFlowConfig
	if err := c.ShouldBindJSON(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flow.ID = id
	if err := UpdateFlow(&flow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flow)
}

type startInstanceReq struct {
	FlowID      string `json:"flow_id" binding:"required"`
	Subject     string `json:"subject"`
	InitiatorID string `json:"initiator_id" binding:"required"`
}

func startInstance(c *gin.Context) {
	var req startInstanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inst, err := StartApproval(StartApprovalRequest{
		FlowID:      req.FlowID,
		Subject:     req.Subject,
		InitiatorID: req.InitiatorID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inst)
}

func listInstances(c *gin.Context) {
	list := ListInstances()
	c.JSON(http.StatusOK, list)
}

func getInstance(c *gin.Context) {
	id := c.Param("id")
	inst, ok := GetApprovalInstance(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "instance not found"})
		return
	}
	c.JSON(http.StatusOK, inst)
}

type approveReq struct {
	ApproverID string `json:"approver_id" binding:"required"`
	Remark     string `json:"remark"`
}

func approveInstance(c *gin.Context) {
	id := c.Param("id")
	var req approveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inst, err := Approve(id, req.ApproverID, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inst)
}

func rejectInstance(c *gin.Context) {
	id := c.Param("id")
	var req approveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inst, err := Reject(id, req.ApproverID, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inst)
}
