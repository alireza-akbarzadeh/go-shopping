package controllers

import (
	"github.com/alireza-akbarzadeh/shopping-platform/views"
	"github.com/gin-gonic/gin"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

// LandingPage serves the main HTML page.
func (pc *PageController) LandingPage(c *gin.Context) {
	if err := views.RenderTemplate(c.Writer, "index.html", nil); err != nil {
		c.String(500, err.Error())
		return
	}
}
