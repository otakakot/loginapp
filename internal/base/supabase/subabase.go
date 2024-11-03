package supabase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/supabase-community/gotrue-go"
)

type Supabase struct {
	client gotrue.Client
}

func New(
	projectReference string,
	apiKey string,
) *Supabase {
	return &Supabase{
		client: gotrue.New(projectReference, apiKey),
	}
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserID string
}

func (sb *Supabase) Auth(
	ctx context.Context,
	request AuthRequest,
) (*AuthResponse, error) {
	res, err := sb.client.SignInWithEmailPassword(request.Email, request.Password)
	if err != nil {
		return nil, fmt.Errorf("sign in with email and password: %w", err)
	}

	slog.Info(fmt.Sprintf("response: %+v", res))

	return &AuthResponse{
		UserID: res.User.ID.String(),
	}, nil
}
