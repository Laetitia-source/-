package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Connexion à la base de données PostgreSQL
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=laetitia dbname=anime_api sslmode=disable")
	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données : ", err)
	}
	defer db.Close()

	// Test de connexion
	if err := db.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à la base de données : ", err)
	}

	log.Println("Connexion à la base de données réussie!")

	// Initialisation du serveur Gin
	router := gin.Default()

	// Routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Bienvenue dans l'API des recommandations d'animés!"})
	})

	router.GET("/animes", getAnimes)
	router.POST("/animes", addAnime)
	router.GET("/animes/search", searchAnimesByGenre)
	router.PUT("/animes/:id", updateAnime)
	router.DELETE("/animes/:id", deleteAnime)

	// Lancer le serveur
	router.Run(":8080")
}

// getAnimes récupère tous les animés
func getAnimes(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, genre, release_year, rating, description FROM animes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des animés"})
		return
	}
	defer rows.Close()

	var animes []map[string]interface{}
	for rows.Next() {
		var id int
		var title, genre, description string
		var releaseYear int
		var rating float64

		rows.Scan(&id, &title, &genre, &releaseYear, &rating, &description)
		animes = append(animes, gin.H{
			"id":          id,
			"title":       title,
			"genre":       genre,
			"releaseYear": releaseYear,
			"rating":      rating,
			"description": description,
		})
	}

	c.JSON(http.StatusOK, animes)
}

// addAnime ajoute un nouvel animé
func addAnime(c *gin.Context) {
	var anime struct {
		Title       string  `json:"title" binding:"required"`
		Genre       string  `json:"genre"`
		ReleaseYear int     `json:"releaseYear"`
		Rating      float64 `json:"rating"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	_, err := db.Exec(
		"INSERT INTO animes (title, genre, release_year, rating, description) VALUES ($1, $2, $3, $4, $5)",
		anime.Title, anime.Genre, anime.ReleaseYear, anime.Rating, anime.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout de l'animé"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Animé ajouté avec succès"})
}

// searchAnimesByGenre permet de filtrer les animés par genre
func searchAnimesByGenre(c *gin.Context) {
	genre := c.Query("genre")

	if genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le paramètre 'genre' est requis"})
		return
	}

	query := "SELECT id, title, genre, release_year, rating, description FROM animes WHERE genre = $1"
	rows, err := db.Query(query, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recherche"})
		return
	}
	defer rows.Close()

	var animes []map[string]interface{}
	for rows.Next() {
		var id int
		var title, genre, description string
		var releaseYear int
		var rating float64

		rows.Scan(&id, &title, &genre, &releaseYear, &rating, &description)
		animes = append(animes, gin.H{
			"id":          id,
			"title":       title,
			"genre":       genre,
			"releaseYear": releaseYear,
			"rating":      rating,
			"description": description,
		})
	}

	c.JSON(http.StatusOK, animes)
}

// updateAnime met à jour un animé
func updateAnime(c *gin.Context) {
	id := c.Param("id")

	var anime struct {
		Title       string  `json:"title"`
		Genre       string  `json:"genre"`
		ReleaseYear int     `json:"releaseYear"`
		Rating      float64 `json:"rating"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	_, err := db.Exec(
		"UPDATE animes SET title = $1, genre = $2, release_year = $3, rating = $4, description = $5 WHERE id = $6",
		anime.Title, anime.Genre, anime.ReleaseYear, anime.Rating, anime.Description, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animé mis à jour avec succès"})
}

// deleteAnime supprime un animé
func deleteAnime(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM animes WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animé supprimé avec succès"})
}
