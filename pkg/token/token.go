package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// AuthRequest represents the request payload for authentication
type AuthRequest struct {
	Auth Auth `json:"auth"`
}

// Auth contains identity and scope for the authentication
type Auth struct {
	Identity Identity `json:"identity"`
	Scope    Scope    `json:"scope"`
}

// Identity contains methods and password information for authentication
type Identity struct {
	Methods  []string `json:"methods"`
	Password Password `json:"password"`
}

// Password contains user details for password authentication
type Password struct {
	User User `json:"user"`
}

// User contains username, domain, and password
type User struct {
	Name     string `json:"name"`
	Domain   Domain `json:"domain"`
	Password string `json:"password"`
}

// Domain contains domain name
type Domain struct {
	Name string `json:"name"`
}

// Scope contains project information
type Scope struct {
	Project Project `json:"project"`
}

// Project contains project name and domain
type Project struct {
	Name   string `json:"name"`
	Domain Domain `json:"domain"`
}

// GetAuthToken retrieves the authentication token from Selectel
func GetAuthToken() (string, error) {
	username := os.Getenv("SELECTEL_USER_NAME")
	password := os.Getenv("SELECTEL_USER_PASSWORD")
	accountID := os.Getenv("SELECTEL_DOMAIN_NAME")
	projectName := os.Getenv("SELECTEL_PROJECT_NAME")

	if username == "" || password == "" || accountID == "" || projectName == "" {
		return "", errors.New("one or more required environment variables are missing")
	}

	authRequest := AuthRequest{
		Auth: Auth{
			Identity: Identity{
				Methods: []string{"password"},
				Password: Password{
					User: User{
						Name:     username,
						Domain:   Domain{Name: accountID},
						Password: password,
					},
				},
			},
			Scope: Scope{
				Project: Project{
					Name:   projectName,
					Domain: Domain{Name: accountID},
				},
			},
		},
	}

	requestBody, err := json.Marshal(authRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://cloud.api.selcloud.ru/identity/v3/auth/tokens", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", errors.New("X-Subject-Token header is missing in the response")
	}

	return token, nil
}
