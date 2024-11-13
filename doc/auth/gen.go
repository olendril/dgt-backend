// Package auth_api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package auth_api

import (
	"github.com/gin-gonic/gin"
)

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Link string `json:"link"`
}

// RedirectResponse defines model for RedirectResponse.
type RedirectResponse struct {
	AccessToken    string `json:"access_token"`
	AvatarLink     string `json:"avatar_link"`
	ExpirationDate string `json:"expiration_date"`
	Username       string `json:"username"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /login)
	GetLogin(c *gin.Context)

	// (GET /redirect)
	GetRedirect(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetLogin operation middleware
func (siw *ServerInterfaceWrapper) GetLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetLogin(c)
}

// GetRedirect operation middleware
func (siw *ServerInterfaceWrapper) GetRedirect(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetRedirect(c)
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

	router.GET(options.BaseURL+"/login", wrapper.GetLogin)
	router.GET(options.BaseURL+"/redirect", wrapper.GetRedirect)
}
