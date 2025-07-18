package server

import (
	// "crypto"
	"bytes"
	"errors"

	"golang.org/x/crypto/sha3"
)

const (
	TEST_USER    string = "nate@test.ch"
	TEST_PW      string = "pw"
	TEST_TOKEN   string = "1234"
	TEST_REFRESH string = "32478"
	TEST_NAME    string = "nate"
)

func (s *HandlerSession) Authenticate(email, pw string) *userReq {

	u, err := s.store.ReadUserByEmail(email)
	if err != nil {
		return emptyUser()
	}

	pwHash, err := HashPassword(pw)
	if err != nil {
		return emptyUser()
	}

	if bytes.Equal(pwHash, u.Pw) {
		userReq := userFromModel(u)
		return userReq
	}

	return emptyUser()
}

func HashPassword(pw string) ([]byte, error) {
	sh := sha3.New256()
	_, errSh := sh.Write([]byte(pw))

	if errSh != nil {
		return nil, errors.New("could not hash password")
	}

	return sh.Sum(nil), nil
}
