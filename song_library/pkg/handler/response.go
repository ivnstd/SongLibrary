package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Возврат в формате сообщения об ошибке
func newErrorResponse(c *gin.Context, statusCode int, message string, description ...string) {
	// Формирование ответа
	response := gin.H{
		"status":  "error",
		"message": message,
	}

	// Возврат ответа
	logrus.Errorf("%v error.\tmessage: %v", statusCode, message)
	if len(description) > 0 {
		logrus.Errorf("Error description: %s", description[0])
	}
	c.AbortWithStatusJSON(statusCode, response)
}

// Возврат в формате сообщения об успешном выполнении запроса
func newSuccessResponse(c *gin.Context, statusCode int, dataKey string, dataValue interface{}) {
	// Формирование ответа
	response := gin.H{
		"status": "success",
		dataKey:  dataValue,
	}

	// Возврат ответа
	logrus.Infof("%v success.\t%v: %v", statusCode, dataKey, dataValue)
	c.JSON(statusCode, response)
}
