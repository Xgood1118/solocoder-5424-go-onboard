package model

import "time"

type OffboardingStatus string

const (
	OffboardApplied        OffboardingStatus = "applied"
	OffboardManagerConfirmed OffboardingStatus = "manager_confirmed"
	OffboardDeptApproved   OffboardingStatus = "dept_approved"
	OffboardHROnApproved   OffboardingStatus = "hr_approved"
	OffboardHandoffIn      OffboardingStatus = "handoff_in"
	OffboardAssetReturnIn  OffboardingStatus = "asset_return_in"
	OffboardFinanceSettled OffboardingStatus = "finance_settled"
	OffboardInterviewDone  OffboardingStatus = "interview_done"
	OffboardCertIssued     OffboardingStatus = "cert_issued"
	OffboardSocialStopped  OffboardingStatus = "social_stopped"
	OffboardCompleted      OffboardingStatus = "completed"
)

type FinancialSettlement struct {
	RemainingSalary float64 `json:"remaining_salary"`
	AnnualLeaveConvert float64 `json:"annual_leave_convert"`
	StockBuyback    float64 `json:"stock_buyback"`
	AssetDeduction  float64 `json:"asset_deduction"`
	UnpaidAssetsDebt float64 `json:"unpaid_assets_debt"`
	TotalPayable    float64 `json:"total_payable"`
}

type OffboardingProcess struct {
	ID                string              `json:"id"`
	EmployeeID        string              `json:"employee_id"`
	EmployeeName      string              `json:"employee_name"`
	DepartmentID      string              `json:"department_id"`
	Level             EmployeeLevel       `json:"level"`
	Status            OffboardingStatus   `json:"status"`
	ApplyDate         time.Time           `json:"apply_date"`
	LastWorkingDate   time.Time           `json:"last_working_date"`
	NoticeDays        int                 `json:"notice_days"`
	Nodes             []ProcessNode       `json:"nodes"`
	ApprovalID        string              `json:"approval_id"`
	HandoffID         string              `json:"handoff_id"`
	ReturnableAssets  []ReturnableAsset   `json:"returnable_assets"`
	Settlement        *FinancialSettlement `json:"settlement"`
	InterviewRemark   string              `json:"interview_remark"`
	CertificateURL    string              `json:"certificate_url"`
	SocialStopDate    *time.Time          `json:"social_stop_date"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type ReminderMessage struct {
	ID         string    `json:"id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	ProcessID  string    `json:"process_id"`
	NodeKey    string    `json:"node_key"`
	CreatedAt  time.Time `json:"created_at"`
	Read       bool      `json:"read"`
}

type PendingDebtItem struct {
	ID         string    `json:"id"`
	EmployeeID string    `json:"employee_id"`
	EmployeeName string   `json:"employee_name"`
	Amount     float64   `json:"amount"`
	AssetIDs   []string  `json:"asset_ids"`
	OffboardID string    `json:"offboard_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
