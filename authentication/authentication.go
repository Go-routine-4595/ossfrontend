package authentication

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

const (
	pkey = `
-----BEGIN CERTIFICATE-----
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
VkXcn7d3Sjgh2I7TKbisU2FvRQ==
-----END CERTIFICATE-----`

	key    = "62erSGDG35wWn55KkE7QPwMNMFs3n/BkdFXNKt7WrHY="
	user2  = "otodo"
	role2  = "partner"
	user_  = "isp1"
	role_  = "admin"
	tenant = "ISPNAME"
	kid    = "key1"
)

// UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"user"`
	Role     string `json:"role"`
}

type Storer interface {
	KeySelector(t string) []byte
	UserSelector(t string) string
	RoleSelector(t string) string
}

type Authentication struct {
	store Storer
}

func NewAuthentication(s interface{}) *Authentication {
	return &Authentication{
		store: s.(Storer),
	}
}

func (a *Authentication) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			ok   bool
			iss  string
			pKey []byte
			user string
			role string
		)
		rawAccessToken := c.GetHeader("Authorization")
		if rawAccessToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//fmt.Println("Token: ", stoken)
		//we should check the token's validity/authentication

		accessToken, _ := stripBearerPrefixFromTokenString(rawAccessToken)
		parser := jwt.NewParser()
		unverifyToken, _, err := parser.ParseUnverified(accessToken, &UserClaims{})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		iss, err = unverifyToken.Claims.GetIssuer()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		pKey = a.store.KeySelector(iss)
		user = a.store.UserSelector(iss)
		role = a.store.RoleSelector(iss)
		ok, err = a.TenantJWTValidator(rawAccessToken, pKey, user, role)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("tenant", iss)
		c.Next()
	}
}

func (a *Authentication) TenantJWTValidator(rawAccessToken string, pKey []byte, user string, role string) (bool, error) {

	accessToken, _ := stripBearerPrefixFromTokenString(rawAccessToken)
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM(pKey)
		})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}

	if strings.ToLower(claims.Username) == user && claims.Role == role {
		return true, nil
	}

	return false, errors.New("invalid user/role")
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}
