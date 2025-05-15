package services

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"net/http"
	"strconv"
	"strings"
)

func Signup(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		//var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user administrativo.Usuario
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if user.Email != "" {
			existUser, _ := Get[administrativo.Usuario](&user, map[string]interface{}{"diretor": true, "ativo": true}, conn)
			if existUser == nil {
				existUser, _ = Get[administrativo.Usuario](&user, map[string]interface{}{"coordenador": true, "ativo": true}, conn)
				fmt.Println(existUser)
				if existUser != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "This Action is disable !!!"})
					c.Abort()
					return
				} else {
					foundUser, _ := Get[administrativo.Usuario](&user, map[string]interface{}{"email": user.Email, "ativo": true}, conn)
					if foundUser != nil {
						c.JSON(http.StatusConflict, gin.H{"error": "User Registred"})
						c.Abort()
						return
					}
				}
			}
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid User !!!"})
			c.Abort()
			return
		}
		userID, err := New[administrativo.Usuario](&user, conn)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		user.ID = userID
		token, err := config.CreateToken(user, 1440, token.GetAccess())
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		err = config.StoreToken(token, strconv.Itoa(int(userID)), 1440, conn)
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
func SignupStatus(conn config.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user administrativo.Usuario
		existUser, _ := Get[administrativo.Usuario](&user, map[string]interface{}{"diretor": true, "ativo": true}, conn)
		if existUser == nil {
			existUser, _ = Get[administrativo.Usuario](&user, map[string]interface{}{"coordenador": true, "ativo": true}, conn)
			fmt.Println(existUser)
			//Verificar a senha com o user
			if existUser == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This Action is disable !!!"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusAccepted, gin.H{"data": "Accept !!!"})
				c.Abort()
				return
			}
		}
	}
}
func Login(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user administrativo.Usuario
		jsonMap := make(map[string]interface{})
		err := json.NewDecoder(c.Request.Body).Decode(&jsonMap)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		password := jsonMap["code"].(string)
		foundUser, err := Get[administrativo.Usuario](&user, map[string]interface{}{"email": jsonMap["email"].(string), "ativo": true}, conn)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			c.Abort()
			return
		}
		err = administrativo.VerifyPassword(conn.Db, foundUser.ID, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			c.Abort()
			return
		}
		//Create an access token
		accesstoken, err := config.CreateToken(*foundUser, 1440, token.GetAccess())
		err = config.StoreToken(accesstoken, strconv.Itoa(int(foundUser.ID)), 1440, conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//Create an refresh token
		refreshtoken, err := config.CreateToken(*foundUser, 180, token.GetRefresh())
		err = config.StoreToken(strconv.Itoa(int(foundUser.ID)), refreshtoken, 180, conn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var perfil string
		if foundUser.Diretor == true {
			perfil += "Diretor; "
		}
		if foundUser.Coordenador == true {
			perfil += "Coordenador; "
		}
		if foundUser.Professor == true {
			perfil += "Professor"
		}

		c.JSON(http.StatusOK, gin.H{"name": foundUser.Nome, "email": foundUser.Email, "profile": perfil, "data": refreshtoken})

	}
}
func Logout(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		userID := c.Request.Header.Get("ID")
		if authHeader == "" || userID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid authorization header format"})
			return
		}
		// Split the header value into 'Bearer' and the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		// Verify the JWT token
		tokenString := tokenParts[1]
		_, err := config.VerifyToken(tokenString, token.GetRefresh())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if err := config.RevokeToken(tokenString, conn); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
			return
		}
		if err := config.RevokeToken(userID, conn); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke user session"})
			return
		}
		c.Redirect(http.StatusFound, "/login")
	}
}
func Authenticate(conn config.Connection, token config.SecretsToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		userID := c.Request.Header.Get("ID")
		if authHeader == "" || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Empty Authorization"})
			c.Abort()
			return
		}
		// Split the header value into 'Bearer' and the token
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization"})
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
		foundUser, err := Get[administrativo.Usuario](new(administrativo.Usuario), map[string]interface{}{"email": userID, "ativo": true}, conn)
		//		foundUser := (*foundUsers)[0]
		token, err := ValidateSession(conn, tokenString, token, *foundUser)

		// Store the user ID in the request context
		c.Set("id", foundUser.ID)
		c.Set("data", token)

		// Proceed to the next handler
		c.Next()

	}
}
func ValidateSession(conn config.Connection, tokenString string, token config.SecretsToken,
	user administrativo.Usuario) (string, error) {
	if conn.NoSql != nil {
		return "", fmt.Errorf("error Database access token")
	}
	userId, err := conn.NoSql.Get(tokenString).Result()
	if err != nil {
		return "", fmt.Errorf("error validate session %w", err)
	}
	if userId == "" {
		tokenAccess, err := conn.NoSql.Get(strconv.Itoa(int(user.ID))).Result()
		if err != nil || tokenAccess == "" {
			return "", fmt.Errorf("error validate session %w", err)
		}
		tk, err := config.VerifyToken(tokenAccess, token.GetAccess())
		if tk == nil {
			_ = config.RevokeToken(strconv.Itoa(int(user.ID)), conn)
			return "", fmt.Errorf("error validate session %w", err)
		}
		refreshtk, err := config.CreateToken(user, 10, token.GetRefresh())
		if err != nil {
			return "", fmt.Errorf("error validate session %w", err)
		}
		err = config.StoreToken(refreshtk, strconv.Itoa(int(user.ID)), 10, conn)
		if err != nil {
			return "", fmt.Errorf("error validate session %w", err)
		}
	}
	if userId == strconv.Itoa(int(user.ID)) {
		tk, err := conn.NoSql.Get(strconv.Itoa(int(user.ID))).Result()
		if err != nil || tk == "" {
			_ = config.RevokeToken(tokenString, conn)
			_ = config.RevokeToken(tk, conn)
			return "", fmt.Errorf("error validate session %w", err)
		}
	}
	return tokenString, nil
}
