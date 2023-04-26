package api

import (
	"github.com/mohammadrabetian/ports/docs"
	"github.com/mohammadrabetian/ports/handlers"
	"github.com/mohammadrabetian/ports/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Store  *Store
	config util.Config
	router *gin.Engine
}

// creates an HTTP server
func NewServer(config util.Config) *Server {
	store := NewStore(config)
	server := &Server{config: config, Store: store}
	server.setupRouter()
	return server

}

func (s *Server) setupRouter() {
	router := gin.Default()

	// Register the request logger middleware
	router.Use(s.SetupRequestLogger())

	portGroup := router.Group("/api/v1/ports")

	err := router.SetTrustedProxies([]string{"192.168.1.2"})
	if err != nil {
		logrus.Fatalf("failed to set trusted proxies")
	}

	// set swagger info
	docs.SwaggerInfo.Title = "Ports Swagger API"
	docs.SwaggerInfo.Description = "Interact with the APIs here"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = s.config.HTTPServer.Address
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// business logic controllers
	{
		portGroup.GET("/:id", handlers.GetPortByID)
		portGroup.GET("/", handlers.ListPorts)
	}

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) SetupRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithFields(logrus.Fields{
			"request_method": c.Request.Method,
			"request_path":   c.Request.URL.Path,
		})

		logger.Info("Started handling request")
		c.Next()
		logger.Info("Completed handling request")
	}
}
