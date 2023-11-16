package controllers

import (
	"context"
	"fmt"
	"github.com/Go-routine-4995/ossfrontend/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	routerfilenametag = "file"
	maxMessageSize    = 8000000 // max message 8M
)

type IService interface {
	CreateRouters(ctx context.Context, r []domain.Router, tenant string) ([]byte, error)
	DeleteRouters(ctx context.Context, r []domain.Router, tenant string) error
	GetRoutersPage(ctx context.Context, paginationByte []byte, tenant string) (domain.Response, error)
	GetRouters(ctx context.Context, r domain.Router, tenant string) (domain.Router, error)
}

type IAuthentication interface {
	AuthMiddleware() gin.HandlerFunc
}

type ApiServer struct {
	port string
	ctx  context.Context
	next IService
	auth IAuthentication
}

func NewApiServer(next interface{}, port string, a IAuthentication) *ApiServer {
	return &ApiServer{
		port: port,
		next: next.(IService),
		ctx:  context.Background(),
		auth: a,
	}
}

func (a *ApiServer) Start() {
	router := gin.Default()

	router.GET("/", a.IndexHandler)

	oss := router.Group("/api/v4/oss")
	oss.Use(a.auth.AuthMiddleware())
	oss.POST("/routers", a.CreateRouters)
	oss.DELETE("/routers", a.DeleteRouters)
	oss.GET("/routers", a.GetRouters)

	err := http.ListenAndServe(":"+a.port, router)
	if err != nil {
		fmt.Println(err)
	}
}

func (a *ApiServer) IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world this is OSS API",
	})
}
