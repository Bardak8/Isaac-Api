package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"isaac-fansite/db"
	"isaac-fansite/handlers"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	h := handlers.New(database)

	r.GET("/", h.Index)
	r.GET("/search", h.Search)
	r.GET("/detail/:title", h.Detail)
	r.GET("/favorites", h.Favorites)
	r.POST("/favorites/add", h.AddFavorite)
	r.POST("/favorites/remove", h.RemoveFavorite)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting on :%s", port)
	r.Run(":" + port)
}