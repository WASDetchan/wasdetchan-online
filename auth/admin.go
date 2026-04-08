package auth

import (
	"context"
	"os"

	"github.com/WASDetchan/wasdetchan-online/repository"
)

func MakeAdmin(email string, q *repository.Queries) {
	q.MakeAdmin(context.Background(), email)
}

func MakeOwnerAdmin(q *repository.Queries) {
	owner := os.Getenv("OWNER_EMAIL")
	if owner != "" {
		MakeAdmin(owner, q)
	}
}

func IsOwner(email string) bool {
	return email == os.Getenv("OWNER_EMAIL")
}
