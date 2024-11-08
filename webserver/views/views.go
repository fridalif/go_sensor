package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(
	  http.StatusOK,
	  "index.html",
	  gin.H{
		"title": "Home Page",
	  },
	)
  }