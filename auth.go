package gleam

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	redirectURI = "http://localhost:8080/callback"
	authURL     = "https://id.twitch.tv/oauth2/authorize"
	tokenURL    = "https://id.twitch.tv/oauth2/token"
	scopes      = "chat:read chat:edit"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func authenticate(bot *Bot) error {
	fmt.Printf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s\n\n",
		authURL,
		bot.options.ClientID,
		redirectURI,
		url.QueryEscape(scopes))

	codeCh := make(chan string, 1)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			_, err := w.Write([]byte("Got it! Head back to your terminal."))
			if err != nil {
				log.Println("Failed to write response:", err)
			}
		} else {
			if _, err := w.Write([]byte("Failed to receive authorization code.")); err != nil {
				log.Println("Failed to write response:", err)
			}
		}

		go func() { codeCh <- code }()
	})

	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	authCode := <-codeCh
	close(codeCh)

	token, err := getToken(bot.options.ClientID, bot.options.ClientSecret, authCode)
	if err != nil {
		return fmt.Errorf("error exchanging auth code for token: %w", err)
	} else if token.AccessToken == "" {
		return fmt.Errorf("error getting token: empty")
	}

	bot.tokens.access = token.AccessToken
	bot.tokens.refresh = token.RefreshToken

	return nil
}

func getToken(clientID, clientSecret, authCode string) (*tokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", authCode)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	res, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func reauthenticate(bot *Bot) (string, string, error) {
	params := url.Values{}
	params.Add("grant_type", `refresh_token`)
	params.Add("refresh_token", bot.tokens.refresh)
	params.Add("client_id", bot.options.ClientID)
	params.Add("client_secret", bot.options.ClientSecret)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", body)
	if err != nil {
		return "", "", fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", "", fmt.Errorf("error decoding token response: %w", err)
	}

	return token.AccessToken, token.RefreshToken, nil
}
