package approval

import (
	"errors"
	"fmt"
	"time"

	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func ValidateFlowNoCycle(flow *model.ApprovalFlowConfig) error {
	if flow == nil || len(flow.NodeKeys) == 0 {
		return errors.New("审批流配置为空")
	}

	nodeSet := make(map[string]bool, len(flow.NodeKeys))
	for _, k := range flow.NodeKeys {
		nodeSet[k] = true
	}

	adj := make(map[string][]string)
	for i := 0; i < len(flow.NodeKeys)-1; i++ {
		from := flow.NodeKeys[i]
		to := flow.NodeKeys[i+1]
		adj[from] = append(adj[from], to)
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		for _, next := range adj[node] {
			if !visited[next] {
				if dfs(next) {
					return true
				}
			} else if recStack[next] {
				return true
			}
		}
		recStack[node] = false
		return false
	}

	for _, node := range flow.NodeKeys {
		if !visited[node] {
			if dfs(node) {
				return errors.New("审批链存在循环引用")
			}
		}
	}

	return nil
}

func CreateFlow(flow *model.ApprovalFlowConfig) error {
	if err := ValidateFlowNoCycle(flow); err != nil {
		return err
	}
	if flow.ID == "" {
		flow.ID = genID("flow")
	}
	s.SaveApprovalFlow(flow)
	return nil
}

func GetFlow(id string) (*model.ApprovalFlowConfig, bool) {
	return s.GetApprovalFlow(id)
}

func ListFlows() []*model.ApprovalFlowConfig {
	var flows []*model.ApprovalFlowConfig
	s.RangeApprovalFlows(func(_ string, f *model.ApprovalFlowConfig) bool {
		flows = append(flows, f)
		return true
	})
	return flows
}

func UpdateFlow(flow *model.ApprovalFlowConfig) error {
	if err := ValidateFlowNoCycle(flow); err != nil {
		return err
	}
	s.SaveApprovalFlow(flow)
	return nil
}

func DeleteFlow(id string) {
	s.SaveApprovalFlow(&model.ApprovalFlowConfig{ID: id, NodeMap: map[string]model.ApprovalNode{}})
}

type StartApprovalRequest struct {
	FlowID      string
	Subject     string
	InitiatorID string
}

func StartApproval(req StartApprovalRequest) (*model.ApprovalInstance, error) {
	flow, ok := s.GetApprovalFlow(req.FlowID)
	if !ok {
		return nil, errors.New("审批流不存在")
	}

	if err := ValidateFlowNoCycle(flow); err != nil {
		return nil, err
	}

	var nodes []model.ApprovalNode
	for _, nk := range flow.NodeKeys {
		n, ok := flow.NodeMap[nk]
		if !ok {
			return nil, fmt.Errorf("节点 %s 配置缺失", nk)
		}
		n.Status = model.ApprovalPending
		n.ApprovedAt = nil
		nodes = append(nodes, n)
	}

	if len(nodes) > 0 {
		nodes[0].Status = model.ApprovalPending
	}

	inst := &model.ApprovalInstance{
		ID:          genID("appr"),
		FlowID:      req.FlowID,
		Subject:     req.Subject,
		InitiatorID: req.InitiatorID,
		CurrentNode: flow.NodeKeys[0],
		Nodes:       nodes,
		Status:      model.ApprovalPending,
		CreatedAt:   time.Now(),
	}

	s.SaveApprovalInstance(inst)
	return inst, nil
}

func GetApprovalInstance(id string) (*model.ApprovalInstance, bool) {
	return s.GetApprovalInstance(id)
}

func Approve(instanceID, approverID, remark string) (*model.ApprovalInstance, error) {
	inst, ok := s.GetApprovalInstance(instanceID)
	if !ok {
		return nil, errors.New("审批实例不存在")
	}
	if inst.Status != model.ApprovalPending {
		return nil, errors.New("审批已结束")
	}

	idx := -1
	for i, n := range inst.Nodes {
		if n.NodeKey == inst.CurrentNode {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, errors.New("当前节点无效")
	}

	currentNode := &inst.Nodes[idx]
	if currentNode.ApproverID != approverID {
		return nil, errors.New("无权审批此节点")
	}
	if currentNode.Status == model.ApprovalApproved {
		return nil, errors.New("节点已审批")
	}

	now := time.Now()
	currentNode.Status = model.ApprovalApproved
	currentNode.ApprovedAt = &now
	currentNode.Remark = remark

	if idx < len(inst.Nodes)-1 {
		inst.CurrentNode = inst.Nodes[idx+1].NodeKey
	} else {
		inst.Status = model.ApprovalApproved
		inst.FinishedAt = &now
		inst.CurrentNode = ""
	}

	s.SaveApprovalInstance(inst)
	return inst, nil
}

func Reject(instanceID, approverID, remark string) (*model.ApprovalInstance, error) {
	inst, ok := s.GetApprovalInstance(instanceID)
	if !ok {
		return nil, errors.New("审批实例不存在")
	}
	if inst.Status != model.ApprovalPending {
		return nil, errors.New("审批已结束")
	}

	idx := -1
	for i, n := range inst.Nodes {
		if n.NodeKey == inst.CurrentNode {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, errors.New("当前节点无效")
	}

	currentNode := &inst.Nodes[idx]
	if currentNode.ApproverID != approverID {
		return nil, errors.New("无权审批此节点")
	}

	now := time.Now()
	currentNode.Status = model.ApprovalRejected
	currentNode.ApprovedAt = &now
	currentNode.Remark = remark

	inst.Status = model.ApprovalRejected
	inst.FinishedAt = &now

	s.SaveApprovalInstance(inst)
	return inst, nil
}

func ListInstances() []*model.ApprovalInstance {
	var list []*model.ApprovalInstance
	s.RangeApprovalInstances(func(_ string, inst *model.ApprovalInstance) bool {
		list = append(list, inst)
		return true
	})
	return list
}
