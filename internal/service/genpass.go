package service

import (
	"crypto/rand"
	"math/big"
)

//go:generate mockery --name GenPass --filename genpass.go
type GenPass interface {
	GeneratePassword(length int) (string, error)
}

type genPassService struct {
}

func NewGenPassService() GenPass {
	return genPassService{}
}

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (g genPassService) GeneratePassword(length int) (string, error) {
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}

	return string(password), nil
}
