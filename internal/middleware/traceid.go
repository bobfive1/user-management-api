package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	xTraceID   = "X-Trace-ID"
	xRequestID = "X-Request-ID"
)

func TraceIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		traceId := ctx.GetHeader(xTraceID)
		if traceId == "" {
			traceId = ctx.GetHeader(xRequestID)
		}

		if traceId == "" {
			traceId = uuid.New().String()
		}

		ctx.Set("trace_id", traceId)
		ctx.Writer.Header().Set("X-Trace-ID", traceId)

		ctx.Next()
	}
}
