package employee

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hr-onboard/internal/model"
)

func RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/employees")
	{
		g.POST("", createEmployee)
		g.GET("", listEmployees)
		g.GET("/:id", getEmployee)
		g.PUT("/:id", updateEmployee)
		g.DELETE("/:id", deleteEmployee)
	}

	dept := r.Group("/departments")
	{
		dept.POST("", createDepartment)
		dept.GET("", listDepartments)
		dept.GET("/:id", getDepartment)
	}

	rost := r.Group("/rosters")
	{
		rost.GET("", listRosters)
		rost.POST("/sync", syncRoster)
	}

	hol := r.Group("/holidays")
	{
		hol.GET("/:employee_id", getHolidayBalance)
		hol.POST("/sync", syncHolidays)
	}
}

func createEmployee(c *gin.Context) {
	var emp model.EmployeeProfile
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := CreateEmployee(&emp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func listEmployees(c *gin.Context) {
	list := ListEmployees()
	c.JSON(http.StatusOK, list)
}

func getEmployee(c *gin.Context) {
	id := c.Param("id")
	emp, ok := GetEmployee(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func updateEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp model.EmployeeProfile
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := UpdateEmployee(id, &emp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := DeleteEmployee(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func createDepartment(c *gin.Context) {
	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := CreateDepartment(&dept)
	c.JSON(http.StatusOK, result)
}

func listDepartments(c *gin.Context) {
	list := ListDepartments()
	c.JSON(http.StatusOK, list)
}

func getDepartment(c *gin.Context) {
	id := c.Param("id")
	dept, ok := GetDepartment(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func listRosters(c *gin.Context) {
	list := ListRosters()
	c.JSON(http.StatusOK, list)
}

type syncRosterReq struct {
	List []*model.EmployeeRoster `json:"list"`
}

func syncRoster(c *gin.Context) {
	var req syncRosterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SyncRoster(req.List)
	c.JSON(http.StatusOK, gin.H{"message": "synced", "count": len(req.List)})
}

func getHolidayBalance(c *gin.Context) {
	empID := c.Param("employee_id")
	h, ok := GetHolidayBalance(empID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, h)
}

type syncHolidayReq struct {
	List []*model.HolidayBalance `json:"list"`
}

func syncHolidays(c *gin.Context) {
	var req syncHolidayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SyncHolidayBalances(req.List)
	c.JSON(http.StatusOK, gin.H{"message": "synced", "count": len(req.List)})
}
