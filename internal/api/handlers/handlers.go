package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/auth"
	"github.com/Marattttt/portfolio/portfolio_back/internal/services"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"gorm.io/gorm"
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
	dbconn, cancelDB := getDbConnCancel(ctx)
	defer cancelDB()

	authRequest := api.GetData[api.AuthRequest](ctx.Request.Body)

	if authRequest == nil {
		ctx.Writer.WriteString("Could not unmarshal json")
		ctx.Status(http.StatusBadRequest)
		return
	}

	service := services.NewGuestsService(dbconn)
	guest := service.GetGuest(authRequest.Id)

	if guest == nil {
		ctx.Writer.WriteString("Guest not found")
		ctx.Status(http.StatusNotFound)
		return
	}

	err := auth.Validate([]byte(guest.Secret), []byte(guest.Salt), authRequest.Password)
	if err != nil {
		applog.Error(applog.Auth, err)

		ctx.Status(http.StatusForbidden)
		return
	}

	ctx.Writer.WriteString("Success!!")
}

func (apiServerCodegenWrapper) GetGuestsGuestId(ctx *gin.Context, guestId int) {
	dbconn, cancelDB := getDbConnCancel(ctx)
	defer cancelDB()

	service := services.NewGuestsService(dbconn)

	guest := service.GetGuest(guestId)
	if guest == nil {
		ctx.Writer.WriteString("Guest not found")
		ctx.Status(http.StatusNotFound)
		return
	}
	json.NewEncoder(ctx.Writer).Encode(api.ToGuestResponse(*guest))
}

func (apiServerCodegenWrapper) PostGuests(ctx *gin.Context) {
	dbconn, cancelDB := getDbConnCancel(ctx)
	defer cancelDB()

	guestRequest := api.GetData[api.GuestRequest](ctx.Request.Body)

	if guestRequest == nil {
		ctx.Writer.WriteString("Could not unmarshal json")
		ctx.Status(http.StatusBadRequest)
		return
	}

	service := services.NewGuestsService(dbconn)

	guest := api.ToGuest(*guestRequest)

	if g, err := service.NewGuest(guest); err != nil {
		applog.Error(applog.DB, err)

		ctx.Writer.WriteString("Could not save guest")
		ctx.Status(http.StatusBadGateway)
		return
	} else {
		guest = *g
	}

	json.NewEncoder(ctx.Writer).Encode(api.ToGuestResponse(guest))
	ctx.Status(http.StatusCreated)
}

func (apiServerCodegenWrapper) PatchGuestsGuestId(ctx *gin.Context, guestId int) {
	text := "This is the patch guest by id handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) GetGuestsGuestIdStats(ctx *gin.Context, guestId int) {
	text := "This is the get guest stats handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) GetStats(ctx *gin.Context) {
	text := "This is the get stats handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func getDbConnCancel(ctx *gin.Context) (*gorm.DB, context.CancelFunc) {
	db := ctx.MustGet("DB").(*gorm.DB)
	cancel := ctx.MustGet("DB_CANCEL").(context.CancelFunc)

	return db, cancel
}
