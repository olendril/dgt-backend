// Package guild_api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package guild_api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Error *string `json:"error,omitempty"`
}

// GuildInfo defines model for GuildInfo.
type GuildInfo struct {
	Name   string `json:"name"`
	Server string `json:"server"`
}

// GuildResponse defines model for GuildResponse.
type GuildResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Server string `json:"server"`
}

// PostGuildsJSONRequestBody defines body for PostGuilds for application/json ContentType.
type PostGuildsJSONRequestBody = GuildInfo

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /guilds)
	GetGuilds(c *gin.Context)

	// (POST /guilds)
	PostGuilds(c *gin.Context)

	// (GET /guilds/{id})
	GetGuildsId(c *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetGuilds operation middleware
func (siw *ServerInterfaceWrapper) GetGuilds(c *gin.Context) {

	c.Set(BasicAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetGuilds(c)
}

// PostGuilds operation middleware
func (siw *ServerInterfaceWrapper) PostGuilds(c *gin.Context) {

	c.Set(BasicAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostGuilds(c)
}

// GetGuildsId operation middleware
func (siw *ServerInterfaceWrapper) GetGuildsId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: false})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BasicAuthScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetGuildsId(c, id)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/guilds", wrapper.GetGuilds)
	router.POST(options.BaseURL+"/guilds", wrapper.PostGuilds)
	router.GET(options.BaseURL+"/guilds/:id", wrapper.GetGuildsId)
}