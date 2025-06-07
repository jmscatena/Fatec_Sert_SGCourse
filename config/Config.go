package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/migrations"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Connection struct {
	Db    *gorm.DB
	NoSql *redis.Client
}

func (c *Connection) InitDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading Configuration File")
	}

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbase := os.Getenv("DB")
	dbServer := os.Getenv("DBSERVER")
	dbPort := os.Getenv("DBPORT")
	dbURL := "postgres://" + dbUser + ":" + dbPass + "@" + dbServer + ":" + dbPort + "/" + dbase

	c.Db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln("Erro no carregamento do SGBD", err)
	}
	migrations.RunMigrate(c.Db)
	return c.Db, err
}

func (c *Connection) InitNoSQL() (*redis.Client, error) {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	c.NoSql = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	var _, err = c.NoSql.Ping().Result()
	if err != nil {
		log.Fatalln("Erro no carregamento do Redis:", err)
	}
	return c.NoSql, nil
}

type Server struct {
	Port   string
	Server *gin.Engine
}

func IPMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém o endereço de IP do cliente.
		// c.ClientIP() tenta extrair o IP de cabeçalhos como X-Forwarded-For ou X-Real-IP,
		// o que é útil se a sua API estiver atrás de um proxy reverso ou load balancer.
		clientIP := c.ClientIP()

		// Verifica se o IP do cliente está na lista de IPs permitidos.
		allowed := false
		for _, ip := range allowedIPs {
			if clientIP == ip {
				allowed = true
				break
			}
		}

		// Se o IP não for permitido, a requisição é abortada.
		if !allowed {
			// c.AbortWithStatusJSON interrompe a cadeia de execução e envia uma resposta JSON.
			// http.StatusForbidden é o código de status HTTP 403.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Acesso negado",
			})
			return // Encerra a função aqui.
		}

		// Se o IP for permitido, a requisição continua para o próximo handler.
		c.Next()
	}
}
func (s *Server) NewServer(port string) {
	allowedIPs := []string{"127.0.0.1", "200.144.13.38"}
	s.Port = port
	s.Server = gin.Default()
	s.Server.Use(IPMiddleware(allowedIPs))
}

func (s *Server) Run() {
	log.Printf("Server running at port: %v", s.Port)

}

type SecretsToken struct {
	secret  string
	refresh string
}

func (t *SecretsToken) GetAccess() string {
	return t.secret
}
func (t *SecretsToken) GetRefresh() string {
	return t.refresh
}

func (t *SecretsToken) GenerateSecret() *SecretsToken {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading Configuration File")
	}
	return &SecretsToken{
		secret:  os.Getenv("TOKEN_SECRET_KEY"),
		refresh: os.Getenv("REFRESH_SECRET_KEY"),
	}
}
