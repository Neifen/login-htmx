package server

import (
	// "crypto"
	"bytes"
	"errors"
	"net/http"

	"golang.org/x/crypto/sha3"
)

func (s *HandlerSession) isLoggedIn(cookie *http.Cookie, err error) bool {
	if err != nil {
		return false
	}

	if cookie == nil {
		return false
	}

	token := cookie.Value
	if err := CheckToken(token); err != nil {
		return false
	}

	return true
}

const (
	TEST_USER    string = "nate@test.ch"
	TEST_PW      string = "pw"
	TEST_TOKEN   string = "1234"
	TEST_REFRESH string = "32478"
	TEST_NAME    string = "nate"
)

func (s *HandlerSession) Authenticate(email, pw string) bool {

	u, err := s.store.ReadUserByEmail(email)
	if err != nil {
		return false
	}

	pwHash, err := HashPassword(pw)
	if err != nil {
		return false
	}

	if bytes.Equal(pwHash, u.pw) {
		s.user = u.ToUserInfo()
		s.user.AddCrypt()
		return true
	}

	return false
}

func HashPassword(pw string) ([]byte, error) {
	sh := sha3.New256()
	_, errSh := sh.Write([]byte(pw))

	if errSh != nil {
		return nil, errors.New("could not hash password")
	}

	return sh.Sum(nil), nil

}
