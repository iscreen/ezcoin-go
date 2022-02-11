package core

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/initialize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunServer() {
	if global.GVA_CONFIG.System.UseMultipoint {
		// 初始化redis服务
		initialize.Redis()
	}

	if global.GVA_DB != nil {
		service.LoadAll()
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		fmt.Printf("start server\n")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.GVA_LOG.Error(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	global.GVA_LOG.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		global.GVA_LOG.Fatal(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	global.GVA_LOG.Info("Server exiting")
}

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
