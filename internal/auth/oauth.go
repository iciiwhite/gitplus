package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/iciwhite/gitplus/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthService struct {
	config *config.Config
	token  *oauth2.Token
	client *http.Client
}

func NewOAuthService(cfg *config.Config) *OAuthService {
	return &OAuthService{
		config: cfg,
	}
}

func (s *OAuthService) IsAuthenticated() bool {
	data, err := os.ReadFile(s.tokenPath())
	if err != nil {
		return false
	}
	var tok oauth2.Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return false
	}
	s.token = &tok
	s.client = s.oauthConfig().Client(context.Background(), s.token)
	return true
}

func (s *OAuthService) StartAuthFlow() error {
	conf := s.oauthConfig()
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Open this URL in your browser:\n%s\n", url)

	codeChan := make(chan string)
	srv := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Authorization successful. You may close this window.")
		codeChan <- code
	})}
	go func() {
		_ = srv.ListenAndServe()
	}()
	defer srv.Close()

	select {
	case code := <-codeChan:
		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			return err
		}
		s.token = tok
		s.client = conf.Client(context.Background(), tok)
		return s.saveToken()
	case <-time.After(2 * time.Minute):
		return fmt.Errorf("authorization timeout")
	}
}

func (s *OAuthService) GetClient() *http.Client {
	return s.client
}

func (s *OAuthService) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.GitHubClientID,
		ClientSecret: s.config.GitHubClientSecret,
		Scopes:       []string{"repo", "user", "read:org"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/callback",
	}
}

func (s *OAuthService) saveToken() error {
	data, err := json.Marshal(s.token)
	if err != nil {
		return err
	}
	return os.WriteFile(s.tokenPath(), data, 0600)
}

func (s *OAuthService) tokenPath() string {
	return os.Getenv("HOME") + "/.gitplus_token.json"
}