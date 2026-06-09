package audit

import (
	"fmt"
	"time"

	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func AddLog(operator, action, targetID, targetType, reason, detail string) {
	log := &model.AuditLog{
		ID:         genID("audit"),
		Operator:   operator,
		Action:     action,
		TargetID:   targetID,
		TargetType: targetType,
		Reason:     reason,
		Detail:     detail,
		Timestamp:  time.Now(),
	}
	s.SaveAuditLog(log)
}

func ListLogs(targetType, targetID string) []*model.AuditLog {
	var list []*model.AuditLog
	s.RangeAuditLogs(func(_ string, log *model.AuditLog) bool {
		if targetType != "" && log.TargetType != targetType {
			return true
		}
		if targetID != "" && log.TargetID != targetID {
			return true
		}
		list = append(list, log)
		return true
	})
	return list
}
