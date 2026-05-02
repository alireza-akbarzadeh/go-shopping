package routes

import "github.com/alireza-akbarzadeh/shopping-platform/middleware"

// RegisterMiddlewares attaches any custom middleware not already applied globally
func (r *Router) RegisterMiddlewares() {
	r.engine.Use(middleware.CORS())
}
