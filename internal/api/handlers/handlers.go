package handlers

import "knowstack/internal/core/services"

type Handlers struct {
	HealthHandler *HealthHandler
	UserHandler   *UserHandler
}

/*
Create a new handlers instance
Returns a pointer to the handlers instance
*/
func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		UserHandler:   NewUserHandler(service.UserService),
	}
}
