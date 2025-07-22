package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"moodleinix/constant"
	"moodleinix/model"

	"github.com/gin-gonic/gin"
)

// SplitName memisahkan nama lengkap jadi firstname dan lastname
func SplitName(fullName string) (string, string) {
	names := strings.Fields(fullName)
	count := len(names)

	switch count {
	case 0:
		return "User", "-"
	case 1:
		return names[0], "-"
	case 2:
		return names[0], names[1]
	default:
		return names[0], names[count-1]
	}
}

// GetUsernameFromEmail mengambil bagian awal email sebelum '@'
func GetUsernameFromEmail(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return email
	}
	return email[:atIndex]
}

func CreateMoodleUser(user model.MoodleUser) (int, error) {
	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionCreate)
	form.Set("moodlewsrestformat", "json")

	form.Set("users[0][username]", user.Username)
	form.Set("users[0][password]", user.Password)
	form.Set("users[0][firstname]", user.Firstname)
	form.Set("users[0][lastname]", user.Lastname)
	form.Set("users[0][email]", user.Email)

	resp, err := http.PostForm(constant.MoodleURL, form)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("gagal membaca response body: %v", err)
	}

	fmt.Println("HTTP Status:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println("Response Body:", string(body))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return 0, fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	// Parsing JSON
	var responseArray []map[string]interface{}
	var responseObject map[string]interface{}

	err = json.Unmarshal(body, &responseArray)
	if err != nil {
		err2 := json.Unmarshal(body, &responseObject)
		if err2 == nil {
			return 0, fmt.Errorf("moodle error: %v", responseObject["message"])
		}
		return 0, fmt.Errorf("gagal parsing respons: %v", err)
	}

	if len(responseArray) == 0 {
		return 0, fmt.Errorf("tidak ada user yang dibuat")
	}

	userID, ok := responseArray[0]["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("format user ID tidak valid")
	}

	return int(userID), nil

}

// EnrolUser mendaftarkan user ke course di Moodle
func EnrolUser(userid, courseID int) error {
	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionEnrol)
	form.Set("moodlewsrestformat", "json")

	form.Set("enrolments[0][roleid]", fmt.Sprintf("%d", constant.RoleID))
	form.Set("enrolments[0][userid]", fmt.Sprintf("%d", userid))
	form.Set("enrolments[0][courseid]", fmt.Sprintf("%d", courseID))

	resp, err := http.PostForm(constant.MoodleURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("enrol gagal: %s", string(body))
	}

	// Optional: bisa parse response jika perlu
	fmt.Println("Enrol response:", string(body))
	return nil
}

// BasicAuthMiddleware adalah middleware untuk memeriksa header Authorization dengan skema Basic Auth
func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Basic ") {
			// Header Authorization tidak ada atau tidak sesuai format Basic
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Decode base64 username:password
		payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		// Split username dan password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != username || pair[1] != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You Haven't authorization"})
			return
		}
		c.Next()
	}
}
func GetUserIDByEmail(email string) (int, error) {
	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionGetUsers)
	form.Set("moodlewsrestformat", "json")
	form.Set("criteria[0][key]", "email")
	form.Set("criteria[0][value]", email)

	log.Printf("Mengirim kriteria pencarian ke Moodle: criteria[0][key]=email, criteria[0][value]=%s\n", email)

	resp, err := http.PostForm(constant.MoodleURL, form)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("gagal membaca response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Moodle response: %s\n", string(body))

	var data struct {
		Users []struct {
			ID int `json:"id"`
		} `json:"users"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("gagal parsing respons: %v", err)
	}

	if len(data.Users) == 0 {
		return 0, fmt.Errorf("user dengan email %s tidak ditemukan", email)
	}

	return data.Users[0].ID, nil
}
func GetGroupIDByName(groupName string, courseID int) (int, error) {
	log.Printf("[INFO] Mencari GroupID untuk nama grup '%s' di course ID: %d", groupName, courseID)

	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionGetNameGroup) // Contoh: "core_group_get_course_groups"
	form.Set("moodlewsrestformat", "json")
	form.Set("courseid", fmt.Sprintf("%d", courseID))

	resp, err := http.PostForm(constant.MoodleURL, form)
	if err != nil {
		log.Printf("[ERROR] Gagal request GetGroupIDByName: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Gagal membaca response body GetGroupIDByName: %v", err)
		return 0, fmt.Errorf("gagal membaca response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Moodle server returned status %d: %s", resp.StatusCode, string(body))
		return 0, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Response Moodle berupa array langsung
	var groups []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &groups)
	if err != nil {
		log.Printf("[ERROR] Gagal parsing respons GetGroupIDByName: %v", err)
		return 0, fmt.Errorf("gagal parsing respons: %v", err)
	}

	for _, group := range groups {
		log.Printf("[DEBUG] Ditemukan grup: ID=%d, Name=%s", group.ID, group.Name)
		if group.Name == groupName {
			log.Printf("[INFO] GroupID ditemukan: %d", group.ID)
			return group.ID, nil
		}
	}

	log.Printf("[WARN] Group dengan nama '%s' tidak ditemukan di course %d", groupName, courseID)
	return 0, fmt.Errorf("group dengan nama %s tidak ditemukan di course %d", groupName, courseID)
}

func AddUserToGroup(userID, groupID int) error {
	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionAddUsertoGroup)
	form.Set("moodlewsrestformat", "json")

	form.Set("members[0][groupid]", strconv.Itoa(groupID))
	form.Set("members[0][userid]", strconv.Itoa(userID))

	resp, err := http.PostForm(constant.MoodleURL, form)
	if err != nil {
		return fmt.Errorf("gagal request ke Moodle: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gagal membaca response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server Moodle mengembalikan status %d: %s", resp.StatusCode, string(body))
	}

	// Opsional: parsing response body untuk cek hasil
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if warning, ok := result["warnings"]; ok && len(warning.([]interface{})) > 0 {
			return fmt.Errorf("warning dari Moodle: %v", warning)
		}
	}

	return nil
}
