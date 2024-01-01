package handlers

import "github.com/gin-gonic/gin"

func (apiServerCodegenWrapper) PostAuthorize(ctx *gin.Context) {
	text := "This is the authorization handler and it is not implemented :)"
	ctx.Writer.Write([]byte(text))
}
