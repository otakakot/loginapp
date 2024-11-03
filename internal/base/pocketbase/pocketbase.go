package pocketbase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Pocketbase struct {
	baseURL string
}

func New(baseURL string) *Pocketbase {
	return &Pocketbase{
		baseURL: baseURL,
	}
}

type AuthRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	Record struct {
		ID              string `json:"id"`
		CollectionID    string `json:"collectionId"`
		CollectionName  string `json:"collectionName"`
		Username        string `json:"username"`
		Verified        bool   `json:"verified"`
		EmailVisibility bool   `json:"emailVisibility"`
		Email           string `json:"email"`
		Created         string `json:"created"`
		Updated         string `json:"updated"`
		Name            string `json:"name"`
		Avatar          string `json:"avatar"`
	} `json:"record"`
}

// Auth.
// ref: https://pocketbase.io/docs/api-records/#auth-with-password
func (pb *Pocketbase) Auth(
	ctx context.Context,
	request AuthRequest,
) (*AuthResponse, error) {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		pb.baseURL+"/api/collections/users/auth-with-password",
		buf,
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
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response AuthResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info(fmt.Sprintf("response: %+v", response))

	return &response, nil
}
