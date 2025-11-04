package handlers

type Handlers struct {
	HealthHandler *HealthHandler
}

/*
Create a new handlers instance
Returns a pointer to the handlers instance
*/
func NewHandlers() *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
	}
}
