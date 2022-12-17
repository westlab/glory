package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "/app/templates/base.html", "/app/templates/index.html")
	r.AddFromFiles("progress", "/app/templates/base.html", "/app/templates/progress.html")
	r.AddFromFiles("notFound", "/app/templates/base.html", "/app/templates/404.html")
	return r
}

func main() {

	done, err := InitializeDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer done()

	Port := os.Getenv("GLORY_PORT")
	if Port == "" {
		Port = ":8080"
	}

	engine := gin.Default()
	//engine.LoadHTMLGlob("templates/*.html")
	if gin.Mode() == "release" {
		logFile := "log/glory_log_" + time.Now().Format(time.RFC3339) + ".log"
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

	engine.GET("/", IndexHandler)
	engine.GET("/data/:id", ProgressHandler)
	engine.NoRoute(NotFoundHandler)

	err = engine.Run(Port)
	if err != nil {
		log.Fatal(err)
	}
}
