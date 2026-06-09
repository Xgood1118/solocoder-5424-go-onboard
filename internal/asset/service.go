package asset

import (
	"errors"
	"fmt"
	"math"
	"time"

	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func CreateAsset(a *model.Asset) (*model.Asset, error) {
	if a.Code == "" || a.Type == "" {
		return nil, errors.New("资产编号和类型不能为空")
	}
	if a.ID == "" {
		a.ID = genID("ast")
	}
	if a.Status == "" {
		a.Status = "available"
	}
	s.SaveAsset(a)
	return a, nil
}

func GetAsset(id string) (*model.Asset, bool) {
	return s.GetAsset(id)
}

func UpdateAsset(id string, upd *model.Asset) (*model.Asset, error) {
	a, ok := s.GetAsset(id)
	if !ok {
		return nil, errors.New("资产不存在")
	}
	if upd.Type != "" {
		a.Type = upd.Type
	}
	if upd.Code != "" {
		a.Code = upd.Code
	}
	if upd.Name != "" {
		a.Name = upd.Name
	}
	if upd.OriginalValue > 0 {
		a.OriginalValue = upd.OriginalValue
	}
	if !upd.PurchaseDate.IsZero() {
		a.PurchaseDate = upd.PurchaseDate
	}
	if upd.HolderID != "" {
		a.HolderID = upd.HolderID
	}
	if upd.Status != "" {
		a.Status = upd.Status
	}
	s.SaveAsset(a)
	return a, nil
}

func ListAssets(assetType model.AssetType, holderID string) []*model.Asset {
	var list []*model.Asset
	s.RangeAssets(func(_ string, a *model.Asset) bool {
		if assetType != "" && a.Type != assetType {
			return true
		}
		if holderID != "" && a.HolderID != holderID {
			return true
		}
		list = append(list, a)
		return true
	})
	return list
}

func AssignAsset(assetID, holderID string) error {
	a, ok := s.GetAsset(assetID)
	if !ok {
		return errors.New("资产不存在")
	}
	if a.HolderID != "" {
		return errors.New("资产已被领用")
	}
	a.HolderID = holderID
	a.Status = "in_use"
	s.SaveAsset(a)
	return nil
}

func ReturnAsset(assetID string) error {
	a, ok := s.GetAsset(assetID)
	if !ok {
		return errors.New("资产不存在")
	}
	a.HolderID = ""
	a.Status = "available"
	s.SaveAsset(a)
	return nil
}

func CalculateDepreciation(originalValue float64, purchaseDate time.Time) float64 {
	years := time.Since(purchaseDate).Hours() / 24 / 365
	depreciationRate := 0.2
	depreciated := originalValue * math.Pow(1-depreciationRate, years)
	if depreciated < 0 {
		depreciated = 0
	}
	return math.Round(depreciated*100) / 100
}

func GetEmployeeAssets(employeeID string) []*model.Asset {
	return ListAssets("", employeeID)
}
