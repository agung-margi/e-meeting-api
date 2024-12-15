package main

import (
	"e-meeting-api/presenter/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// tokenString := "YOUR_JWT_TOKEN_HERE"
	// secret := "JWT_SECRET" // Nama environment variable yang menyimpan secret key

	// // Verifikasi token
	// token, err := auth.VerifyToken(secret, tokenString)
	// if err != nil {
	// 	log.Fatalf("Error verifying token: %v", err)
	// }

	// fmt.Println("Token is valid. Claims:", token.Claims)

	if err := handler.RoutingRestAPI(e); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":8080"))
}
