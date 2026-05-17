package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	db *gorm.DB
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{db: db}
}

// Check godoc
// @Summary      Health check
// @Description  Returns the health status of the API and database
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /health [get]
func (hc *HealthController) Check(c *gin.Context) {
	response := gin.H{
		"status":  "ok",
		"message": "Service is up and running",
	}

	// Check database connectivity
	var result int
	if err := hc.db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		response["status"] = "degraded"
		response["db_ok"] = false
		response["db_error"] = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response["db_ok"] = true
	c.JSON(http.StatusOK, response)
}
