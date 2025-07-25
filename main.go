package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"moodleinix/constant"
	"moodleinix/model"
	"moodleinix/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/api/register/JDA", service.BasicAuthMiddleware(constant.BasicAuthUserName, constant.BasicAuthUPassword), func(c *gin.Context) {
		var req model.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedURL, err := url.Parse(req.ClassLink)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classLink URL"})
			return
		}

		queryParams := parsedURL.Query()
		courseIDStr := queryParams.Get("id")
		if courseIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "course ID not found in classLink"})
			return
		}

		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}

		// Cek user berdasarkan email apakah sudah ada
		userID, err := service.GetUserIDByEmail(req.Email)
		if err != nil {
			// user belum ada, buat user baru
			username := service.GetUsernameFromEmail(req.Email)
			firstname, lastname := service.SplitName(req.Name)

			user := model.MoodleUser{
				Username:  username,
				Password:  req.Password,
				Firstname: firstname,
				Lastname:  lastname,
				Email:     req.Email,
				City:      req.City,
				Department: req.GroupName,
			}

			userID, err = service.CreateMoodleUser(user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Enroll user ke course
		if err := service.EnrolUser(userID, courseID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if req.GroupName != "" {
			log.Printf("[INFO] Permintaan penambahan user ID %d ke group '%s' di course ID %d", userID, req.GroupName, courseID)

			groupID, err := service.GetGroupIDByName(req.GroupName, courseID)
			if err != nil {
				log.Printf("[ERROR] Group not found: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "group not found: " + err.Error()})
				return
			}

			err = service.AddUserToGroup(userID, groupID)
			if err != nil {
				log.Printf("[ERROR] Failed to add user to group: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add user to group: " + err.Error()})
				return
			}
			log.Printf("[INFO] User ID %d berhasil ditambahkan ke group ID %d", userID, groupID)
		}

		message := "User registered, enrolled, and processed add user to group successfully"

		now := time.Now().UTC()

		response := model.ReturnResponse{
			Data: model.ResponseData{
				ID:        userID,
				Name:      req.Name,
				Email:     req.Email,
				CreatedAt: now,
				UpdatedAt: now,
				IsExists:  true,
				Token:     "-",
			},
			Message: message,
			Status:  true,
		}

		c.JSON(http.StatusOK, response)
	})

	log.Fatal(r.Run(constant.HOST))
}
