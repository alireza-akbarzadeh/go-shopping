// Package routes contains all route definitions and middleware registrations for the shopping platform API
package routes

import "github.com/alireza-akbarzadeh/shopping-platform/middleware"

// RegisterMiddlewares attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddlewares() {
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.RateLimitMiddleware(100, 200))
}
