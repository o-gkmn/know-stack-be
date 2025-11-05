package handlers

import "knowstack/internal/core/service"

type Handlers struct {
	HealthHandler *HealthHandler
	UserHandler   *UserHandler
}

/*
Create a new handlers instance
Returns a pointer to the handlers instance
*/
func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		UserHandler:   NewUserHandler(service.UserService),
	}
}
