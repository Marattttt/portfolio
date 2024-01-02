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

func (apiServerCodegenWrapper) PostAuthorize(ctx *gin.Context) {
	text := "This is the authorization handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) PostGuests(ctx *gin.Context) {
	text := "This is the post guests handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) GetGuestsGuestId(ctx *gin.Context, guestId int) {
	text := "This is the get guest by id handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) GetStats(ctx *gin.Context) {
	text := "This is the get stats handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}
