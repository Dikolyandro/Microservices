package server
import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/recktt77/Microservices-First-/api-gateway/config"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func New(cfg *config.Config) *Server {
	r := gin.Default()

	s := &Server{
		router: r,
		config: cfg,
	}

	s.routes()
	return s
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Starting API Gateway on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) routes() {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/products", s.handleProxyGetRaw(s.config.ProductServiceURL, "/api/v1/products"))
		v1.GET("/products/:id", s.handleProxyGet(s.config.ProductServiceURL, "/api/v1/products/"))
		v1.POST("/products", s.handleProxyPost(s.config.ProductServiceURL, "/api/v1/products"))
		v1.PATCH("/products/:id", s.handleProxyPatch(s.config.ProductServiceURL, "/api/v1/products/"))
		v1.DELETE("/products/:id", s.handleProxyDelete(s.config.ProductServiceURL, "/api/v1/products/"))

		v1.GET("/orders", s.handleProxyGetRaw(s.config.OrderServiceURL, "/api/v1/orders"))
		v1.GET("/orders/:id", s.handleProxyGet(s.config.OrderServiceURL, "/api/v1/orders/"))
		v1.POST("/orders", s.handleProxyPost(s.config.OrderServiceURL, "/api/v1/orders"))
		v1.PATCH("/orders/:id", s.handleProxyPatch(s.config.OrderServiceURL, "/api/v1/orders/"))
	}
}

func (s *Server) handleProxyGet(serviceURL string, pathPrefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		url := fmt.Sprintf("%s%s%s", serviceURL, pathPrefix, id)

		resp, err := http.Get(url)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

func (s *Server) handleProxyGetRaw(serviceURL string, path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := fmt.Sprintf("%s%s", serviceURL, path)
		resp, err := http.Get(url)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()
		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

func (s *Server) handleProxyPost(serviceURL, path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		resp, err := http.Post(fmt.Sprintf("%s%s", serviceURL, path), "application/json", bytes.NewReader(body))
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()
		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

func (s *Server) handleProxyPatch(serviceURL, pathPrefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		url := fmt.Sprintf("%s%s%s", serviceURL, pathPrefix, id)
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()
		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

func (s *Server) handleProxyDelete(serviceURL, pathPrefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		url := fmt.Sprintf("%s%s%s", serviceURL, pathPrefix, id)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create request"})
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()
		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}
