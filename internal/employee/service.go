package employee

import (
	"encoding/base64"
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

// TODO: 当前使用 base64 占位加密，密钥管理方案待定后替换为 AES-GCM 等安全加密
func encryptCardNo(plain string) string {
	if plain == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(plain))
}

// TODO: 对应解密函数，密钥管理方案待定后替换
func decryptCardNo(encoded string) string {
	if encoded == "" {
		return ""
	}
	dec, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return encoded
	}
	return string(dec)
}

func CreateEmployee(emp *model.EmployeeProfile) (*model.EmployeeProfile, error) {
	if emp.Name == "" {
		return nil, errors.New("姓名不能为空")
	}
	if emp.ID == "" {
		emp.ID = genID("emp")
	}
	now := time.Now()
	emp.CreatedAt = now
	emp.UpdatedAt = now
	emp.IsFormal = false
	emp.BankCard.CardNo = encryptCardNo(emp.BankCard.CardNo)
	s.SaveEmployee(emp)
	emp.BankCard.CardNo = decryptCardNo(emp.BankCard.CardNo)
	return emp, nil
}

func GetEmployee(id string) (*model.EmployeeProfile, bool) {
	emp, ok := s.GetEmployee(id)
	if ok {
		emp.BankCard.CardNo = decryptCardNo(emp.BankCard.CardNo)
	}
	return emp, ok
}

func UpdateEmployee(id string, upd *model.EmployeeProfile) (*model.EmployeeProfile, error) {
	emp, ok := s.GetEmployee(id)
	if !ok {
		return nil, errors.New("员工不存在")
	}
	if upd.Name != "" {
		emp.Name = upd.Name
	}
	if upd.Phone != "" {
		emp.Phone = upd.Phone
	}
	if upd.IDCard.IDNumber != "" {
		emp.IDCard = upd.IDCard
	}
	if upd.EmergencyContact != "" {
		emp.EmergencyContact = upd.EmergencyContact
	}
	if upd.EmergencyPhone != "" {
		emp.EmergencyPhone = upd.EmergencyPhone
	}
	if upd.BankCard.CardNo != "" {
		emp.BankCard = upd.BankCard
		emp.BankCard.CardNo = encryptCardNo(upd.BankCard.CardNo)
	}
	if upd.PhotoBase64 != "" {
		emp.PhotoBase64 = upd.PhotoBase64
	}
	if upd.ESignatureBase64 != "" {
		emp.ESignatureBase64 = upd.ESignatureBase64
	}
	emp.SpecialNeeds = upd.SpecialNeeds
	if upd.DepartmentID != "" {
		emp.DepartmentID = upd.DepartmentID
	}
	if upd.Position != "" {
		emp.Position = upd.Position
	}
	if upd.Level != "" {
		emp.Level = upd.Level
	}
	if upd.DirectManagerID != "" {
		emp.DirectManagerID = upd.DirectManagerID
	}
	if upd.BaseSalary > 0 {
		emp.BaseSalary = upd.BaseSalary
	}
	if upd.OnboardDate != nil {
		emp.OnboardDate = upd.OnboardDate
	}
	emp.UpdatedAt = time.Now()
	s.SaveEmployee(emp)
	return emp, nil
}

func DeleteEmployee(id string) error {
	_, ok := s.GetEmployee(id)
	if !ok {
		return errors.New("员工不存在")
	}
	s.DeleteEmployee(id)
	return nil
}

func ListEmployees() []*model.EmployeeProfile {
	var list []*model.EmployeeProfile
	s.RangeEmployees(func(_ string, emp *model.EmployeeProfile) bool {
		emp.BankCard.CardNo = decryptCardNo(emp.BankCard.CardNo)
		list = append(list, emp)
		return true
	})
	return list
}

func CountEmployees() int {
	return s.EmployeeCount()
}

func SyncRoster(list []*model.EmployeeRoster) {
	for _, r := range list {
		s.SaveRoster(r)
	}
}

func GetRoster(id string) (*model.EmployeeRoster, bool) {
	return s.GetRoster(id)
}

func ListRosters() []*model.EmployeeRoster {
	var list []*model.EmployeeRoster
	s.RangeRosters(func(_ string, r *model.EmployeeRoster) bool {
		list = append(list, r)
		return true
	})
	return list
}

func CreateDepartment(dept *model.Department) *model.Department {
	if dept.ID == "" {
		dept.ID = genID("dept")
	}
	s.SaveDepartment(dept)
	return dept
}

func GetDepartment(id string) (*model.Department, bool) {
	return s.GetDepartment(id)
}

func ListDepartments() []*model.Department {
	var list []*model.Department
	s.RangeDepartments(func(_ string, d *model.Department) bool {
		list = append(list, d)
		return true
	})
	return list
}

func SyncHolidayBalances(list []*model.HolidayBalance) {
	for _, h := range list {
		s.SaveHolidayBalance(h)
	}
}

func GetHolidayBalance(empID string) (*model.HolidayBalance, bool) {
	return s.GetHolidayBalance(empID)
}
