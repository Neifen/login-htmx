package crypto

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/pkg/errors"
)

type Token struct {
	Token      paseto.Token
	Expiration time.Time
	refresh    bool
	Encrypted  string
}

func (t *Token) UserName() (string, error) {
	return t.Token.GetString("user-name")
}

func (t *Token) UserID() (string, error) {
	return t.Token.GetString("user-id")
}

// A key that can only be verified with the private key
// Is faster than Asym creation and verification
func encryptedKey(t paseto.Token) (string, error) {
	priv := os.Getenv("TOKEN_LOCAL_KEY")

	secretKey, err := paseto.V4SymmetricKeyFromHex(priv)
	if err != nil {
		return "", errors.Wrap(err, "could not create private key")
	}

	encrypted := t.V4Encrypt(secretKey, nil)
	return encrypted, nil
}

func (t *Token) AddToCookie() *http.Cookie {
	/*
		Set-Cookie: access_token=eyJ…; HttpOnly; Secure
		Set-Cookie: refresh_token=…; Max-Age=31536000; Path=/api/auth/refresh; HttpOnly; Secure

	*/
	tCookie := new(http.Cookie)
	tCookie.Value = t.Encrypted
	tCookie.Expires = t.Expiration
	// tCookie.HttpOnly = true
	// tCookie.Secure = true
	tCookie.Name = "token"
	tCookie.Path = "/"
	if t.refresh {
		tCookie.Name = "refresh"
		tCookie.Path = "/token/refresh"
	}
	return tCookie

	// lets do without token-expires for now
	// expCookie := new(http.Cookie)
	// expCookie.Name = "token-expires"
	// expCookie.Value = t.Expiration.String()
	// c.SetCookie(expCookie)

	// return nil
}

// For Accesstoken normally use public keys (asymetric encryption) are used so that thirdparties can also verify the token.
// asymetric is slower but allows verification with public key.
func ValidTokenFromCookies(cookie *http.Cookie) (*Token, error) {
	if cookie == nil {
		return nil, fmt.Errorf("cannot validate token from cookie, cookie empty")
	}

	priv := os.Getenv("TOKEN_LOCAL_KEY")
	symKey, err := paseto.V4SymmetricKeyFromHex(priv)
	if err != nil {
		return nil, errors.Wrap(err, "could not create private key")
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(symKey, cookie.Value, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "token invalid (expiration: %s)", cookie.Expires)
	}
	exp, err := token.GetExpiration()
	if err != nil {
		return nil, errors.Wrapf(err, "token without expiration date")
	}
	return &Token{Token: *token, Expiration: exp, Encrypted: cookie.Value}, nil
}

func NewAccessToken(uid, name string) (*Token, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	addTime := 2 * time.Hour
	exp := time.Now().Add(addTime)
	token.SetExpiration(exp)

	token.SetString("user-id", uid)
	token.SetString("user-name", name)

	symKey, err := encryptedKey(token)
	if err != nil {
		return nil, errors.Wrap(err, "could not create symetric key for new refresh token")
	}

	return &Token{
		Token:      token,
		Encrypted:  symKey,
		Expiration: exp,
		refresh:    false,
	}, nil
}

func NewRefreshToken(uid, name string, stayloggedin bool) (*Token, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	addTime := 7 * 24 * time.Hour
	if stayloggedin {
		addTime = 40 * 24 * time.Hour
	}

	exp := time.Now().Add(addTime)
	token.SetExpiration(exp)

	token.SetString("user-id", uid)
	token.SetString("user-name", name)

	symKey, err := encryptedKey(token)
	if err != nil {
		return nil, errors.Wrap(err, "could not create symetric key for new refresh token")
	}

	return &Token{
		Token:      token,
		Encrypted:  symKey,
		Expiration: exp,
		refresh:    true,
	}, nil
}
