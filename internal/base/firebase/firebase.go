package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const endpoint = "https://identitytoolkit.googleapis.com"

type Option func(*Firebase)

func Endpoint(endpoint string) Option {
	return func(f *Firebase) {
		if endpoint == "" {
			return
		}

		f.endpoint = endpoint
	}
}

type Firebase struct {
	endpoint string
	apikey   string
}

func New(apikey string, options ...Option) *Firebase {
	fb := &Firebase{
		apikey:   apikey,
		endpoint: endpoint,
	}

	for _, option := range options {
		option(fb)
	}

	return fb
}

type AuthRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSerureToken bool   `json:"returnSecureToken"`
}

type AuthResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	Registered   bool   `json:"registered"`
}

// Auth.
// ref: https://firebase.google.com/docs/reference/rest/auth#section-sign-in-email-password
func (fb *Firebase) Auth(
	ctx context.Context,
	request AuthRequest,
) (*AuthResponse, error) {
	url := fmt.Sprintf("%s/v1/accounts:signInWithPassword?key=%s", fb.endpoint, fb.apikey)

	body := new(bytes.Buffer)

	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	var response AuthResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info(fmt.Sprintf("response: %+v", response))

	return &response, nil
}
