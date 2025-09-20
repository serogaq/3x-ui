// Package global provides global variables and interfaces for accessing web and subscription servers.
package global

import (
	"context"
	_ "unsafe"

	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
)

var (
	webServer WebServer
	subServer SubServer
	caching   Cache
)

// WebServer interface defines methods for accessing the web server instance.
type WebServer interface {
	GetCron() *cron.Cron     // Get the cron scheduler
	GetCtx() context.Context // Get the server context
}

// SubServer interface defines methods for accessing the subscription server instance.
type SubServer interface {
	GetCtx() context.Context // Get the server context
}

type Cache interface {
	Memory() *cache.Cache
	GetCtx() context.Context
}

// SetWebServer sets the global web server instance.
func SetWebServer(s WebServer) {
	webServer = s
}

// GetWebServer returns the global web server instance.
func GetWebServer() WebServer {
	return webServer
}

// SetSubServer sets the global subscription server instance.
func SetSubServer(s SubServer) {
	subServer = s
}

// GetSubServer returns the global subscription server instance.
func GetSubServer() SubServer {
	return subServer
}

func SetCache(c Cache) {
	caching = c
}

func GetCache() Cache {
	return caching
}
