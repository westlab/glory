package main

import (
	"github.com/westlab/glory/config"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"github.com/westlab/glory/db"
)

var (
	GloryConfig *config.Conf
	dsn         string
	port        string
)

func init() {

	var err error
	if GloryConfig, err = config.LoadConfig("/app/config.json"); err != nil {
		log.Fatal(err)
	}

	dsn = os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("environment variables: DSN is not defined")
	}

	port = os.Getenv("GLORY_PORT")
	if port == "" {
		port = ":8080"
	}

}

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "/app/templates/base.html", "/app/templates/index.html")
	r.AddFromFiles("progress", "/app/templates/base.html", "/app/templates/progress.html")
	r.AddFromFiles("notFound", "/app/templates/base.html", "/app/templates/404.html")
	return r
}

func main() {

	done, err := db.InitializeDB(dsn)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer done()

	engine := gin.Default()
	//engine.LoadHTMLGlob("templates/*.html")
	if gin.Mode() == "release" {
		logFile := "/app/log/glory_log_" + time.Now().Format(time.RFC3339) + ".log"
		f, err := os.Create(logFile)
		if err != nil {
			log.Fatalf("log file '%s' is invalid: %v", logFile, err)
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

	err = engine.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
