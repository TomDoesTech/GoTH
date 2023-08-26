package config

import (
	"crypto/rsa"
	"fmt"
	"log"

	b64 "encoding/base64"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string

	JWTPrivateKey *rsa.PrivateKey
	JWTPublicKey  *rsa.PublicKey
}

func Must() Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file", err)
	}

	viper.AutomaticEnv()

	port := viper.GetString("PORT")

	if port == "" {
		port = ":8080"
	}

	// make sure port starts with a colon
	if port[0] != ':' {
		port = ":" + port
	}

	JWTPrivateKeyEncoded := viper.GetString("JWT_PRIVATE_KEY")
	JWTPrivateKeyStr, err := b64.URLEncoding.DecodeString(JWTPrivateKeyEncoded)
	if err != nil {
		log.Fatal("Error decoding JWT_PRIVATE_KEY", err)
	}

	JWTPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(JWTPrivateKeyStr))
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		log.Fatal("Error parsing private key:", err)
	}

	JWTPublicKeyEncoded := viper.GetString("JWT_PUBLIC_KEY")
	JWTPublicKeyStr, err := b64.URLEncoding.DecodeString(JWTPublicKeyEncoded)
	if err != nil {
		log.Fatal("Error decoding JWT_PUBLIC_KEY", err)
	}

	JWTPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(JWTPublicKeyStr))
	if err != nil {
		log.Fatal("Error parsing public key:", err)
	}

	return Config{
		DBHost:        viper.GetString("DATABASE_HOST"),
		DBUser:        viper.GetString("DATABASE_USER"),
		DBPassword:    viper.GetString("DATABASE_PASSWORD"),
		DBName:        viper.GetString("DATABASE_NAME"),
		JWTPrivateKey: JWTPrivateKey,
		JWTPublicKey:  JWTPublicKey,
		Port:          port,
	}
}
