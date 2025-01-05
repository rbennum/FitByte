package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Service interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewJWTService() *jwtService {
	return &jwtService{}
}

const ENV_PATH = ".env"

var _ = godotenv.Load(ENV_PATH)
var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

func (s *jwtService) GenerateToken(userID string) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func ParseToken(tokeString string) (id string, err error) {
	token, err := jwt.Parse(tokeString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("invalid token")
	}

	userID := claims["user_id"].(string)

	return userID, nil
}
