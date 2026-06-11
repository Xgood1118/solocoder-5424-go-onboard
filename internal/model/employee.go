package model

import "time"

type NodeStatus string

const (
	StatusPending   NodeStatus = "pending"
	StatusInProgress NodeStatus = "in_progress"
	StatusCompleted NodeStatus = "completed"
	StatusSkipped   NodeStatus = "skipped"
	StatusOverdue   NodeStatus = "overdue"
)

type EmployeeLevel string

const (
	LevelJunior EmployeeLevel = "junior"
	LevelMid    EmployeeLevel = "mid"
	LevelSenior EmployeeLevel = "senior"
	LevelExecutive EmployeeLevel = "executive"
)

type ProcessNode struct {
	NodeKey        string     `json:"node_key"`
	NodeName       string     `json:"node_name"`
	Assignee       string     `json:"assignee"`
	PlannedStartAt *time.Time `json:"planned_start_at"`
	PlannedEndAt   *time.Time `json:"planned_end_at"`
	ActualStartAt  *time.Time `json:"actual_start_at"`
	ActualEndAt    *time.Time `json:"actual_end_at"`
	Status         NodeStatus `json:"status"`
	Remark         string     `json:"remark"`
}

type IDCard struct {
	IDNumber    string   `json:"id_number"`
	RegisteredResidence string `json:"registered_residence"`
}

type BankCard struct {
	BankName string `json:"bank_name"`
	CardNo   string `json:"card_no"`
}

// 银行卡号使用 AES-256-GCM 加密存储，密钥通过 pkg/crypto 包管理
// TODO: 密钥应从环境变量 HR_AES_KEY 或密钥管理服务加载，当前硬编码仅用于开发

type SpecialNeeds struct {
	DietRestrictions string `json:"diet_restrictions"`
	Religion         string `json:"religion"`
	Accessibility    string `json:"accessibility"`
	Other            string `json:"other"`
}

type EmployeeProfile struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Phone           string       `json:"phone"`
	IDCard          IDCard       `json:"id_card"`
	EmergencyContact string      `json:"emergency_contact"`
	EmergencyPhone  string       `json:"emergency_phone"`
	BankCard        BankCard     `json:"bank_card"`
	PhotoBase64     string       `json:"photo_base64"`
	ESignatureBase64 string      `json:"e_signature_base64"`
	SpecialNeeds    SpecialNeeds `json:"special_needs"`
	DepartmentID    string       `json:"department_id"`
	Position        string       `json:"position"`
	Level           EmployeeLevel `json:"level"`
	DirectManagerID string       `json:"direct_manager_id"`
	BaseSalary      float64      `json:"base_salary"`
	OnboardDate     *time.Time   `json:"onboard_date"`
	ProbationEndDate *time.Time  `json:"probation_end_date"`
	IsFormal        bool         `json:"is_formal"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type Department struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
	ManagerID string `json:"manager_id"`
	HRBPID   string `json:"hrbp_id"`
}

type EmployeeRoster struct {
	EmployeeID string       `json:"employee_id"`
	Name       string       `json:"name"`
	DepartmentID string     `json:"department_id"`
	Position   string       `json:"position"`
	Level      EmployeeLevel `json:"level"`
	Status     string       `json:"status"`
	OnboardDate *time.Time  `json:"onboard_date"`
}

type HolidayBalance struct {
	EmployeeID  string  `json:"employee_id"`
	AnnualLeave float64 `json:"annual_leave"`
	SickLeave   float64 `json:"sick_leave"`
}
