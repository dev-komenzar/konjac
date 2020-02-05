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
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	router.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")
	router.StaticFile("/favicon-16x16.png", "./resources/favicon-16x16.png")
	router.StaticFile("/manifest.json", "./resources/manifest.json")
	router.StaticFile("/safari-pinned-tab.svg", "./resources/safari-pinned-tab.svg")
	router.StaticFile("/android-chrome-192x192.png", "./resources/android-chrome-192x192.png")
	router.StaticFile("/android-chrome-256x256.png", "./resources/android-chrome-256x256.png")

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
