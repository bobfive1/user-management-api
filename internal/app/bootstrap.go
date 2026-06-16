package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bobfive1/user-management-api/internal/config"
	"github.com/bobfive1/user-management-api/internal/db"
	"github.com/bobfive1/user-management-api/internal/logger"
	"github.com/bobfive1/user-management-api/internal/middleware"
	"github.com/bobfive1/user-management-api/internal/userprofile"
	"github.com/bobfive1/user-management-api/internal/validation"

	errorInt "github.com/bobfive1/user-management-api/internal/error"

	"github.com/gin-gonic/gin"
)

var (
	Logger    = logger.GetDefaultLogger()
	startTime = time.Now()
)

type serverAPI struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

type App interface {
	Start()
	Stop()
	Serv() *http.Server
}

func Bootstrap() (App, *sync.WaitGroup) {
	configApp, err := config.LoadConfig()
	if err != nil {
		Logger.Fatalf("LoadConfig error: %v", err)
	}

	if err := logger.ApplyConfig(configApp.App.Name, configApp.Logging.Level); err != nil {
		Logger.Fatalf("logger ApplyConfig error: %v", err)
	}

	validation.Init()

	database, err := db.SetupClientPostgres(*configApp)
	if err != nil {
		Logger.Fatalf("connect database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.Ping(ctx); err != nil {
		Logger.Fatalf("ping database: %v", err)
	}

	repo := userprofile.NewUserProfileRepository(database)
	service := userprofile.NewUserProfileService(repo)

	apiServer := ApiServer(configApp, service)

	var wg sync.WaitGroup
	addHookShutdown(&wg, func() {
		apiServer.Stop()
		db.StopPostgres(database)
	})
	return apiServer, &wg
}

func ApiServer(config *config.AppConfig, service userprofile.UserProfileService) App {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.TraceIDMiddleware())
	router.Use(errorInt.MiddlewareErrorHandler())

	handler := userprofile.NewUserProfileHandler(service)
	handler.RegisterRoutes(router.Group("/api/v1"))

	router.GET("/health", func(ctx *gin.Context) {
		duration := time.Since(startTime)
		uptimeStr := formatUptime(duration)

		ctx.JSON(http.StatusOK, gin.H{
			"status":    "UP",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"uptime":    uptimeStr, // แสดงผลที่นี่
		})
	})

	server := &http.Server{
		Addr:              config.ServerAPI.Port,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       config.ServerAPI.ReadTimeout,
		WriteTimeout:      config.ServerAPI.WriteTimeout,
		IdleTimeout:       config.ServerAPI.IdleTimeout,
	}

	shutdownTimeout := config.ServerAPI.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 15 * time.Second
	}

	return &serverAPI{server: server, shutdownTimeout: shutdownTimeout}
}

func (s *serverAPI) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		Logger.Fatalf("Start server http error :%v", err)
	}
}

func (s *serverAPI) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		Logger.Fatalf("Shutdown api server error cause by: %s", err)
	}
	Logger.Info("Shutdown api server")
}

func (s *serverAPI) Serv() *http.Server {
	return s.server
}

func addHookShutdown(wg *sync.WaitGroup, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = <-c
		f()
	}()
}

func formatUptime(d time.Duration) string {
	// คำนวณจำนวน วัน, ชั่วโมง, นาที, วินาที
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
