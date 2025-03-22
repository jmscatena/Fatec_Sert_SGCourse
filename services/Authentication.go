package services

import (
	"encoding/json"
	"fmt"
	gin "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"net/http"
	"strings"
)

var validate = validator.New()

/*
	func Auth() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Extract the JWT token from the Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.String(http.StatusUnauthorized, "Authorization header is missing")
				c.Abort()
				return
			}

			// Split the header value into "Bearer " and the token
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				c.String(http.StatusUnauthorized, "Invalid authorization header format")
				c.Abort()
				return
			}

			// Verify the JWT token
			tokenString := tokenParts[1]
			userId, err := VerifyToken(tokenString, "refresh")
			if err != nil {
				c.String(http.StatusUnauthorized, "Invalid or expired token")
				c.Abort()
				return
			}

			// Store the user ID in the request context
			c.Set("userId", userId)

			// Proceed to the next handler
			c.Next()
		}
	}
*/
func Signup(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		//var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user administrativo.Usuario
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}
		foundUser, err := Get[administrativo.Usuario](&user, "email=?", user.Email, conn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			c.Abort()
			return
		}
		if foundUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User Registred"})
			c.Abort()
			return
		}

		userID, err := New[administrativo.Usuario](&user, conn)
		user.UID = userID
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}
		token, err := config.CreateToken(user, 1440, token.GetAccess())
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		err = config.StoreToken(token, userID.String(), 1440, conn)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Could Not Signup."})
			c.Abort()
			return

		}
		c.JSON(http.StatusOK, gin.H{"data": userID, "token": token})
		//defer cancel()
		//Redirect to login if successfull
		c.Next()
	}

}
func Login(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user administrativo.Usuario
		json_map := make(map[string]interface{})
		err := json.NewDecoder(c.Request.Body).Decode(&json_map)
		if err != nil {
			//c.ShouldBindJSON(&json_map);
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		println(json_map["email"])
		println(json_map["code"])
		password := json_map["code"].(string)
		condition := "Email=?"
		foundUser, err := Get[administrativo.Usuario](&user, condition, json_map["email"].(string), conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			c.Abort()
			return
		}
		println(foundUser.Nome)
		println(password)
		//foundUser := (*foundUsers)[0]
		err = administrativo.VerifyPassword(password, foundUser.Senha)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			c.Abort()
			return
		}
		//Create access token
		accesstoken, err := config.CreateToken(*foundUser, 1440, token.GetAccess())
		err = config.StoreToken(accesstoken, foundUser.UID.String(), 1440, conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//Create refresh token
		refreshtoken, err := config.CreateToken(*foundUser, 180, token.GetRefresh())
		err = config.StoreToken(foundUser.UID.String(), refreshtoken, 180, conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": foundUser.UID, "data": refreshtoken})

	}
}
func Logout(conn config.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		userID := c.Request.Header.Get("ID")
		if authHeader == "" || userID == "" {
			c.Redirect(http.StatusFound, "/login")
		}
		// Split the header value into "Bearer " and the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		// Verify the JWT token
		tokenString := tokenParts[1]
		config.RevokeToken(tokenString, conn)
		config.RevokeToken(userID, conn)
		c.Redirect(http.StatusFound, "/login")
	}
}
func Authenticate(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		userID := c.Request.Header.Get("ID")

		if authHeader == "" || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}
		// Split the header value into "Bearer " and the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
			c.Abort()
			return
		}
		// Verify the JWT token
		tokenString := tokenParts[1]
		_, err := config.VerifyToken(tokenString, token.GetRefresh())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired Access"})
			c.Abort()
			return
		}
		condition := "UID=?"
		foundUser, err := Get[administrativo.Usuario](new(administrativo.Usuario), condition, userID, conn)
		//		foundUser := (*foundUsers)[0]
		token, err := ValidateSession(conn, tokenString, token, *foundUser)

		// Store the user ID in the request context
		//c.JSON(http.StatusOK, gin.H{"data": foundUser.UID, "token": token})
		c.Set("id", foundUser.UID)
		c.Set("data", token)

		// Proceed to the next handler
		c.Next()

	}
}

func ValidateSession(conn config.Connection, tokenString string, token config.SecretsToken,
	user administrativo.Usuario) (string, error) {

	if conn.NoSql != nil {
		return "", fmt.Errorf("Error Database access token")
	}
	userId, err := conn.NoSql.Get(tokenString).Result()
	if err != nil {
		return "", fmt.Errorf("Error validate session: %w", err)
	}
	if userId == "" {
		tokenAccess, err := conn.NoSql.Get(user.UID.String()).Result()
		if err != nil || tokenAccess == "" {
			return "", fmt.Errorf("Error validate session: %w", err)
		}
		tk, err := config.VerifyToken(tokenAccess, token.GetAccess())
		if tk == nil {
			config.RevokeToken(user.UID.String(), conn)
			return "", fmt.Errorf("Error validate session: %w", err)
		}
		refreshtk, err := config.CreateToken(user, 10, token.GetRefresh())
		if err != nil {
			return "", fmt.Errorf("Error validate session: %w", err)
		}
		err = config.StoreToken(refreshtk, user.UID.String(), 10, conn)
		if err != nil {
			return "", fmt.Errorf("Error validate session: %w", err)
		}
	}
	if userId == user.UID.String() {
		tk, err := conn.NoSql.Get(user.UID.String()).Result()
		if err != nil || tk == "" {
			config.RevokeToken(tokenString, conn)
			config.RevokeToken(tk, conn)
			return "", fmt.Errorf("Error validate session: %w", err)
		}
	}
	//defer conn.NoSql.Close()
	return tokenString, nil
}
