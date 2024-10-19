package main

import (
	"go_proj_example/internal/database"
	"go_proj_example/internal/handlers"
	"go_proj_example/internal/routes"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	router := gin.Default()
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get dir:", err)
	}
	templatesPath := filepath.Join(basePath, "www", "templates")
	assetsPath := filepath.Join(basePath, "www", "assets")

	router.HTMLRender = handlers.LoadTemplates(templatesPath)
	router.Static("/static", assetsPath)
	routes.SetupRouter(router)

	router.Run(":3000")
}
