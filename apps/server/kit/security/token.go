package security

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrParsingKey   = errors.New("error parsing RSA key")
)

// JWTClaims contains the claims embedded in every access token.
// There is only one token type (access); no refresh tokens are issued.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func ParseRSAPrivateKeyFromBase64(base64Key string) (*rsa.PrivateKey, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, ErrParsingKey
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, ErrParsingKey
	}
	return privateKey, nil
}

func ParseRSAPublicKeyFromBase64(base64Key string) (*rsa.PublicKey, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, ErrParsingKey
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, ErrParsingKey
	}
	return publicKey, nil
}

// GenerateAccessToken issues a single RS256-signed access token valid for expiration duration.
// The token is returned as a raw JWT string to be sent in the JSON response body.
func GenerateAccessToken(userID string, role string, privateKey *rsa.PrivateKey, expiration time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func ValidateAccessToken(tokenString string, publicKey *rsa.PublicKey) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
