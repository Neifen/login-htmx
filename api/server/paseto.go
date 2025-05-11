package server

import (
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func NewToken(u *UserInfo) (string, *time.Time, error) {
	// Access token are signed (asymetric / public enrypt)
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	exp := time.Now().Add(2 * time.Hour)
	token.SetExpiration(exp)

	token.SetString("user-id", u.userName)

	priv := os.Getenv("TOKEN_PRIV_KEY")

	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(priv)
	if err != nil {
		return "", nil, fmt.Errorf("could not create private key")
	}

	signed := token.V4Sign(secretKey, nil)

	return signed, &exp, nil
}

func CheckToken(signed string) error {
	pub := os.Getenv("TOKEN_PUB_KEY")
	public, err := paseto.NewV4AsymmetricPublicKeyFromHex(pub)
	if err != nil {
		return fmt.Errorf("could not create public key")
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(public, signed, nil)
	fmt.Println(token)
	if err != nil {
		return err
	}
	return nil
}

func NewRefreshToken(u *UserInfo) string {
	// Access token are encrypted (symetric / private enrypt)

	//todo do this
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.SetString("user-id", u.userName)

	secretKey := paseto.NewV4SymmetricKey() // do not share
	encrypted := token.V4Encrypt(secretKey, nil)

	return encrypted
}
