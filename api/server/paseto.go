package server

import (
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func NewToken(uid, name string) (string, *time.Time, error) {
	// Access token are signed (asymetric / public enrypt)
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	exp := time.Now().Add(2 * time.Hour)
	token.SetExpiration(exp)

	token.SetString("user-id", uid)
	token.SetString("user-name", name)

	priv := os.Getenv("TOKEN_PRIV_KEY")

	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(priv)
	if err != nil {
		return "", nil, fmt.Errorf("could not create private key")
	}

	signed := token.V4Sign(secretKey, nil)

	return signed, &exp, nil
}

func CheckToken(signed string) (*paseto.Token, error) {
	pub := os.Getenv("TOKEN_PUB_KEY")
	public, err := paseto.NewV4AsymmetricPublicKeyFromHex(pub)
	if err != nil {
		return nil, fmt.Errorf("could not create public key")
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(public, signed, nil)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func NewRefreshToken(uid string) (string, *time.Time, error) {
	// Access token are encrypted (symetric / private enrypt)

	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	exp := time.Now().Add(7 * 24 * time.Hour)
	token.SetExpiration(exp)

	token.SetString("user-id", uid)

	priv := os.Getenv("TOKEN_LOCAL_KEY")
	secretKey, err := paseto.V4SymmetricKeyFromHex(priv)
	if err != nil {
		return "", nil, err
	}

	encrypted := token.V4Encrypt(secretKey, nil)

	return encrypted, &exp, nil
}

func CheckRefreshToken(encrypted string) (*paseto.Token, error) {
	priv := os.Getenv("TOKEN_LOCAL_KEY")
	secretKey, err := paseto.V4SymmetricKeyFromHex(priv)
	if err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(secretKey, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return token, nil
}