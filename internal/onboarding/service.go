package onboarding

import (
	"errors"
	"fmt"
	"time"

	"hr-onboard/internal/audit"
	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

var stateTransitionMap = map[model.OnboardingStatus][]model.OnboardingStatus{
	model.OnboardOfferConfirmed:  {model.OnboardHealthCheckIn},
	model.OnboardHealthCheckIn:   {model.OnboardHealthCheckDone},
	model.OnboardHealthCheckDone: {model.OnboardBgCheckIn, model.OnboardContractPending},
	model.OnboardBgCheckIn:       {model.OnboardBgCheckDone},
	model.OnboardBgCheckDone:     {model.OnboardContractPending},
	model.OnboardContractPending: {model.OnboardContractSigned},
	model.OnboardContractSigned:  {model.OnboardAccountSetupIn},
	model.OnboardAccountSetupIn:  {model.OnboardAccountDone},
	model.OnboardAccountDone:     {model.OnboardWorkspaceDone},
	model.OnboardWorkspaceDone:   {model.OnboardTrainingIn},
	model.OnboardTrainingIn:      {model.OnboardTrainingDone},
	model.OnboardTrainingDone:    {model.OnboardDeptIntroDone},
	model.OnboardDeptIntroDone:   {model.OnboardManagerMeetDone},
	model.OnboardManagerMeetDone: {model.OnboardProbationIn},
	model.OnboardProbationIn:     {model.OnboardDefenseIn},
	model.OnboardDefenseIn:       {model.OnboardFormal, model.OnboardExtended, model.OnboardDismissed},
	model.OnboardExtended:        {model.OnboardDefenseIn},
}

func canTransition(from, to model.OnboardingStatus) bool {
	nextStates, ok := stateTransitionMap[from]
	if !ok {
		return false
	}
	for _, ns := range nextStates {
		if ns == to {
			return true
		}
	}
	return false
}

func buildDefaultNodes(level model.EmployeeLevel) []model.ProcessNode {
	nodes := []model.ProcessNode{
		{NodeKey: "offer_confirm", NodeName: "Offer确认", Status: model.StatusPending},
		{NodeKey: "health_check", NodeName: "体检", Status: model.StatusPending},
	}
	if level == model.LevelExecutive {
		nodes = append(nodes, model.ProcessNode{NodeKey: "bg_check", NodeName: "背景调查", Status: model.StatusPending})
	}
	nodes = append(nodes, []model.ProcessNode{
		{NodeKey: "sign_contract", NodeName: "签合同", Status: model.StatusPending},
		{NodeKey: "account_setup", NodeName: "开通账号", Status: model.StatusPending},
		{NodeKey: "workspace", NodeName: "工位安排", Status: model.StatusPending},
		{NodeKey: "training", NodeName: "入职培训", Status: model.StatusPending},
		{NodeKey: "dept_intro", NodeName: "部门介绍", Status: model.StatusPending},
		{NodeKey: "manager_meet", NodeName: "直属上级见面", Status: model.StatusPending},
		{NodeKey: "probation", NodeName: "试用期", Status: model.StatusPending},
		{NodeKey: "defense", NodeName: "转正答辩", Status: model.StatusPending},
		{NodeKey: "formal", NodeName: "转正", Status: model.StatusPending},
	}...)
	return nodes
}

type CreateOnboardingRequest struct {
	EmployeeID   string
	EmployeeName string
	DepartmentID string
	Level        model.EmployeeLevel
	Assignees    map[string]string
}

func CreateOnboarding(req CreateOnboardingRequest) (*model.OnboardingProcess, error) {
	if req.EmployeeID == "" {
		return nil, errors.New("员工ID不能为空")
	}
	probationMonths := 3
	if req.Level == model.LevelExecutive {
		probationMonths = 6
	}

	now := time.Now()
	probStart := now.AddDate(0, 0, 0)
	probEnd := probStart.AddDate(0, probationMonths, 0)

	nodes := buildDefaultNodes(req.Level)
	for k, v := range req.Assignees {
		for i := range nodes {
			if nodes[i].NodeKey == k {
				nodes[i].Assignee = v
				break
			}
		}
	}

	proc := &model.OnboardingProcess{
		ID:                 genID("onb"),
		EmployeeID:         req.EmployeeID,
		EmployeeName:       req.EmployeeName,
		DepartmentID:       req.DepartmentID,
		Level:              req.Level,
		Status:             model.OnboardOfferConfirmed,
		Nodes:              nodes,
		ProbationMonths:    probationMonths,
		ProbationStartDate: &probStart,
		ProbationEndDate:   &probEnd,
		OfferConfirmDate:   &now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	markNodeStarted(proc, "offer_confirm", &now)
	markNodeCompleted(proc, "offer_confirm", &now)

	s.SaveOnboarding(proc)
	return proc, nil
}

func markNodeStarted(proc *model.OnboardingProcess, nodeKey string, t *time.Time) {
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			proc.Nodes[i].Status = model.StatusInProgress
			if proc.Nodes[i].ActualStartAt == nil {
				proc.Nodes[i].ActualStartAt = t
			}
			break
		}
	}
}

func markNodeCompleted(proc *model.OnboardingProcess, nodeKey string, t *time.Time) {
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			proc.Nodes[i].Status = model.StatusCompleted
			proc.Nodes[i].ActualEndAt = t
			break
		}
	}
}

func markNodeSkipped(proc *model.OnboardingProcess, nodeKey string, t *time.Time) {
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			proc.Nodes[i].Status = model.StatusSkipped
			proc.Nodes[i].ActualEndAt = t
			break
		}
	}
}

func findNode(proc *model.OnboardingProcess, nodeKey string) *model.ProcessNode {
	for i := range proc.Nodes {
		if proc.Nodes[i].NodeKey == nodeKey {
			return &proc.Nodes[i]
		}
	}
	return nil
}

func TransitionStatus(id string, targetStatus model.OnboardingStatus, operator, remark string) (*model.OnboardingProcess, error) {
	proc, ok := s.GetOnboarding(id)
	if !ok {
		return nil, errors.New("入职流程不存在")
	}
	if !canTransition(proc.Status, targetStatus) {
		return nil, fmt.Errorf("不能从 %s 转移到 %s", proc.Status, targetStatus)
	}

	now := time.Now()

	statusToNode := map[model.OnboardingStatus]string{
		model.OnboardHealthCheckIn:   "health_check",
		model.OnboardHealthCheckDone: "health_check",
		model.OnboardBgCheckIn:       "bg_check",
		model.OnboardBgCheckDone:     "bg_check",
		model.OnboardContractPending: "sign_contract",
		model.OnboardContractSigned:  "sign_contract",
		model.OnboardAccountSetupIn:  "account_setup",
		model.OnboardAccountDone:     "account_setup",
		model.OnboardWorkspaceDone:   "workspace",
		model.OnboardTrainingIn:      "training",
		model.OnboardTrainingDone:    "training",
		model.OnboardDeptIntroDone:   "dept_intro",
		model.OnboardManagerMeetDone: "manager_meet",
		model.OnboardProbationIn:     "probation",
		model.OnboardDefenseIn:       "defense",
		model.OnboardFormal:          "formal",
	}

	isStartState := map[model.OnboardingStatus]bool{
		model.OnboardHealthCheckIn:   true,
		model.OnboardBgCheckIn:       true,
		model.OnboardContractPending: true,
		model.OnboardAccountSetupIn:  true,
		model.OnboardTrainingIn:      true,
		model.OnboardDefenseIn:       true,
	}

	nodeKey, hasNode := statusToNode[targetStatus]
	if hasNode {
		if isStartState[targetStatus] {
			markNodeStarted(proc, nodeKey, &now)
		} else {
			markNodeCompleted(proc, nodeKey, &now)
		}
	}

	if targetStatus == model.OnboardFormal {
		proc.IsFormal = true
		proc.FormalDate = proc.ProbationEndDate
		markNodeCompleted(proc, "formal", proc.ProbationEndDate)
	}

	proc.Status = targetStatus
	proc.UpdatedAt = now
	if remark != "" {
		proc.Remarks = remark
	}

	s.SaveOnboarding(proc)
	return proc, nil
}

func ForceSkip(id, nodeKey, operator, reason string) (*model.OnboardingProcess, error) {
	proc, ok := s.GetOnboarding(id)
	if !ok {
		return nil, errors.New("入职流程不存在")
	}

	node := findNode(proc, nodeKey)
	if node == nil {
		return nil, errors.New("节点不存在")
	}

	now := time.Now()
	markNodeSkipped(proc, nodeKey, &now)

	nextStatusMap := map[string]model.OnboardingStatus{
		"offer_confirm":  model.OnboardHealthCheckIn,
		"health_check":   model.OnboardBgCheckIn,
		"bg_check":       model.OnboardContractPending,
		"sign_contract":  model.OnboardAccountSetupIn,
		"account_setup":  model.OnboardWorkspaceDone,
		"workspace":      model.OnboardTrainingIn,
		"training":       model.OnboardDeptIntroDone,
		"dept_intro":     model.OnboardManagerMeetDone,
		"manager_meet":   model.OnboardProbationIn,
		"probation":      model.OnboardDefenseIn,
		"defense":        model.OnboardFormal,
	}

	if nextStatus, ok := nextStatusMap[nodeKey]; ok {
		if nodeKey == "health_check" && proc.Level == model.LevelExecutive {
			proc.Status = model.OnboardBgCheckDone
			markNodeSkipped(proc, "bg_check", &now)
			proc.Status = model.OnboardContractPending
		} else {
			proc.Status = nextStatus
		}
	}

	proc.UpdatedAt = now

	audit.AddLog(operator, "force_skip", id, "onboarding", reason,
		fmt.Sprintf("跳过节点: %s", nodeKey))

	s.SaveOnboarding(proc)
	return proc, nil
}

func GetOnboarding(id string) (*model.OnboardingProcess, bool) {
	return s.GetOnboarding(id)
}

func ListOnboardings(status model.OnboardingStatus) []*model.OnboardingProcess {
	var list []*model.OnboardingProcess
	s.RangeOnboardings(func(_ string, o *model.OnboardingProcess) bool {
		if status != "" && o.Status != status {
			return true
		}
		list = append(list, o)
		return true
	})
	return list
}

type DefensePanelConfig struct {
	IncludeManager      bool
	IncludeDeptHead     bool
	IncludeHRBP         bool
	CrossDepartmentCount int
}

func GetDefaultPanelConfig() DefensePanelConfig {
	return DefensePanelConfig{
		IncludeManager:       true,
		IncludeDeptHead:      true,
		IncludeHRBP:          true,
		CrossDepartmentCount: 1,
	}
}

func CreateDefense(onboardingID string, panel []model.DefensePanelMember) (*model.ProbationDefense, error) {
	if len(panel) < 3 {
		return nil, errors.New("答辩小组至少需要3人")
	}
	proc, ok := s.GetOnboarding(onboardingID)
	if !ok {
		return nil, errors.New("入职流程不存在")
	}
	if proc.Status != model.OnboardProbationIn && proc.Status != model.OnboardExtended {
		return nil, errors.New("当前状态不能创建答辩")
	}

	defense := &model.ProbationDefense{
		ID:           genID("def"),
		OnboardingID: onboardingID,
		EmployeeID:   proc.EmployeeID,
		Panel:        panel,
	}
	s.SaveDefense(defense)

	now := time.Now()
	proc.Status = model.OnboardDefenseIn
	proc.UpdatedAt = now
	markNodeStarted(proc, "defense", &now)
	s.SaveOnboarding(proc)

	return defense, nil
}

func SubmitDefenseResult(defenseID string, result model.DefenseResult, extendMonths int, remark string) (*model.ProbationDefense, error) {
	def, ok := s.GetDefense(defenseID)
	if !ok {
		return nil, errors.New("答辩记录不存在")
	}
	if def.Result != "" {
		return nil, errors.New("答辩结果已提交")
	}
	if result == model.DefenseExtend && (extendMonths < 1 || extendMonths > 3) {
		return nil, errors.New("延长期限必须是1-3个月")
	}

	proc, ok := s.GetOnboarding(def.OnboardingID)
	if !ok {
		return nil, errors.New("入职流程不存在")
	}

	now := time.Now()
	def.Result = result
	def.ExtendMonths = extendMonths
	def.DefenseDate = &now
	def.Remark = remark
	s.SaveDefense(def)

	markNodeCompleted(proc, "defense", &now)

	switch result {
	case model.DefensePass:
		proc.Status = model.OnboardFormal
		proc.IsFormal = true
		proc.FormalDate = proc.ProbationEndDate
		markNodeCompleted(proc, "formal", proc.ProbationEndDate)
	case model.DefenseExtend:
		proc.Status = model.OnboardExtended
		newEnd := proc.ProbationEndDate.AddDate(0, extendMonths, 0)
		proc.ProbationEndDate = &newEnd
	case model.DefenseFail:
		proc.Status = model.OnboardDismissed
	}

	proc.UpdatedAt = now
	s.SaveOnboarding(proc)

	return def, nil
}

func GetTodayTodo() *model.TodayTodo {
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	var workspaceIDs []string
	var contractIDs []string
	var probationIDs []string

	oneMonthLater := today.AddDate(0, 1, 0)

	s.RangeOnboardings(func(_ string, o *model.OnboardingProcess) bool {
		if o.Status == model.OnboardAccountDone {
			workspaceIDs = append(workspaceIDs, o.EmployeeID)
		}
		if o.Status == model.OnboardContractPending {
			contractIDs = append(contractIDs, o.EmployeeID)
		}
		if (o.Status == model.OnboardProbationIn || o.Status == model.OnboardExtended) &&
			o.ProbationEndDate != nil &&
			o.ProbationEndDate.Before(oneMonthLater) &&
			o.ProbationEndDate.After(today) {
			probationIDs = append(probationIDs, o.EmployeeID)
		}
		return true
	})

	return &model.TodayTodo{
		WorkspacePrepare: &model.TodoItem{
			Date:        todayStr,
			Count:       len(workspaceIDs),
			Category:    "工位准备",
			EmployeeIDs: workspaceIDs,
		},
		ContractPending: &model.TodoItem{
			Date:        todayStr,
			Count:       len(contractIDs),
			Category:    "待签合同",
			EmployeeIDs: contractIDs,
		},
		ProbationDue: &model.TodoItem{
			Date:        todayStr,
			Count:       len(probationIDs),
			Category:    "试用期即将到期",
			EmployeeIDs: probationIDs,
		},
	}
}
