package main
import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)
var db *pgxpool.Pool

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/users", getUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9003"
	}

	log.Println("🚀 Starting server on port", port)
	router.Run(":" + port)
}

func getUsers(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		users = append(users, gin.H{"id": id, "name": name})
	}
	c.JSON(http.StatusOK, users)
}
