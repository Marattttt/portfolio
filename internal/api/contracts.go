package api

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/auth"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
)

func ToGuest(r GuestRequest) models.Guest {
	secret, hash := auth.HashSecret(r.Secret)
	return models.Guest{
		Name:   r.Name,
		Secret: secret, Salt: hash,
	}
}

func ToGuestResponse(g models.Guest) GuestResponse {
	return GuestResponse{
		Id:   int(g.ID),
		Name: g.Name,
	}
}
