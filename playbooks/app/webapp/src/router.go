package main

import (
	"web-sqlx-sqlite-user/src/controllers"
	tmpl "web-sqlx-sqlite-user/src/lib/template"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	// テンプレートロード
	t, err := tmpl.LoadTemplates("src/templates")
	if err != nil {
		panic("テンプレートのロードに失敗しました: " + err.Error())
	}
	router.SetHTMLTemplate(t)

	// 静的ファイル
	router.Static("/assets", "public/assets")

	// ルータ登録
	registerUserRouter(router)

	// 404 ハンドラ
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{"page_title": "Not Found"})
	})

	return router
}

func registerUserRouter(router *gin.Engine) {
	uc := &controllers.UserController{}

	// リダイレクト
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/users")
	})

	// ユーザ管理
	router.GET("/users", uc.Index)
	router.GET("/users/new", uc.New)
	router.POST("/users", uc.Create)
	router.GET("/users/:id/edit", uc.Edit)
	router.POST("/users/:id", uc.Update)
	router.POST("/users/:id/delete", uc.Delete)
}
