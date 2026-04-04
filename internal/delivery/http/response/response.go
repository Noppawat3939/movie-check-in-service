package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, statusCode int, message string, data ...any) {
	body := gin.H{"success": false, "message": message}
	if len(data) > 0 && data[0] != nil {
		body["data"] = data[0]
	}

	c.JSON(statusCode, body)
}
