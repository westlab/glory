package main

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/westlab/glory/web"
	"log"
)

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("progress", "templates/base.html", "templates/progress.html")
	r.AddFromFiles("notFound", "templates/base.html", "templates/404.html")
	return r
}

func main() {
	engine := gin.Default()
	//engine.LoadHTMLGlob("templates/*.html")
	engine.HTMLRender = createRender()
	engine.Static("/public", "./public")
	engine.StaticFile("/favicon.ico", "./public/favicon.ico")

	engine.GET("/", web.IndexHandler)
	engine.GET("/data/:id", web.ProgressHandler)
	engine.NoRoute(web.NotFoundHandler)

	err := engine.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
