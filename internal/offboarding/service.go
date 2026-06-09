package offboarding

import (
	"errors"
	"fmt"
	"math"
	"time"

	"hr-onboard/internal/approval"
	"hr-onboard/internal/asset"
	"hr-onboard/internal/audit"
	"hr-onboard/internal/handoff"
	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func getNoticeDays(level model.EmployeeLevel) int {
	if level == model.LevelExecutive {
		return 90
	}
	return 30
}

func buildOffboardingNodes() []model.ProcessNode {
	return []model.ProcessNode{
		{NodeKey: "apply", NodeName: "离职申请", Status: model.StatusCompleted},
		{NodeKey: "manager_confirm", NodeName: "直属上级确认", Status: model.StatusPending},
		{NodeKey: "dept_approve", NodeName: "部门负责人审批", Status: model.StatusPending},
		{NodeKey: "hr_approve", NodeName: "HR审批", Status: model.StatusPending},
		{NodeKey: "handoff", NodeName: "工作交接", Status: model.StatusPending},
		{NodeKey: "asset_return", NodeName: "资产归还", Status: model.StatusPending},
		{NodeKey: "finance_settle", NodeName: "财务结算", Status: model.StatusPending},
		{NodeKey: "exit_interview", NodeName: "离职面谈", Status: model.StatusPending},
		{NodeKey: "cert_issue", NodeName: "离职证明开具", Status: model.StatusPending},
		{NodeKey: "social_stop", NodeName: "社保停缴", Status: model.StatusPending},
	}
}

type ApplyOffboardingRequest struct {
	EmployeeID   string
	EmployeeName string
	DepartmentID string
	Level        model.EmployeeLevel
	ApplyDate    time.Time
}

func ApplyOffboarding(req ApplyOffboardingRequest) (*model.OffboardingProcess, error) {
	if req.EmployeeID == "" {
		return nil, errors.New("员工ID不能为空")
	}

	noticeDays := getNoticeDays(req.Level)
	lastWorkingDate := req.ApplyDate.AddDate(0, 0, noticeDays)

	nodes := buildOffboardingNodes()
	now := time.Now()
	nodes[0].ActualStartAt = &now
	nodes[0].ActualEndAt = &now

	proc := &model.OffboardingProcess{
		ID:              genID("off"),
		EmployeeID:      req.EmployeeID,
		EmployeeName:    req.EmployeeName,
		DepartmentID:    req.DepartmentID,
		Level:           req.Level,
		Status:          model.OffboardApplied,
		ApplyDate:       req.ApplyDate,
		LastWorkingDate: lastWorkingDate,
		NoticeDays:      noticeDays,
		Nodes:           nodes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	s.SaveOffboarding(proc)
	return proc, nil
}

func GetOffboarding(id string) (*model.OffboardingProcess, bool) {
	return s.GetOffboarding(id)
}

func ListOffboardings(status model.OffboardingStatus) []*model.OffboardingProcess {
	var list []*model.OffboardingProcess
	s.RangeOffboardings(func(_ string, o *model.OffboardingProcess) bool {
		if status != "" && o.Status != status {
			return true
		}
		list = append(list, o)
		return true
	})
	return list
}

func updateNodeStatus(proc *model.OffboardingProcess, nodeKey string, status model.NodeStatus) {
	now := time.Now()
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			if status == model.StatusInProgress && proc.Nodes[i].ActualStartAt == nil {
				proc.Nodes[i].ActualStartAt = &now
			}
			if status == model.StatusCompleted {
				proc.Nodes[i].ActualEndAt = &now
			}
			proc.Nodes[i].Status = status
			break
		}
	}
}

func ApproveManager(id, operator string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardApplied {
		return nil, errors.New("当前状态不能进行直属上级确认")
	}
	proc.Status = model.OffboardManagerConfirmed
	updateNodeStatus(proc, "manager_confirm", model.StatusCompleted)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func ApproveDept(id, operator string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardManagerConfirmed {
		return nil, errors.New("当前状态不能进行部门负责人审批")
	}
	proc.Status = model.OffboardDeptApproved
	updateNodeStatus(proc, "dept_approve", model.StatusCompleted)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func ApproveHR(id, operator, approvalFlowID string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardDeptApproved {
		return nil, errors.New("当前状态不能进行HR审批")
	}

	if approvalFlowID != "" {
		inst, err := approval.StartApproval(approval.StartApprovalRequest{
			FlowID:      approvalFlowID,
			Subject:     fmt.Sprintf("离职审批: %s", proc.EmployeeName),
			InitiatorID: operator,
		})
		if err != nil {
			return nil, err
		}
		proc.ApprovalID = inst.ID
	}

	proc.Status = model.OffboardHROnApproved
	updateNodeStatus(proc, "hr_approve", model.StatusCompleted)
	updateNodeStatus(proc, "handoff", model.StatusInProgress)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func StartHandoff(id string, items []model.HandoffItem) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardHROnApproved {
		return nil, errors.New("当前状态不能发起交接")
	}

	h, err := handoff.CreateHandoff(items)
	if err != nil {
		return nil, err
	}

	proc.HandoffID = h.ID
	updateNodeStatus(proc, "handoff", model.StatusInProgress)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func CompleteHandoff(id string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.HandoffID == "" {
		return nil, errors.New("未发起交接")
	}
	allDone, err := handoff.IsAllCompleted(proc.HandoffID)
	if err != nil {
		return nil, err
	}
	if !allDone {
		return nil, errors.New("还有未完成的交接项")
	}
	updateNodeStatus(proc, "handoff", model.StatusCompleted)
	proc.Status = model.OffboardAssetReturnIn
	updateNodeStatus(proc, "asset_return", model.StatusInProgress)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func InitReturnableAssets(id string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}

	empAssets := asset.GetEmployeeAssets(proc.EmployeeID)
	var returnable []model.ReturnableAsset
	for _, a := range empAssets {
		depVal := asset.CalculateDepreciation(a.OriginalValue, a.PurchaseDate)
		returnable = append(returnable, model.ReturnableAsset{
			AssetID:          a.ID,
			AssetType:        a.Type,
			AssetName:        a.Name,
			Returned:         false,
			DepreciatedValue: depVal,
		})
	}
	proc.ReturnableAssets = returnable
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func ReturnAsset(id, assetID string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}

	now := time.Now()
	found := false
	for i := range proc.ReturnableAssets {
		if proc.ReturnableAssets[i].AssetID == assetID {
			proc.ReturnableAssets[i].Returned = true
			proc.ReturnableAssets[i].ReturnedAt = &now
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("资产不在待归还列表中")
	}

	if err := asset.ReturnAsset(assetID); err != nil {
		return nil, err
	}

	proc.UpdatedAt = now
	s.SaveOffboarding(proc)
	return proc, nil
}

func allAssetsReturned(proc *model.OffboardingProcess) bool {
	if len(proc.ReturnableAssets) == 0 {
		return true
	}
	for _, a := range proc.ReturnableAssets {
		if !a.Returned {
			return false
		}
	}
	return true
}

func CalculateSettlement(id string) (*model.FinancialSettlement, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}

	emp, empOk := s.GetEmployee(proc.EmployeeID)
	holiday, holOk := s.GetHolidayBalance(proc.EmployeeID)

	dailySalary := 0.0
	if empOk {
		dailySalary = emp.BaseSalary / 21.75
	}

	annualLeaveConvert := 0.0
	if holOk {
		annualLeaveConvert = holiday.AnnualLeave * dailySalary
	}

	remainingSalary := 0.0
	if empOk {
		remainingSalary = emp.BaseSalary / 30 * float64(proc.NoticeDays) / 3
	}

	stockBuyback := 0.0

	var unreturnedValue float64
	var unreturnedIDs []string
	for _, a := range proc.ReturnableAssets {
		if !a.Returned {
			unreturnedValue += a.DepreciatedValue
			unreturnedIDs = append(unreturnedIDs, a.AssetID)
		}
	}

	totalPayable := remainingSalary + annualLeaveConvert + stockBuyback

	assetDeduction := unreturnedValue
	if assetDeduction > totalPayable {
		assetDeduction = totalPayable
	}

	unpaidDebt := unreturnedValue - assetDeduction
	if unpaidDebt < 0 {
		unpaidDebt = 0
	}

	totalPayable = totalPayable - assetDeduction
	if totalPayable < 0 {
		totalPayable = 0
	}

	settlement := &model.FinancialSettlement{
		RemainingSalary:    math.Round(remainingSalary*100) / 100,
		AnnualLeaveConvert: math.Round(annualLeaveConvert*100) / 100,
		StockBuyback:       math.Round(stockBuyback*100) / 100,
		AssetDeduction:     math.Round(assetDeduction*100) / 100,
		UnpaidAssetsDebt:   math.Round(unpaidDebt*100) / 100,
		TotalPayable:       math.Round(totalPayable*100) / 100,
	}

	if unpaidDebt > 0 {
		debt := &model.PendingDebtItem{
			ID:           genID("debt"),
			EmployeeID:   proc.EmployeeID,
			EmployeeName: proc.EmployeeName,
			Amount:       math.Round(unpaidDebt*100) / 100,
			AssetIDs:     unreturnedIDs,
			OffboardID:   id,
			Status:       "pending",
			CreatedAt:    time.Now(),
		}
		s.SavePendingDebt(debt)
	}

	proc.Settlement = settlement
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)

	return settlement, nil
}

func CompleteAssetReturnAndSettle(id string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardAssetReturnIn {
		return nil, errors.New("当前状态不能完成资产归还")
	}

	updateNodeStatus(proc, "asset_return", model.StatusCompleted)

	_, err := CalculateSettlement(id)
	if err != nil {
		return nil, err
	}

	proc.Status = model.OffboardFinanceSettled
	updateNodeStatus(proc, "finance_settle", model.StatusCompleted)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func DoExitInterview(id, remark string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardFinanceSettled {
		return nil, errors.New("当前状态不能进行离职面谈")
	}
	proc.InterviewRemark = remark
	proc.Status = model.OffboardInterviewDone
	updateNodeStatus(proc, "exit_interview", model.StatusCompleted)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func IssueCertificate(id, certURL string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardInterviewDone {
		return nil, errors.New("当前状态不能开具离职证明")
	}
	proc.CertificateURL = certURL
	proc.Status = model.OffboardCertIssued
	updateNodeStatus(proc, "cert_issue", model.StatusCompleted)
	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func StopSocialInsurance(id string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}
	if proc.Status != model.OffboardCertIssued {
		return nil, errors.New("当前状态不能停缴社保")
	}
	now := time.Now()
	proc.SocialStopDate = &now
	proc.Status = model.OffboardSocialStopped
	updateNodeStatus(proc, "social_stop", model.StatusCompleted)
	proc.Status = model.OffboardCompleted
	proc.UpdatedAt = now
	s.SaveOffboarding(proc)
	return proc, nil
}

func CheckOverdueAndRemind() []*model.ReminderMessage {
	now := time.Now()
	var reminders []*model.ReminderMessage

	s.RangeOffboardings(func(_ string, proc *model.OffboardingProcess) bool {
		if proc.Status == model.OffboardCompleted {
			return true
		}

		for _, node := range proc.Nodes {
			if node.Status == model.StatusPending || node.Status == model.StatusInProgress {
				if node.PlannedEndAt != nil && node.PlannedEndAt.Before(now) {
					msg := &model.ReminderMessage{
						ID:         genID("rem"),
						ReceiverID: node.Assignee,
						Content:    fmt.Sprintf("节点[%s]已超期，请尽快处理", node.NodeName),
						ProcessID:  proc.ID,
						NodeKey:    node.NodeKey,
						CreatedAt:  now,
						Read:       false,
					}
					s.SaveReminder(msg)
					reminders = append(reminders, msg)
				}
				break
			}
		}
		return true
	})

	return reminders
}

func ListPendingDebts() []*model.PendingDebtItem {
	var list []*model.PendingDebtItem
	s.RangePendingDebts(func(_ string, d *model.PendingDebtItem) bool {
		list = append(list, d)
		return true
	})
	return list
}

func ForceSkip(id, nodeKey, operator, reason string) (*model.OffboardingProcess, error) {
	proc, ok := s.GetOffboarding(id)
	if !ok {
		return nil, errors.New("离职流程不存在")
	}

	found := false
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			now := time.Now()
			proc.Nodes[i].Status = model.StatusSkipped
			proc.Nodes[i].ActualEndAt = &now
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("节点不存在")
	}

	audit.AddLog(operator, "force_skip_offboard", id, "offboarding", reason,
		fmt.Sprintf("跳过节点: %s", nodeKey))

	proc.UpdatedAt = time.Now()
	s.SaveOffboarding(proc)
	return proc, nil
}

func ListReminders(receiverID string) []*model.ReminderMessage {
	var list []*model.ReminderMessage
	s.RangeReminders(func(_ string, r *model.ReminderMessage) bool {
		if receiverID != "" && r.ReceiverID != receiverID {
			return true
		}
		list = append(list, r)
		return true
	})
	return list
}
