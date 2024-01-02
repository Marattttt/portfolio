package handlers

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/api"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

type apiServerCodegenWrapper struct{}

func SetupHandlers(r *gin.Engine) error {
	api.RegisterHandlers(r, apiServerCodegenWrapper{})

	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}

	r.Use(middleware.OapiRequestValidator(swagger))
	return nil
}
