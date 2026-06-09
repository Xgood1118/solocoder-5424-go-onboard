package onboarding

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/onboarding")
	{
		g.POST("", createOnboarding)
		g.GET("", listOnboardings)
		g.GET("/:id", getOnboarding)
		g.POST("/:id/transition", transitionStatus)
		g.POST("/:id/force_skip", forceSkip)

		g.POST("/:id/defense", createDefense)
		g.POST("/defense/:defense_id/result", submitDefenseResult)

		g.GET("/todo/today", getTodayTodo)
	}
}

type createOnboardingReq struct {
	EmployeeID   string            `json:"employee_id" binding:"required"`
	EmployeeName string            `json:"employee_name"`
	DepartmentID string            `json:"department_id"`
	Level        model.EmployeeLevel `json:"level"`
	Assignees    map[string]string `json:"assignees"`
}

func createOnboarding(c *gin.Context) {
	var req createOnboardingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := CreateOnboarding(CreateOnboardingRequest{
		EmployeeID:   req.EmployeeID,
		EmployeeName: req.EmployeeName,
		DepartmentID: req.DepartmentID,
		Level:        req.Level,
		Assignees:    req.Assignees,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func listOnboardings(c *gin.Context) {
	status := c.Query("status")
	list := ListOnboardings(model.OnboardingStatus(status))
	c.JSON(http.StatusOK, list)
}

func getOnboarding(c *gin.Context) {
	id := c.Param("id")
	proc, ok := GetOnboarding(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, proc)
}

type transitionReq struct {
	TargetStatus model.OnboardingStatus `json:"target_status" binding:"required"`
	Operator     string                 `json:"operator"`
	Remark       string                 `json:"remark"`
}

func transitionStatus(c *gin.Context) {
	id := c.Param("id")
	var req transitionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := TransitionStatus(id, req.TargetStatus, req.Operator, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

type forceSkipReq struct {
	NodeKey  string `json:"node_key" binding:"required"`
	Operator string `json:"operator" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

func forceSkip(c *gin.Context) {
	id := c.Param("id")
	var req forceSkipReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := ForceSkip(id, req.NodeKey, req.Operator, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

type createDefenseReq struct {
	Panel []model.DefensePanelMember `json:"panel" binding:"required"`
}

func createDefense(c *gin.Context) {
	id := c.Param("id")
	var req createDefenseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	def, err := CreateDefense(id, req.Panel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

type defenseResultReq struct {
	Result       model.DefenseResult `json:"result" binding:"required"`
	ExtendMonths int                 `json:"extend_months"`
	Remark       string              `json:"remark"`
}

func submitDefenseResult(c *gin.Context) {
	defenseID := c.Param("defense_id")
	var req defenseResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	def, err := SubmitDefenseResult(defenseID, req.Result, req.ExtendMonths, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

func getTodayTodo(c *gin.Context) {
	todo := GetTodayTodo()
	c.JSON(http.StatusOK, todo)
}
