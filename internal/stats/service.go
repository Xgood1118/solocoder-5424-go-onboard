package stats

import (
	"math"
	"time"

	"hr-onboard/internal/model"
	"hr-onboard/internal/store"
)

var s = store.Get()

func GenerateMonthlyReport(year int, month time.Month) *model.MonthlyReport {
	monthStr := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Format("2006-01")
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	var onboardCount int
	var offboardCount int
	var probationPass int
	var probationTotal int
	var totalOnboardDays float64
	var onboardForCalc int

	s.RangeOnboardings(func(_ string, o *model.OnboardingProcess) bool {
		if o.FormalDate != nil && o.FormalDate.After(monthStart) && o.FormalDate.Before(nextMonthStart) {
			onboardCount++
		}
		if o.Status == model.OnboardFormal && o.OfferConfirmDate != nil && o.FormalDate != nil {
			days := o.FormalDate.Sub(*o.OfferConfirmDate).Hours() / 24
			if days > 0 {
				totalOnboardDays += days
				onboardForCalc++
			}
		}
		if o.ProbationEndDate != nil &&
			o.ProbationEndDate.After(monthStart) &&
			o.ProbationEndDate.Before(nextMonthStart) {
			probationTotal++
			if o.Status == model.OnboardFormal &&
				o.FormalDate != nil &&
				o.FormalDate.Equal(*o.ProbationEndDate) {
				probationPass++
			}
		}
		return true
	})

	s.RangeOffboardings(func(_ string, o *model.OffboardingProcess) bool {
		if o.LastWorkingDate.After(monthStart) && o.LastWorkingDate.Before(nextMonthStart) {
			if o.Status == model.OffboardCompleted {
				offboardCount++
			}
		}
		return true
	})

	startHeadcount := 0
	s.RangeEmployees(func(_ string, e *model.EmployeeProfile) bool {
		if e.OnboardDate != nil && e.OnboardDate.Before(monthStart) {
			if !e.IsFormal {
			} else {
				startHeadcount++
			}
		}
		return true
	})

	plannedOnboard := 20
	if plannedOnboard == 0 {
		plannedOnboard = 1
	}

	onboardRate := float64(onboardCount) / float64(plannedOnboard)
	offboardRate := 0.0
	if startHeadcount > 0 {
		offboardRate = float64(offboardCount) / float64(startHeadcount)
	}
	probationPassRate := 0.0
	if probationTotal > 0 {
		probationPassRate = float64(probationPass) / float64(probationTotal)
	}
	avgOnboardDays := 0.0
	if onboardForCalc > 0 {
		avgOnboardDays = totalOnboardDays / float64(onboardForCalc)
	}

	report := &model.MonthlyReport{
		Month:              monthStr,
		OnboardCount:       onboardCount,
		PlannedOnboard:     plannedOnboard,
		OnboardRate:        math.Round(onboardRate*10000) / 10000,
		OffboardCount:      offboardCount,
		StartHeadcount:     startHeadcount,
		OffboardRate:       math.Round(offboardRate*10000) / 10000,
		ProbationPass:      probationPass,
		ProbationTotal:     probationTotal,
		ProbationPassRate:  math.Round(probationPassRate*10000) / 10000,
		AvgOnboardDays:     math.Round(avgOnboardDays*100) / 100,
	}

	s.SaveMonthlyReport(report)
	return report
}

func GetMonthlyReport(month string) (*model.MonthlyReport, bool) {
	return s.GetMonthlyReport(month)
}

func ListMonthlyReports() []*model.MonthlyReport {
	var list []*model.MonthlyReport
	s.RangeMonthlyReports(func(_ string, r *model.MonthlyReport) bool {
		list = append(list, r)
		return true
	})
	return list
}

func GetCurrentMonthStats() map[string]interface{} {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	var onboardThisMonth int
	var offboardThisMonth int
	var inProbation int
	var pendingApproval int

	s.RangeOnboardings(func(_ string, o *model.OnboardingProcess) bool {
		if o.FormalDate != nil && o.FormalDate.After(monthStart) {
			onboardThisMonth++
		}
		if o.Status == model.OnboardProbationIn || o.Status == model.OnboardExtended {
			inProbation++
		}
		return true
	})

	s.RangeOffboardings(func(_ string, o *model.OffboardingProcess) bool {
		if o.LastWorkingDate.After(monthStart) {
			offboardThisMonth++
		}
		if o.Status != model.OffboardCompleted {
			pendingApproval++
		}
		return true
	})

	totalEmployees := s.EmployeeCount()

	return map[string]interface{}{
		"total_employees":   totalEmployees,
		"onboard_this_month": onboardThisMonth,
		"offboard_this_month": offboardThisMonth,
		"in_probation":      inProbation,
		"pending_offboarding": pendingApproval,
	}
}
