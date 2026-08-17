package payment

import (
	"github.com/subhasundardass/retui/ent"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		client: client,
	}
}
