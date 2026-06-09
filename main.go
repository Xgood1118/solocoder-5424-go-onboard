package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/approval"
	"hr-onboard/internal/asset"
	"hr-onboard/internal/audit"
	"hr-onboard/internal/employee"
	"hr-onboard/internal/handoff"
	"hr-onboard/internal/offboarding"
	"hr-onboard/internal/onboarding"
	"hr-onboard/internal/scheduler"
	"hr-onboard/internal/stats"
	"hr-onboard/internal/store"
	"hr-onboard/pkg/config"
)

func main() {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		employee.RegisterRoutes(api)
		onboarding.RegisterRoutes(api)
		offboarding.RegisterRoutes(api)
		asset.RegisterRoutes(api)
		handoff.RegisterRoutes(api)
		approval.RegisterRoutes(api)
		stats.RegisterRoutes(api)
		audit.RegisterRoutes(api)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	port := config.GetPort()
	empCount := store.Get().EmployeeCount()

	log.Println("==================================")
	log.Printf("服务启动端口: %s", port)
	log.Printf("当前员工档案数: %d 条", empCount)
	log.Println("==================================")

	scheduler.Start()

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
