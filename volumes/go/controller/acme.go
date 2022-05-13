package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AcmeHandle(c *gin.Context) {
	// acmename := c.Param("acmename")
	c.String(http.StatusOK, "QHI33DW6CTddO7olw8wKOKvQuV4vcqYH3QCO_GNAmMM.NW7vp2syLk8vHn3Tt6GpQs9-M-wlEZg-GP-NlS65HTE")
}
