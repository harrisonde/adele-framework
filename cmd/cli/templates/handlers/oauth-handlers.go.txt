package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"myapp/data"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
)

func (h *Handlers) PasswordGrantExchange(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Validate the client id and secret is valid
	clientId := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	cID, err := strconv.Atoi(clientId)
	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Invalid client provided."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	client, err := h.Models.Clients.Get(cID)
	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Invalid client provided."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	isValidClient := client.CheckIsValid(cID, clientSecret)
	if !isValidClient {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Invalid client provided."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Validate the user
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Invalid user provided."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Validate the users credentials
	matches, err := user.PasswordMatches(password)
	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Error validating password."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	if !matches {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Invalid password"
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Generate a token for the user and persist in the database
	expiry := 24 * time.Hour
	token, err := h.Models.Tokens.GenerateToken(user.ID, expiry)
	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Error creating token."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Persist the token
	err = h.Models.Tokens.Insert(*token, *user)

	if err != nil {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = true
		payload.Message = "Error creating token."
		h.App.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	var payload struct {
		TokenType   string        `json:"token_type"`
		ExpiresIn   time.Duration `json:"expires_in"`
		AccessToken string        `json:"access_token"`
	}

	payload.TokenType = "Bearer"
	payload.ExpiresIn = expiry
	payload.AccessToken = token.PlainText
	h.App.WriteJSON(w, http.StatusOK, payload)
}

func (h *Handlers) InitSocialAuth() {
	scope := []string{"user"}
	gScope := []string{"email", "profile"}
	mScope := []string{"openid", "offline_access", "user.read"}

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK"), scope...),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_CALLBACK"), gScope...),
		microsoftonline.New(os.Getenv("MICROSOFT_KEY"), os.Getenv("MICROSOFT_SECRET"), os.Getenv("MICROSOFT_CALLBACK"), mScope...),
	)

	key := os.Getenv("KEY")
	maxAge := 86400 * 30
	st := sessions.NewCookieStore([]byte(key))
	st.MaxAge(maxAge)
	st.Options.Path = "/"
	st.Options.HttpOnly = true
	st.Options.Secure = false

	gothic.Store = st
}

func (h *Handlers) SocialLogin(w http.ResponseWriter, r *http.Request) {
	// Get provider and put it into the session so can logout with it later
	provider := chi.URLParam(r, "provider")
	h.App.Session.Put(r.Context(), "social_provider", provider)
	h.InitSocialAuth()

	// Is the user already logged into the system?
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handlers) SocialMediaCallback(w http.ResponseWriter, r *http.Request) {
	h.InitSocialAuth()
	gUser, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		h.App.Session.Put(r.Context(), "error", err.Error())
		h.App.ErrorLog.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Look up the user with their email address
	var u data.User
	var testUser *data.User

	testUser, err = u.GetByEmail(strings.ToLower(gUser.Email))

	h.App.InfoLog.Println("SocialMediaCallback run and gUser email is: " + strings.ToLower(gUser.Email))
	h.App.InfoLog.Println(testUser)

	if err != nil {
		log.Println(err)

		// Get provider and cast to string
		provider := h.App.Session.Get(r.Context(), "social_provider").(string)

		// we do not have a user
		var newUser data.User
		switch provider {
		case "github":
			exploded := strings.Split(gUser.Name, " ")
			newUser.FirstName = exploded[0]
			if len(exploded) > 1 {
				newUser.LastName = exploded[1]
			}
		case "google":
			newUser.FirstName = gUser.FirstName
			newUser.LastName = gUser.LastName
		}
		newUser.Email = strings.ToLower(gUser.Email)
		newUser.Active = 1
		newUser.Password = h.randomString(20)
		newUser.CreatedAt = time.Now()
		newUser.UpdatedAt = time.Now()

		// Insert the user into the database
		_, err := newUser.Insert(newUser)
		if err != nil {
			// TO DO: Do we want to redirect to login with flash?
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		testUser, _ = u.GetByEmail(strings.ToLower(gUser.Email))
	}

	// Sign the user into the system
	h.App.Session.Put(r.Context(), "userID", testUser.ID)
	h.App.Session.Put(r.Context(), "social_token", gUser.AccessToken)
	h.App.Session.Put(r.Context(), "social_email", gUser.Email)

	// Redirect to the application root
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) socialLogout(w http.ResponseWriter, r *http.Request) string {
	h.App.ErrorLog.Println("social logout")

	provider, ok := h.App.Session.Get(r.Context(), "social_provider").(string)

	h.App.ErrorLog.Println(!ok, provider)

	if !ok {
		return ""
	}

	// Call the proper provider and revoke the auth token.
	switch provider {
	case "github":
		clientID := os.Getenv("GITHUB_KEY")
		clientSecret := os.Getenv("GITHUB_SECRET")

		token := h.App.Session.Get(r.Context(), "social_token").(string)

		var payload struct {
			AccessToken string `json:"access_token"`
		}
		payload.AccessToken = token

		jsonReq, err := json.Marshal(payload)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return ""
		}
		req, err := http.NewRequest(http.MethodDelete,
			fmt.Sprintf("https://%s:%s@api.github.com/applications/%s/grant", clientID, clientSecret, clientID),
			bytes.NewBuffer(jsonReq))
		if err != nil {
			h.App.ErrorLog.Println(err)
			return ""
		}

		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			h.App.ErrorLog.Println("Error logging out of Github:", err)
			return ""
		}

	case "google":
		token := h.App.Session.Get(r.Context(), "social_token").(string)
		_, err := http.PostForm(fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?%s", token), nil)
		if err != nil {
			h.App.ErrorLog.Println("Error logging out of Google:", err)
			return ""
		}

	case "microsoftonline":
		return fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/logout?post_logout_redirect_uri=%s", h.App.Server.URL+"/login")
	}

	return ""
}
