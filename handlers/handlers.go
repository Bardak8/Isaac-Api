package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"isaac-fansite/wiki"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.HTML(http.StatusOK, "search.html", gin.H{"query": "", "results": nil})
		return
	}

	results, err := wiki.Search(query)
	if err != nil {
		c.HTML(http.StatusOK, "search.html", gin.H{
			"query": query,
			"error": "Erreur lors de la recherche",
		})
		return
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"query":   query,
		"results": results,
	})
}

func (h *Handler) Detail(c *gin.Context) {
	title := c.Param("title")

	page, err := wiki.GetPage(title)
	if err != nil {
		c.HTML(http.StatusNotFound, "detail.html", gin.H{"error": "Page introuvable"})
		return
	}

	var isFav bool
	h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM favorites WHERE title=$1)", title).Scan(&isFav)

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"page":      page,
		"isFav":     isFav,
	})
}

func (h *Handler) Favorites(c *gin.Context) {
	rows, err := h.db.Query("SELECT title, thumbnail FROM favorites ORDER BY created_at DESC")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "favorites.html", gin.H{"error": "Erreur BDD"})
		return
	}
	defer rows.Close()

	type Fav struct {
		Title     string
		Thumbnail string
	}

	var favs []Fav
	for rows.Next() {
		var f Fav
		rows.Scan(&f.Title, &f.Thumbnail)
		favs = append(favs, f)
	}

	c.HTML(http.StatusOK, "favorites.html", gin.H{"favorites": favs})
}

func (h *Handler) AddFavorite(c *gin.Context) {
	title := c.PostForm("title")
	thumbnail := c.PostForm("thumbnail")

	if title == "" {
		c.Redirect(http.StatusFound, "/favorites")
		return
	}

	h.db.Exec("INSERT INTO favorites (title, thumbnail) VALUES ($1, $2) ON CONFLICT (title) DO NOTHING", title, thumbnail)
	c.Redirect(http.StatusFound, "/detail/"+title)
}

func (h *Handler) RemoveFavorite(c *gin.Context) {
	title := c.PostForm("title")
	h.db.Exec("DELETE FROM favorites WHERE title=$1", title)
	c.Redirect(http.StatusFound, "/detail/"+title)
}