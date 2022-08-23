package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"

	"github.com/pandaychen/goes-wrapper/pymicrosvc/loger/xzap"
)

const (
	CtxRequestID = "X-REQ-ID"
)

//RequestId 中间件
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, if exists then use it
		requestId := c.Request.Header.Get(CtxRequestID)
		//requestId := c.GetHeader("traceId")
		// Create request id with UUID4
		if requestId == "" {
			uid := uuid.New()
			requestId = uid.String()
			fmt.Println("xx")
		}
		c.Set(CtxRequestID, requestId)
		// Set X-Request-Id header with HTTP response
		c.Writer.Header().Set(CtxRequestID, requestId)
		c.Next()
	}
}

func TraceIdLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(CtxRequestID)
		traceId, ok := c.MustGet(CtxRequestID).(string)
		if !ok {
			panic("fatal errors")
		}
		ctx, log := xzap.GetCtxZapLoger().WithContextAndFields(c.Request.Context(), zap.Any("traceId", traceId))
		c.Request = c.Request.WithContext(ctx)
		log.Info("Set TraceIdLog done")
		c.Next()
	}
}

func main() {
	//设置全局日志属性
	xzap.NewCtxZapLoger(xzap.SetDevelopment(false), xzap.SetWriteFile(true))
	g := gin.New()
	g.Use(RequestId(), TraceIdLog())
	g.GET("/test", func(context *gin.Context) {
		log := xzap.GetCtxZapLoger().GetCtx(context.Request.Context())
		log.Info("test")
		log.Debug("test")
		context.JSON(200, "success")
	})
	xzap.GetLoger().Info("hconf example success")
	http.ListenAndServe(":8888", g)
}
