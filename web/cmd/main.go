package main

import (
	"log"
  "os"
  "time"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/westlab/glory/web"
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
  if gin.Mode() == "release" {
      logFile := "log/glory_log_"+time.Now().Format(time.RFC3339)+".log"
      f, err := os.Create(logFile)
      if err != nil {
          log.Fatalf("log file %s is invalid", logFile)
      }
      gin.DefaultErrorWriter = f
      gin.DefaultWriter = f
      engine.Use(gin.LoggerWithWriter(f))
  }
	engine.HTMLRender = createRender()
	engine.Static("/public", "./public")
	engine.StaticFile("/favicon.ico", "./public/favicon.ico")

	engine.GET("/", web.IndexHandler)
	engine.GET("/data/:id", web.ProgressHandler)
	engine.NoRoute(web.NotFoundHandler)

	err := engine.Run(":80")
	if err != nil {
		log.Fatal(err)
	}
}
