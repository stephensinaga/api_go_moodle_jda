package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"bytes"
	"encoding/csv"

	"moodleinix/constant"
	"moodleinix/model"
	"moodleinix/service"
	"moodleinix/database"
	"moodleinix/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	defer database.Close()
	
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

	// Endpoint untuk mengakses data user course quiz dalam format JSON
	r.GET("/api/user-course-quiz-data", func(c *gin.Context) {
		data, err := service.GetUserCourseQuizData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	// Endpoint baru untuk mengunduh data user course quiz dalam format CSV
	r.GET("/api/user-course-quiz-data/csv", middleware.TokenAuthMiddleware, func(c *gin.Context) {
		data, err := service.GetUserCourseQuizData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		b := &bytes.Buffer{}
		writer := csv.NewWriter(b)

		header := []string{
			"UserID", "FullName", "Email", "City", "Department", "CourseID", "CourseName",
			"QuizDNSDHCP", "QuizRouting", "QuizIPAddress", "QuizOSIModel", "QuizTopologiJaringan",
			"QuizAncamanDuniaDigital", "QuizHukumAturanDuniaDigital", "QuizTujuanMetodeEthicalHacking",
			"QuizDasarEthicalHackingTipeHacker", "QuizUUPDP", "QuizUUITE", "QuizCyberWarfare",
			"QuizKerentananCybersecurity", "QuizSumberAncamanCybersecurity", "QuizTipeAncamanCybersecurity",
			"TestJaringanKomputerDasar", "TestAncamanDuniaDigital", "TestHukumAturandalamDuniaDigital","TestKonsepEthicalHacking", "TestAkhir", "Durasi Penyelesaian (jam)",
			"Tanggal Selesai",
		}
		if err := writer.Write(header); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
			return
		}

		floatToString := func(f *float64) string {
			if f == nil {
				return ""
			}
			return strconv.FormatFloat(*f, 'f', 2, 64)
		}

		timestampToString := func(ts *float64) string {
			if ts == nil || *ts == 0 {
				return ""
			}
			sec := int64(*ts)
			t := time.Unix(sec, 0)
			return t.Format("2006-01-02 15:04:05")
		}

		for _, row := range data {
			record := []string{
				strconv.Itoa(row.UserID),
				row.Fullname,
				row.Email,
				row.City,
				row.Department,
				strconv.Itoa(row.CourseID),
				row.CourseName,
				floatToString(row.QuizDNSDHCP),
				floatToString(row.QuizRouting),
				floatToString(row.QuizIPAddress),
				floatToString(row.QuizOSIModel),
				floatToString(row.QuizTopologiJaringan),
				floatToString(row.QuizAncamanDuniaDigital),
				floatToString(row.QuizHukumAturanDuniaDigital),
				floatToString(row.QuizTujuanMetodeEthicalHacking),
				floatToString(row.QuizDasarEthicalHackingTipeHacker),
				floatToString(row.QuizUUPDP),
				floatToString(row.QuizUUITE),
				floatToString(row.QuizCyberWarfare),
				floatToString(row.QuizKerentananCybersecurity),
				floatToString(row.QuizSumberAncamanCybersecurity),
				floatToString(row.QuizTipeAncamanCybersecurity),
				floatToString(row.TestJaringanKomputerDasar),
				floatToString(row.TestAncamanDuniaDigital),
				floatToString(row.TestHukumAturandalamDuniaDigital),	
				floatToString(row.TestKonsepEthicalHacking),
				floatToString(row.TestAkhir),
				floatToString(row.CourseCompletionDurationHours),
				timestampToString(row.EndDate), // Panggil konversi untuk EndDate di sini
			}

			if err := writer.Write(record); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV record"})
				return
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "CSV write error"})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename=report_user_LMS_inixindo.csv.csv")
		c.Data(http.StatusOK, "text/csv", b.Bytes())
	})

	log.Fatal(r.Run(constant.HOST))

	// Production
	// 	log.Fatal(r.RunTLS(
	// 	constant.PORT,
	// 	constant.EncryptFullChain,
	// 	constant.EncryptPriveKey,
	// ))
}
