package error

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bobfive1/user-management-api/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	log = logger.GetDefaultLogger()
)

func MiddlewareErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		c.Next()

		latency := time.Since(startAt)

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			httpCode, reponse := getErrorResponse(err)

			traceID, _ := c.Get("trace_id")

			log.With(
				zap.Any("trace_id", traceID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", httpCode),
				zap.String("client_ip", c.ClientIP()),
				zap.Duration("latency", latency),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("body", structToJsonString(reponse))).Error("response")

			c.AbortWithStatusJSON(httpCode, reponse)
		}

	}
}

func getErrorResponse(err error) (httpCode int, reponse ErrorResponse) {
	switch e := err.(type) {
	case FieldValidationError:
		return 400, ErrorResponse{
			TimeStamp:    time.Now(),
			ErrorCode:    e.ErrorCode,
			ErrorMessage: e.ErrorMessage,
			ErrorDetail:  e.ErrorDetail,
		}
	case NotFoundError:
		return 404, ErrorResponse{
			TimeStamp:    time.Now(),
			ErrorCode:    e.ErrorCode,
			ErrorMessage: e.ErrorMessage,
			ErrorDetail:  e.ErrorDetail,
		}
	case ApiServerError:
		return 500, ErrorResponse{
			TimeStamp:    time.Now(),
			ErrorCode:    e.ErrorCode,
			ErrorMessage: e.ErrorMessage,
			ErrorDetail:  e.ErrorDetail,
		}
	case InternalServerError:
		return 500, ErrorResponse{
			TimeStamp:    time.Now(),
			ErrorCode:    e.ErrorCode,
			ErrorMessage: e.ErrorMessage,
			ErrorDetail:  e.ErrorDetail,
		}
	default:
		return 500, ErrorResponse{
			TimeStamp:    time.Now(),
			ErrorCode:    "500",
			ErrorMessage: "Internal server error",
			ErrorDetail:  err.Error(),
		}
	}

}

func structToJsonString(body any) string {
	json, err := json.Marshal(body)
	if err != nil {
		return fmt.Sprintf("%+v", body)
	}
	return string(json)
}
