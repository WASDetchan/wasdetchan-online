package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha3"
	"encoding/base64"
	"net/http"

	"github.com/WASDetchan/wasdetchan-online/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func authToken(token string, q *repository.Queries) AuthInfo {
	dtoken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return AuthInfo{}
	}
	token_hash := sha3.Sum512(dtoken)
	rtoken, err := q.GetToken(context.Background(), token_hash[:])
	if err != nil {
		return AuthInfo{}
	}

	user, err := q.GetUser(context.Background(), rtoken.UserID)
	if err != nil {
		return AuthInfo{}
	}

	c := AuthInfo{User: &user}
	c.Capabilities.ReadBytes(rtoken.Capabilities)

	return c
}

func tokenCapabilities(token string, userID pgtype.UUID, q *repository.Queries) Capabilities {
	info := authToken(token, q)
	if info.User.ID != userID {
		return Capabilities{}
	} else {
		return info.Capabilities
	}
}

func CreateToken(c *gin.Context, capabilities ...Capability) {
	user := AssertAuth(c)

	token := make([]byte, 48)
	rand.Read(token)
	hash := sha3.Sum512(token)

	caps := Capabilities{CapabilityReceiptWrite}

	q := repository.GetQueries(c)
	_, err := q.CreateToken(context.Background(), repository.CreateTokenParams{UserID: user.ID, TokenHash: hash[:], Capabilities: caps.IntoBytes()})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(token))

	c.String(http.StatusOK, encoded)
}
