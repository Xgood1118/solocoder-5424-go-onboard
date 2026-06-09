package offboarding

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/offboarding")
	{
		g.POST("", applyOffboarding)
		g.GET("", listOffboardings)
		g.GET("/:id", getOffboarding)
		g.POST("/:id/force_skip", forceSkip)

		g.POST("/:id/approve/manager", approveManager)
		g.POST("/:id/approve/dept", approveDept)
		g.POST("/:id/approve/hr", approveHR)

		g.POST("/:id/handoff/start", startHandoff)
		g.POST("/:id/handoff/complete", completeHandoff)

		g.POST("/:id/assets/init", initAssets)
		g.POST("/:id/assets/:asset_id/return", returnAsset)
		g.POST("/:id/assets/settle", settleAssets)

		g.POST("/:id/settlement/calc", calcSettlement)

		g.POST("/:id/interview", doInterview)
		g.POST("/:id/certificate", issueCert)
		g.POST("/:id/social_stop", stopSocial)

		g.GET("/debts/pending", listPendingDebts)
		g.GET("/reminders", listReminders)
		g.POST("/reminders/check", checkReminders)
	}
}

type applyReq struct {
	EmployeeID   string              `json:"employee_id" binding:"required"`
	EmployeeName string              `json:"employee_name"`
	DepartmentID string              `json:"department_id"`
	Level        model.EmployeeLevel `json:"level"`
	ApplyDate    string              `json:"apply_date"`
}

func applyOffboarding(c *gin.Context) {
	var req applyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	applyDate := time.Now()
	if req.ApplyDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ApplyDate)
		if err == nil {
			applyDate = parsed
		}
	}
	proc, err := ApplyOffboarding(ApplyOffboardingRequest{
		EmployeeID:   req.EmployeeID,
		EmployeeName: req.EmployeeName,
		DepartmentID: req.DepartmentID,
		Level:        req.Level,
		ApplyDate:    applyDate,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func listOffboardings(c *gin.Context) {
	status := c.Query("status")
	list := ListOffboardings(model.OffboardingStatus(status))
	c.JSON(http.StatusOK, list)
}

func getOffboarding(c *gin.Context) {
	id := c.Param("id")
	proc, ok := GetOffboarding(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
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

type approveReq struct {
	Operator string `json:"operator"`
	FlowID   string `json:"flow_id"`
}

func approveManager(c *gin.Context) {
	id := c.Param("id")
	var req approveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := ApproveManager(id, req.Operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func approveDept(c *gin.Context) {
	id := c.Param("id")
	var req approveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := ApproveDept(id, req.Operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func approveHR(c *gin.Context) {
	id := c.Param("id")
	var req approveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := ApproveHR(id, req.Operator, req.FlowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

type startHandoffReq struct {
	Items []model.HandoffItem `json:"items"`
}

func startHandoff(c *gin.Context) {
	id := c.Param("id")
	var req startHandoffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := StartHandoff(id, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func completeHandoff(c *gin.Context) {
	id := c.Param("id")
	proc, err := CompleteHandoff(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func initAssets(c *gin.Context) {
	id := c.Param("id")
	proc, err := InitReturnableAssets(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func returnAsset(c *gin.Context) {
	id := c.Param("id")
	assetID := c.Param("asset_id")
	proc, err := ReturnAsset(id, assetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func settleAssets(c *gin.Context) {
	id := c.Param("id")
	proc, err := CompleteAssetReturnAndSettle(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func calcSettlement(c *gin.Context) {
	id := c.Param("id")
	settle, err := CalculateSettlement(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settle)
}

type interviewReq struct {
	Remark string `json:"remark"`
}

func doInterview(c *gin.Context) {
	id := c.Param("id")
	var req interviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := DoExitInterview(id, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

type certReq struct {
	CertURL string `json:"cert_url"`
}

func issueCert(c *gin.Context) {
	id := c.Param("id")
	var req certReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proc, err := IssueCertificate(id, req.CertURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func stopSocial(c *gin.Context) {
	id := c.Param("id")
	proc, err := StopSocialInsurance(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proc)
}

func listPendingDebts(c *gin.Context) {
	list := ListPendingDebts()
	c.JSON(http.StatusOK, list)
}

func listReminders(c *gin.Context) {
	receiverID := c.Query("receiver_id")
	list := ListReminders(receiverID)
	c.JSON(http.StatusOK, list)
}

func checkReminders(c *gin.Context) {
	reminders := CheckOverdueAndRemind()
	c.JSON(http.StatusOK, gin.H{"new_reminders": len(reminders), "items": reminders})
}
