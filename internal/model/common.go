package model

import "time"

type AssetType string

const (
	AssetTypeLaptop   AssetType = "laptop"
	AssetTypeBadge    AssetType = "badge"
	AssetTypeAccessCard AssetType = "access_card"
	AssetTypeKey      AssetType = "key"
	AssetTypeBook     AssetType = "book"
)

type Asset struct {
	ID           string    `json:"id"`
	Type         AssetType `json:"type"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	OriginalValue float64  `json:"original_value"`
	PurchaseDate time.Time `json:"purchase_date"`
	HolderID     string    `json:"holder_id"`
	Status       string    `json:"status"`
}

type ReturnableAsset struct {
	AssetID   string    `json:"asset_id"`
	AssetType AssetType `json:"asset_type"`
	AssetName string    `json:"asset_name"`
	Returned  bool      `json:"returned"`
	ReturnedAt *time.Time `json:"returned_at"`
	DepreciatedValue float64 `json:"depreciated_value"`
}

type HandoffItem struct {
	ID          string `json:"id"`
	GiverID     string `json:"giver_id"`
	ReceiverID  string `json:"receiver_id"`
	Content     string `json:"content"`
	Completed   bool   `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
}

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

type ApprovalNode struct {
	NodeKey      string         `json:"node_key"`
	NodeName     string         `json:"node_name"`
	ApproverID   string         `json:"approver_id"`
	ApproverName string         `json:"approver_name"`
	Status       ApprovalStatus `json:"status"`
	ApprovedAt   *time.Time     `json:"approved_at"`
	Remark       string         `json:"remark"`
}

type ApprovalFlowConfig struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	NodeKeys  []string         `json:"node_keys"`
	NodeMap   map[string]ApprovalNode `json:"node_map"`
}

type ApprovalInstance struct {
	ID         string           `json:"id"`
	FlowID     string           `json:"flow_id"`
	Subject    string           `json:"subject"`
	InitiatorID string          `json:"initiator_id"`
	CurrentNode string          `json:"current_node"`
	Nodes      []ApprovalNode   `json:"nodes"`
	Status     ApprovalStatus   `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	FinishedAt *time.Time       `json:"finished_at"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	Operator  string    `json:"operator"`
	Action    string    `json:"action"`
	TargetID  string    `json:"target_id"`
	TargetType string   `json:"target_type"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
	Detail    string    `json:"detail"`
}

type MonthlyReport struct {
	Month           string  `json:"month"`
	OnboardCount    int     `json:"onboard_count"`
	PlannedOnboard  int     `json:"planned_onboard"`
	OnboardRate     float64 `json:"onboard_rate"`
	OffboardCount   int     `json:"offboard_count"`
	StartHeadcount  int     `json:"start_headcount"`
	OffboardRate    float64 `json:"offboard_rate"`
	ProbationPass   int     `json:"probation_pass"`
	ProbationTotal  int     `json:"probation_total"`
	ProbationPassRate float64 `json:"probation_pass_rate"`
	AvgOnboardDays  float64 `json:"avg_onboard_days"`
}
