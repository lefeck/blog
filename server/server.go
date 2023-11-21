package server

import (
	"blog/common"
	"blog/config"
	"blog/controller"
	"blog/database"
	"blog/repository"
	"blog/service"
	"blog/utils/set"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine      *gin.Engine
	config      *config.Config
	logger      *logrus.Logger
	controllers []controller.Controller
	repository  repository.Repository
}

func New(config *config.Config, logger *logrus.Logger) (*Server, error) {
	// intiable mysql
	db, err := database.NewMysql(&config.DB)
	if err != nil {
		return nil, err
	}
	//new initable, insert sql table
	repository := repository.NewRepository(db)
	if config.DB.Migrate {
		if err = repository.Migrate(); err != nil {
			return nil, err
		}
	}

	// user
	userService := service.NewUserService(repository.User())
	userController := controller.NewUserController(userService)

	controllers := []controller.Controller{userController}

	gin.SetMode(config.Server.ENV)

	e := gin.New()
	e.Use(
		gin.Recovery(),
	)
	e.LoadHTMLFiles("static/index.html")

	return &Server{
		engine:      e,
		config:      config,
		logger:      logger,
		controllers: controllers,
		repository:  repository,
	}, nil
}

func (s *Server) Run() error {
	defer s.Close()
	s.Routers()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server, %v", err)
		}
	}()

	// 平滑关闭进程
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.GracefulShutdownPeriod)*time.Second)
	defer cancel()

	ch := <-sig
	log.Fatalf("Receive signal: %s", ch)

	return server.Shutdown(ctx)
}

func (s *Server) Close() {
	if err := s.repository.Close(); err != nil {
		s.logger.Warnf("failed to close repository, %v", err)
	}
}

func (s *Server) Routers() {
	root := s.engine

	// register non-resource routers
	root.GET("/", common.WrapFunc(s.getRoutes))
	root.GET("/index", controller.Index)
	root.GET("/healthz", common.WrapFunc(s.Ping))
	//root.GET("/version", common.WrapFunc(version.Get))
	//root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	// swagger doc
	if gin.Mode() != gin.ReleaseMode {
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	api := root.Group("/api/v1")
	controllers := make([]string, 0, len(s.controllers))
	for _, router := range s.controllers {
		router.RegisterRoute(api)
		controllers = append(controllers, router.Name())
	}
	logrus.Infof("server enabled controllers: %v", controllers)
}

// @Summary Healthz
// @Produce json
// @Tags healthz
// @Success 200 {string}  string    "ok"
// @Router /healthz [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (s *Server) getRoutes() []string {
	paths := set.NewString()
	for _, r := range s.engine.Routes() {
		if r.Path != "" {
			paths.Insert(r.Path)
		}
	}
	return paths.Slice()
}

type ServerStatus struct {
	Ping         bool `json:"ping"`
	DBRepository bool `json:"dbRepository"`
}

func (s *Server) Ping() *ServerStatus {
	status := &ServerStatus{Ping: true}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.repository.Ping(ctx); err != nil {
		status.DBRepository = true
	}
	return status
}
