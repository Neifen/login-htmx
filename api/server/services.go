package server

import (
	// "crypto"
	"bytes"

	"github.com/neifen/htmx-login/api/crypto"
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

	pwHash, err := crypto.HashPassword(pw)
	if err != nil {
		return emptyUser()
	}

	if bytes.Equal(pwHash, u.Pw) {
		userReq := userFromModel(u)
		return userReq
	}

	return emptyUser()
}