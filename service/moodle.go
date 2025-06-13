package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"

	"moodleinix/constant"
	"moodleinix/model"
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
func EnrolUser(userid int) error {
	form := url.Values{}
	form.Set("wstoken", constant.MoodleToken)
	form.Set("wsfunction", constant.WSFunctionEnrol)
	form.Set("moodlewsrestformat", "json")

	form.Set("enrolments[0][roleid]", fmt.Sprintf("%d", constant.RoleID))
	form.Set("enrolments[0][userid]", fmt.Sprintf("%d", userid))
	form.Set("enrolments[0][courseid]", fmt.Sprintf("%d", constant.CourseID))

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
