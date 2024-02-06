package middleware

import (
	"context"
	"log"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Adds a db connection from the config as "DB", with a 1 second timeout and a cancel func to call in the request as "DB_CANCEl";
func AddMiddleware(r *gin.Engine, conf *config.AppConfig) {
	addDbConnPooling(r, conf.DB)
}

func addDbConnPooling(r *gin.Engine, dbconf config.DbConfig) {
	r.Use(func(ctx *gin.Context) {
		var timeoutContext context.Context
		var cancel context.CancelFunc
		var dbconn *gorm.DB

		timeoutContext, cancel = context.WithTimeout(context.Background(), time.Second*2)

		if conn, err := dbconf.Connect(); err != nil || conn == nil {
			log.Fatal("Error occured while adding db pooling middleware\n", err)
		} else {
			dbconn = conn.WithContext(timeoutContext)
		}

		ctx.Set("DB", dbconn)
		ctx.Set("DB_CANCEL", cancel)
	})
}
