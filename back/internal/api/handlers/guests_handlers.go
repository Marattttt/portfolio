package handlers

import (
	"github.com/gin-gonic/gin"
)

func (apiServerCodegenWrapper) PostGuests(ctx *gin.Context) {
	text := "This is the post guests handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}

func (apiServerCodegenWrapper) GetGuestsGuestId(ctx *gin.Context, guestId int) {
	text := "This is the get guest by id handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}
