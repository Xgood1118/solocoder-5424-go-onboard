package model

import "time"

type OnboardingStatus string

const (
	OnboardOfferConfirmed  OnboardingStatus = "offer_confirmed"
	OnboardHealthCheckIn   OnboardingStatus = "health_check_in"
	OnboardHealthCheckDone OnboardingStatus = "health_check_done"
	OnboardBgCheckIn       OnboardingStatus = "bg_check_in"
	OnboardBgCheckDone     OnboardingStatus = "bg_check_done"
	OnboardContractPending OnboardingStatus = "contract_pending"
	OnboardContractSigned  OnboardingStatus = "contract_signed"
	OnboardAccountSetupIn  OnboardingStatus = "account_setup_in"
	OnboardAccountDone     OnboardingStatus = "account_done"
	OnboardWorkspaceDone   OnboardingStatus = "workspace_done"
	OnboardTrainingIn      OnboardingStatus = "training_in"
	OnboardTrainingDone    OnboardingStatus = "training_done"
	OnboardDeptIntroDone   OnboardingStatus = "dept_intro_done"
	OnboardManagerMeetDone OnboardingStatus = "manager_meet_done"
	OnboardProbationIn     OnboardingStatus = "probation_in"
	OnboardDefenseIn       OnboardingStatus = "defense_in"
	OnboardFormal          OnboardingStatus = "formal"
	OnboardExtended        OnboardingStatus = "extended"
	OnboardDismissed       OnboardingStatus = "dismissed"
)

type OnboardingProcess struct {
	ID                 string             `json:"id"`
	EmployeeID         string             `json:"employee_id"`
	EmployeeName       string             `json:"employee_name"`
	DepartmentID       string             `json:"department_id"`
	Level              EmployeeLevel      `json:"level"`
	Status             OnboardingStatus   `json:"status"`
	Nodes              []ProcessNode      `json:"nodes"`
	ProbationMonths    int                `json:"probation_months"`
	ProbationStartDate *time.Time         `json:"probation_start_date"`
	ProbationEndDate   *time.Time         `json:"probation_end_date"`
	OfferConfirmDate   *time.Time         `json:"offer_confirm_date"`
	FormalDate         *time.Time         `json:"formal_date"`
	IsFormal           bool               `json:"is_formal"`
	Remarks            string             `json:"remarks"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type DefenseResult string

const (
	DefensePass     DefenseResult = "pass"
	DefenseExtend   DefenseResult = "extend"
	DefenseFail     DefenseResult = "fail"
)

type DefensePanelMember struct {
	EmployeeID string `json:"employee_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
}

type ProbationDefense struct {
	ID             string             `json:"id"`
	OnboardingID   string             `json:"onboarding_id"`
	EmployeeID     string             `json:"employee_id"`
	Panel          []DefensePanelMember `json:"panel"`
	Result         DefenseResult      `json:"result"`
	ExtendMonths   int                `json:"extend_months"`
	DefenseDate    *time.Time         `json:"defense_date"`
	Remark         string             `json:"remark"`
}

type TodoItem struct {
	Date       string   `json:"date"`
	Count      int      `json:"count"`
	Category   string   `json:"category"`
	EmployeeIDs []string `json:"employee_ids"`
}

type TodayTodo struct {
	WorkspacePrepare *TodoItem `json:"workspace_prepare"`
	ContractPending  *TodoItem `json:"contract_pending"`
	ProbationDue     *TodoItem `json:"probation_due"`
}
