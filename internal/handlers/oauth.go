package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/internal/auth"
	"forum/internal/models"
	"forum/internal/services"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	GithubOauthConfig *oauth2.Config
)

func InitOAuth() {
	GithubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
}

const oauthStateCookieName = "oauthstate"

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: oauthStateCookieName, Value: state, Expires: expiration, Path: "/"}
	http.SetCookie(w, &cookie)

	return state
}

func (h *Handlers) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/auth/")
	provider = strings.TrimSuffix(provider, "/login")

	oauthStateString := generateStateOauthCookie(w)

	var url string
	if provider == "github" {
		url = GithubOauthConfig.AuthCodeURL(oauthStateString)
	} else {
		http.Error(w, "Provider not supported", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/auth/")
	provider = strings.TrimSuffix(provider, "/callback")

	oauthState, _ := r.Cookie(oauthStateCookieName)
	if oauthState == nil || r.FormValue("state") != oauthState.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	var oauthConfig *oauth2.Config
	if provider == "github" {
		oauthConfig = GithubOauthConfig
	} else {
		http.Error(w, "Provider not supported", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	var email, id, username string

	if provider == "github" {
		req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
		req.Header.Set("Authorization", "token "+token.AccessToken)
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		var userInfo struct {
			Id    int    `json:"id"`
			Login string `json:"login"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}
		
		email = userInfo.Email
		id = fmt.Sprintf("%d", userInfo.Id)
		username = userInfo.Login
		
		if email == "" {
			// Fetch emails separately for github if private
			reqEmails, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
			reqEmails.Header.Set("Authorization", "token "+token.AccessToken)
			respEmails, err := client.Do(reqEmails)
			if err == nil {
				defer respEmails.Body.Close()
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				json.NewDecoder(respEmails.Body).Decode(&emails)
				for _, e := range emails {
					if e.Primary {
						email = e.Email
						break
					}
				}
			}
		}
	}

	if email == "" {
		http.Error(w, "Failed to retrieve email from provider", http.StatusBadRequest)
		return
	}

	user, err := services.GetUserByEmail(h.DB, email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Create user if not exists
	if user == nil {
		newUser := models.User{
			Username:      username + "_" + id[0:4], // ensure unique username
			Email:         email,
			Password:      "", // No password for oauth
			OAuthProvider: provider,
			OAuthID:       id,
		}
		err = services.CreateUser(h.DB, newUser)
		if err != nil {
			http.Error(w, "Failed to create user account", http.StatusInternalServerError)
			return
		}
		user, _ = services.GetUserByEmail(h.DB, email)
	}

	sessionToken, err := services.CreateSession(h.DB, user.ID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, sessionToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
