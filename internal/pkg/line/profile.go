package line

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type UserProfile struct {
	UserId        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Language      string `json:"language"`
}

func GetUserProfile(userId string) (*UserProfile, error) {
	err := godotenv.Load("../pkg/.env")

	url := fmt.Sprintf("https://api.line.me/v2/bot/profile/%s", userId)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CHANNEL_TOKEN")))
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	defer res.Body.Close()

	userProfile := &UserProfile{}
	err = json.NewDecoder(res.Body).Decode(userProfile)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
