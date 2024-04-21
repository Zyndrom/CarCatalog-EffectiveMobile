package main

import (
	"CarsCatalog/internal/repository"
	"CarsCatalog/internal/router"
	"CarsCatalog/internal/service/car"
	"fmt"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/swaggo/gin-swagger"

	_ "github.com/swaggo/files"
)

// swagger embed files
//	@title			Cars Catalog
//	@version		1.0
//	@description	Cars Catalog.

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	setLogger()
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file not found.")
	}
	repository := repository.New()
	userService := car.New(repository)
	router := router.New(userService)
	router.StartServer()
}

func setLogger() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}
