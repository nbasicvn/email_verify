package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
)
import "github.com/gin-gonic/gin"
import "net/http"
import "github.com/joho/godotenv"
import emailverifier "github.com/AfterShip/email-verifier"

var (
	verifier = emailverifier.
		NewVerifier().
		EnableSMTPCheck().
		EnableCatchAllCheck().
		EnableGravatarCheck().
		EnableAutoUpdateDisposable()
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

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
		ret, err := verifier.Verify(email)
		if err != nil {
			fmt.Println("verify email address failed, error is: ", err)
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "Verify email address failed",
			})
			return
		}
		if !ret.Syntax.Valid {
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "Email address invalid",
			})
			return
		}
		fmt.Println(ret)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"status": ret.Reachable != "no",
				"info":   ret,
			},
		})
	})
	_ = r.Run(":" + os.Getenv("PORT"))
}

func isValidEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
