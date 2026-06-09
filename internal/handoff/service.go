package handoff

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

func CreateHandoff(items []model.HandoffItem) (*store.Handoff, error) {
	if len(items) == 0 {
		return nil, errors.New("交接项不能为空")
	}
	h := &store.Handoff{
		ID:    genID("hof"),
		Items: make([]model.HandoffItem, len(items)),
	}
	for i, item := range items {
		if item.ID == "" {
			item.ID = genID("hitem")
		}
		h.Items[i] = item
	}
	s.SaveHandoff(h)
	return h, nil
}

func GetHandoff(id string) (*store.Handoff, bool) {
	return s.GetHandoff(id)
}

func AddHandoffItem(handoffID string, item *model.HandoffItem) (*store.Handoff, error) {
	h, ok := s.GetHandoff(handoffID)
	if !ok {
		return nil, errors.New("交接清单不存在")
	}
	if item.ID == "" {
		item.ID = genID("hitem")
	}
	h.Items = append(h.Items, *item)
	s.SaveHandoff(h)
	return h, nil
}

func CompleteHandoffItem(handoffID, itemID string) (*store.Handoff, error) {
	h, ok := s.GetHandoff(handoffID)
	if !ok {
		return nil, errors.New("交接清单不存在")
	}
	found := false
	now := time.Now()
	for i := range h.Items {
		if h.Items[i].ID == itemID {
			h.Items[i].Completed = true
			h.Items[i].CompletedAt = &now
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("交接项不存在")
	}
	s.SaveHandoff(h)
	return h, nil
}

func RemoveHandoffItem(handoffID, itemID string) (*store.Handoff, error) {
	h, ok := s.GetHandoff(handoffID)
	if !ok {
		return nil, errors.New("交接清单不存在")
	}
	foundIdx := -1
	for i := range h.Items {
		if h.Items[i].ID == itemID {
			foundIdx = i
			break
		}
	}
	if foundIdx == -1 {
		return nil, errors.New("交接项不存在")
	}
	h.Items = append(h.Items[:foundIdx], h.Items[foundIdx+1:]...)
	s.SaveHandoff(h)
	return h, nil
}

func IsAllCompleted(handoffID string) (bool, error) {
	h, ok := s.GetHandoff(handoffID)
	if !ok {
		return false, errors.New("交接清单不存在")
	}
	if len(h.Items) == 0 {
		return true, nil
	}
	for _, item := range h.Items {
		if !item.Completed {
			return false, nil
		}
	}
	return true, nil
}
