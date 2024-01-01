package handlers

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/api"
	"github.com/gin-gonic/gin"
)

type apiServerCodegenWrapper struct{}

func SetupHandlers(r *gin.Engine) {
	api.RegisterHandlers(r, apiServerCodegenWrapper{})
}
