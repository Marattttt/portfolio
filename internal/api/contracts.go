package api

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/auth"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
)

func ToGuest(r GuestRequest) models.Guest {
	pass, hash := auth.HashSecret([]byte(r.Secret))
	return models.Guest{
		Name:   r.Name,
		Secret: string(pass),
		Salt:   string(hash),
	}
}

func ToGuestResponse(g models.Guest) GuestResponse {
	return GuestResponse{
		Id:     int(g.ID),
		Name:   g.Name,
		Secret: g.Secret,
	}
}
