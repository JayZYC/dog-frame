package main

import (
	"context"
	"github.com/dog-frame/api/router"
	"github.com/dog-frame/common/enum"
	"github.com/dog-frame/common/logger"
	"github.com/dog-frame/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 配置初始化
	// config.init()

	// 日志初始化
	// logger.init()

	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// 路由初始化
	g := gin.New()
	router.Register(g)
	server := http.Server{
		Addr:    ":8080",
		Handler: g,
	}

	ctx := context.Background()

	// 创建系统信号接收器
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		if err := server.Shutdown(ctx); err != nil {
			logger.Error(ctx, "ShutdownServerError", "err", err)
		}
	}()

	logger.Info(ctx, "Starting DOG FRAME HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// 服务正常收到关闭信号后Close
			logger.Info(ctx, "Server closed under request")
		} else {
			// 服务异常关闭
			logger.Error(ctx, "Server closed unexpected", "err", err)
		}
	}
}
