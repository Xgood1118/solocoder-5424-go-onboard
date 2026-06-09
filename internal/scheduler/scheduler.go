package scheduler

import (
	"log"
	"time"

	"hr-onboard/internal/stats"
)

func Start() {
	go func() {
		hourlyTicker := time.NewTicker(1 * time.Hour)
		defer hourlyTicker.Stop()

		for range hourlyTicker.C {
			log.Println("[scheduler] 执行小时级同步任务")
		}
	}()

	go func() {
		dailyTicker := time.NewTicker(24 * time.Hour)
		defer dailyTicker.Stop()

		for range dailyTicker.C {
			now := time.Now()
			if now.Day() == 1 {
				log.Println("[scheduler] 生成上月统计报表")
				lastMonth := now.AddDate(0, -1, 0)
				stats.GenerateMonthlyReport(lastMonth.Year(), lastMonth.Month())
			}
		}
	}()

	log.Println("[scheduler] 定时任务已启动")
}
