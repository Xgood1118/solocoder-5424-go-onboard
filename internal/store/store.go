package store

import (
	"sync"

	"hr-onboard/internal/model"
)

type Store struct {
	employees        sync.Map
	departments      sync.Map
	rosters          sync.Map
	holidayBalances  sync.Map
	assets           sync.Map
	onboardings      sync.Map
	offboardings     sync.Map
	approvalFlows    sync.Map
	approvalInstances sync.Map
	handoffs         sync.Map
	auditLogs        sync.Map
	monthlyReports   sync.Map
	reminders        sync.Map
	pendingDebts     sync.Map
	defenses         sync.Map
}

var globalStore = &Store{}

func Get() *Store {
	return globalStore
}

func (s *Store) SaveEmployee(emp *model.EmployeeProfile) {
	s.employees.Store(emp.ID, emp)
}

func (s *Store) GetEmployee(id string) (*model.EmployeeProfile, bool) {
	v, ok := s.employees.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.EmployeeProfile), true
}

func (s *Store) DeleteEmployee(id string) {
	s.employees.Delete(id)
}

func (s *Store) RangeEmployees(fn func(id string, emp *model.EmployeeProfile) bool) {
	s.employees.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.EmployeeProfile))
	})
}

func (s *Store) EmployeeCount() int {
	count := 0
	s.employees.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *Store) SaveDepartment(dept *model.Department) {
	s.departments.Store(dept.ID, dept)
}

func (s *Store) GetDepartment(id string) (*model.Department, bool) {
	v, ok := s.departments.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.Department), true
}

func (s *Store) RangeDepartments(fn func(id string, dept *model.Department) bool) {
	s.departments.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.Department))
	})
}

func (s *Store) SaveRoster(r *model.EmployeeRoster) {
	s.rosters.Store(r.EmployeeID, r)
}

func (s *Store) GetRoster(id string) (*model.EmployeeRoster, bool) {
	v, ok := s.rosters.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.EmployeeRoster), true
}

func (s *Store) RangeRosters(fn func(id string, r *model.EmployeeRoster) bool) {
	s.rosters.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.EmployeeRoster))
	})
}

func (s *Store) SaveHolidayBalance(h *model.HolidayBalance) {
	s.holidayBalances.Store(h.EmployeeID, h)
}

func (s *Store) GetHolidayBalance(empID string) (*model.HolidayBalance, bool) {
	v, ok := s.holidayBalances.Load(empID)
	if !ok {
		return nil, false
	}
	return v.(*model.HolidayBalance), true
}

func (s *Store) SaveAsset(a *model.Asset) {
	s.assets.Store(a.ID, a)
}

func (s *Store) GetAsset(id string) (*model.Asset, bool) {
	v, ok := s.assets.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.Asset), true
}

func (s *Store) RangeAssets(fn func(id string, a *model.Asset) bool) {
	s.assets.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.Asset))
	})
}

func (s *Store) SaveOnboarding(o *model.OnboardingProcess) {
	s.onboardings.Store(o.ID, o)
}

func (s *Store) GetOnboarding(id string) (*model.OnboardingProcess, bool) {
	v, ok := s.onboardings.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.OnboardingProcess), true
}

func (s *Store) RangeOnboardings(fn func(id string, o *model.OnboardingProcess) bool) {
	s.onboardings.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.OnboardingProcess))
	})
}

func (s *Store) SaveOffboarding(o *model.OffboardingProcess) {
	s.offboardings.Store(o.ID, o)
}

func (s *Store) GetOffboarding(id string) (*model.OffboardingProcess, bool) {
	v, ok := s.offboardings.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.OffboardingProcess), true
}

func (s *Store) RangeOffboardings(fn func(id string, o *model.OffboardingProcess) bool) {
	s.offboardings.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.OffboardingProcess))
	})
}

func (s *Store) SaveApprovalFlow(f *model.ApprovalFlowConfig) {
	s.approvalFlows.Store(f.ID, f)
}

func (s *Store) GetApprovalFlow(id string) (*model.ApprovalFlowConfig, bool) {
	v, ok := s.approvalFlows.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.ApprovalFlowConfig), true
}

func (s *Store) RangeApprovalFlows(fn func(id string, f *model.ApprovalFlowConfig) bool) {
	s.approvalFlows.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.ApprovalFlowConfig))
	})
}

func (s *Store) SaveApprovalInstance(inst *model.ApprovalInstance) {
	s.approvalInstances.Store(inst.ID, inst)
}

func (s *Store) GetApprovalInstance(id string) (*model.ApprovalInstance, bool) {
	v, ok := s.approvalInstances.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.ApprovalInstance), true
}

func (s *Store) RangeApprovalInstances(fn func(id string, inst *model.ApprovalInstance) bool) {
	s.approvalInstances.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.ApprovalInstance))
	})
}

type Handoff struct {
	ID    string            `json:"id"`
	Items []model.HandoffItem `json:"items"`
}

func (s *Store) SaveHandoff(h *Handoff) {
	s.handoffs.Store(h.ID, h)
}

func (s *Store) GetHandoff(id string) (*Handoff, bool) {
	v, ok := s.handoffs.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*Handoff), true
}

func (s *Store) SaveAuditLog(log *model.AuditLog) {
	s.auditLogs.Store(log.ID, log)
}

func (s *Store) RangeAuditLogs(fn func(id string, log *model.AuditLog) bool) {
	s.auditLogs.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.AuditLog))
	})
}

func (s *Store) SaveMonthlyReport(r *model.MonthlyReport) {
	s.monthlyReports.Store(r.Month, r)
}

func (s *Store) GetMonthlyReport(month string) (*model.MonthlyReport, bool) {
	v, ok := s.monthlyReports.Load(month)
	if !ok {
		return nil, false
	}
	return v.(*model.MonthlyReport), true
}

func (s *Store) RangeMonthlyReports(fn func(month string, r *model.MonthlyReport) bool) {
	s.monthlyReports.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.MonthlyReport))
	})
}

func (s *Store) SaveReminder(r *model.ReminderMessage) {
	s.reminders.Store(r.ID, r)
}

func (s *Store) RangeReminders(fn func(id string, r *model.ReminderMessage) bool) {
	s.reminders.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.ReminderMessage))
	})
}

func (s *Store) SavePendingDebt(d *model.PendingDebtItem) {
	s.pendingDebts.Store(d.ID, d)
}

func (s *Store) RangePendingDebts(fn func(id string, d *model.PendingDebtItem) bool) {
	s.pendingDebts.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.PendingDebtItem))
	})
}

func (s *Store) SaveDefense(d *model.ProbationDefense) {
	s.defenses.Store(d.ID, d)
}

func (s *Store) GetDefense(id string) (*model.ProbationDefense, bool) {
	v, ok := s.defenses.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*model.ProbationDefense), true
}

func (s *Store) RangeDefenses(fn func(id string, d *model.ProbationDefense) bool) {
	s.defenses.Range(func(key, value interface{}) bool {
		return fn(key.(string), value.(*model.ProbationDefense))
	})
}
