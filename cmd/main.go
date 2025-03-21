package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmscatena/Fatec_Sert_SGCourse/infra"
	"github.com/jmscatena/Fatec_Sert_SGCourse/routes"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading Configuration File")
	}
	gin.SetMode(os.Getenv("SET_MODE"))
	dbConn := infra.Connection{Db: nil, NoSql: nil}
	_, err = dbConn.InitDB()
	if err != nil {
		log.Fatalf("Error Loading Database Connection")
	}
	_, err = dbConn.InitNoSQL()
	if err != nil {
		log.Fatalf("Error Loading Redis Connection")
	}
	token := (&infra.SecretsToken{}).GenerateSecret()

	server := infra.Server{}
	server.NewServer("8000")
	router := routes.ConfigRoutes(server.Server, dbConn, *token)
	log.Fatal(router.Run(":" + server.Port))

}
