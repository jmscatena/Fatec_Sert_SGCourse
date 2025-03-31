package config

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"strconv"
	"time"
)

func CreateToken(user administrativo.Usuario, expire int, secretkey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uID":  strconv.Itoa(int(user.ID)), // Subject (user identifier)
		"name": user.Nome,                  // Issuer
		//"exp":  time.Now().Add(time.Hour).Unix(), // Expiration time
		"exp": time.Now().Add(time.Duration(expire) * time.Minute).Unix(), // Expiration time
		//"aud":  getRole(username),                // Audience (user role)
		//"iat":  time.Now().Unix(),                // Issued at
	})
	tokenString, err := claims.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func StoreToken(key string, value string, expire int, conn Connection) error {
	redisClient := conn.NoSql
	if redisClient == nil {
		return fmt.Errorf("Error Access Database to store token")
	}
	err := redisClient.Set(key, value, time.Minute*time.Duration(expire)).Err()
	if err != nil {
		return fmt.Errorf("Error storing token: %w", err)
	}
	//defer redisClient.Close()
	return nil
}

func VerifyToken(tokenString string, secretKey string) (*jwt.Token, error) {
	sk := secretKey
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		ok := token.Method.Alg() == jwt.SigningMethodHS256.Alg()
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sk), nil
	}
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, keyfunc)
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, fmt.Errorf("Invalid Token")
	}
	return jwtToken, nil
	/*
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		println(token.Valid)
		if !token.Valid {
			return nil, fmt.Errorf("invalid token")
		}
		return token, nil
	*/
}

func RevokeToken(token string, conn Connection) error {
	//redisClient, err := database.InitDF()
	if conn.NoSql == nil {
		return fmt.Errorf("Error Database Connection to Data token revoke")
	}
	err := conn.NoSql.Del(token).Err()
	if err != nil {
		return fmt.Errorf("Error revoke token: %w", err)
	}
	//defer redisClient.Close()
	return nil
}
