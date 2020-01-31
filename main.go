package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/tuckKome/konjac/db"
	"github.com/tuckKome/konjac/handler"
	"github.com/tuckKome/konjac/render"
)

func main() {
	db.Init()
	router := gin.Default()
	router.Static("/templates", "./templates")
	router.Static("/bootstrap-4", "./bootstrap-4")
	router.Static("/jsmind", "./jsmind")

	router.HTMLRender = render.LoadTemplates("./templates") // 事前にテンプレートをロード multitemplateで
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
