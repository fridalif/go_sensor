package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	c.HTML(
	  http.StatusOK,
	  "index.html",
	  gin.H{
		"title": "Home Page",
		"db":    db,
	  },
	)
}
