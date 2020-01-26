package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/tuckKome/konjac/handler"
)

func main() {

	router := gin.Default()
	router.Static("/templates", "./templates")
	router.LoadHTMLGlob("templates/*.html") // 事前にテンプレートをロード 相対パス
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", handler.Index)

	router.POST("/new", handler.New)

	router.GET("/translate/:text", handler.GetTranslation)

	router.GET("/login", handler.GetLogin)

	router.POST("/login", handler.PostLogin)

	router.GET("/logout", handler.Logout)

	router.GET("/signup", handler.GetSignup)

	router.POST("/signup", handler.PostSignup)

	router.GET("/my", handler.GetMy)

	router.GET("/error", handler.GetError)

	router.Run()
}
