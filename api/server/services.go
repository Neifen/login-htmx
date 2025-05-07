package server

func (s *HandlerSession) isLoggedIn() bool {
	if s.user.crypt == nil {
		return false
	}

	return s.user.crypt.token == TEST_TOKEN
}

const (
	TEST_USER string = "nate@test.ch"
	TEST_PW string = "pw"
	TEST_TOKEN string = "1234"
	TEST_REFRESH string = "32478"
	TEST_NAME string = "nate"
)

func (s *HandlerSession) Authenticate(user, pw string) bool {

	if user == TEST_USER && pw == TEST_PW {
		s.user.crypt = NewCrypt(TEST_TOKEN)
		s.user.userName = TEST_NAME
		return true
	}

	return false
}