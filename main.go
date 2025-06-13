package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"moodleinix/constant"
	"moodleinix/model"
	"moodleinix/service"
)

func main() {
	r := gin.Default()

	r.POST("/api/register/JDA", func(c *gin.Context) {
		var req model.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		firstname, lastname := service.SplitName(req.Name)
		username := service.GetUsernameFromEmail(req.Email)

		user := model.MoodleUser{
			Username:  username,
			Password:  req.Password,
			Firstname: firstname,
			Lastname:  lastname,
			Email:     req.Email,
		}

		userID, err := service.CreateMoodleUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := service.EnrolUser(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		now := time.Now().UTC()

		response := model.ReturnResponse{
			Data: model.ResponseData{
				ID:        userID,
				Name:      req.Name,
				Email:     req.Email,
				CreatedAt: now,
				UpdatedAt: now,
				IsExists:  true,
				Token:     constant.MoodleToken,
			},
			Message: "User registered and enrolled successfully",
			Status:  true,
		}

		c.JSON(http.StatusOK, response)
	})

	log.Fatal(r.Run(constant.HOST))
}
