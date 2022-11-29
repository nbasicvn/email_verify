package main

import (
	"github.com/truemail-rb/truemail-go"
	"log"
	"net/mail"
	"os"
)
import "github.com/gin-gonic/gin"
import "net/http"
import "github.com/joho/godotenv"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	configuration, _ := truemail.NewConfiguration(
		truemail.ConfigurationAttr{
			VerifierEmail: "verifier@htpland.com",
			SmtpFailFast:  true,
		},
	)

	r := gin.Default()
	r.GET("/email/verify/:email", func(c *gin.Context) {
		email := c.Param("email")
		if !isValidEmailAddress(email) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "Email address invalid",
			})
			return
		}
		status := truemail.IsValid(email, configuration)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"status": status,
			},
		})
	})
	_ = r.Run(":" + os.Getenv("PORT"))
}

func isValidEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
