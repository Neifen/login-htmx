package crypto

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-test/deep"
)

func TestInvalidPrivateKey(t *testing.T) {
	os.Setenv("TOKEN_LOCAL_KEY", "♥")

	access, err := NewAccessToken("what", "ever")
	if err == nil {
		t.Errorf(`NewAccessToken: Having ♥ as a Private key should return an error`)
	}

	expError := "encoding/hex: invalid byte: U+00E2 'â'"
	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`NewAccessToken: Having ♥ as a Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if access != nil {
		t.Errorf(`NewAccessToken: Having ♥ as a Private key should not return an access token: %v`, access)
	}

	key, err := encryptedKey(paseto.NewToken())
	if err == nil {
		t.Errorf(`SymKey: Having ♥ as a Private key should return an error`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`SymKey: Having ♥ as a Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if len(key) != 0 {
		t.Errorf(`SymKey: Having ♥ as a Private key should not return a key: %s`, key)
	}

	token, err := ValidTokenFromCookies(new(http.Cookie))
	if err == nil {
		t.Errorf(`ValidSynmTokenFromCookies: Having ♥ as a Private key should return an error for ValidSynmTokenFromCookies`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`ValidSynmTokenFromCookies: Having ♥ as a Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if token != nil {
		t.Errorf(`ValidSynmTokenFromCookies: Having ♥ as a Private key should not return a token: %v`, token)
	}

	refresh, err := NewRefreshToken("what", "ever", false)
	if err == nil {
		t.Errorf(`NewRefreshToken: Having ♥ as a Private key should return an error`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`NewRefreshToken: Having ♥ as a Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if refresh != nil {
		t.Errorf(`NewRefreshToken: Having ♥ as a Private key should not return a refresh token: %v`, refresh)
	}
}

func TestInvalidPrivateKeyEmpty(t *testing.T) {
	os.Setenv("TOKEN_LOCAL_KEY", "")

	access, err := NewAccessToken("what", "ever")
	if err == nil {
		t.Errorf(`NewAccessToken: Having an empty Private key should return an error`)
	}

	expError := "key length incorrect (0), expected 32"

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`NewAccessToken: Having an empty Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if access != nil {
		t.Errorf(`NewAccessToken: Having an empty Private key should not return an access token: %v`, access)
	}

	key, err := encryptedKey(paseto.NewToken())
	if err == nil {
		t.Errorf(`SymKey: Having an empty Private key should return an error`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`SymKey: Having an empty Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if len(key) != 0 {
		t.Errorf(`SymKey: Having an empty Private key should not return a key: %s`, key)
	}

	token, err := ValidTokenFromCookies(new(http.Cookie))
	if err == nil {
		t.Errorf(`ValidSynmTokenFromCookies: Having an empty Private key should return an error for ValidSynmTokenFromCookies`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`ValidSynmTokenFromCookies: Having an empty Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if token != nil {
		t.Errorf(`ValidSynmTokenFromCookies: Having an empty Private key should not return a token: %v`, token)
	}

	refresh, err := NewRefreshToken("what", "ever", true)
	if err == nil {
		t.Errorf(`NewRefreshToken: Having an empty Private key should return an error`)
	}

	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`NewRefreshToken: Having an empty Private key returned the following error %q expected %q`, err.Error(), expError)
	}

	if refresh != nil {
		t.Errorf(`NewRefreshToken: Having an empty Private key should not return a refresh token: %v`, refresh)
	}
}

func TestAccessAndRefreshTokenNotEqual(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	access, err := NewAccessToken(uid, name)
	if err != nil {
		t.Fatalf(`NewAccessToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	refresh, err := NewRefreshToken(uid, name, false)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	accessKey := access.Encrypted
	refreshKey := refresh.Encrypted

	if accessKey == refreshKey {
		t.Errorf(`Access token and Refresh token with same parameters (%q, %q) were equal: %q`, uid, name, refreshKey)
	}
}

func TestAccessAndRefreshTokenNotEqualEmpty(t *testing.T) {
	uid := ""
	name := ""
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	access, err := NewAccessToken(uid, name)
	if err != nil {
		t.Fatalf(`NewAccessToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	refresh, err := NewRefreshToken(uid, name, false)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	accessKey := access.Encrypted
	refreshKey := refresh.Encrypted

	if accessKey == refreshKey {
		t.Errorf(`Access token and Refresh token with same parameters (%q, %q) were equal: %q`, uid, name, refreshKey)
	}
}

func TestAccessTokenSymetricKey(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	access, err := NewAccessToken(uid, name)
	if err != nil {
		t.Fatalf(`NewAccessToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	parsedKey, err := paseto.NewParser().ParseV4Local(symKey, access.Encrypted, nil)
	if err != nil {
		t.Fatalf(`Parsing NewAccessToken(%q, %q).SymKey() failed with the following error %v`, uid, name, err)
	}

	diffs := deep.Equal(*parsedKey, access.Token)
	if len(diffs) != 0 {
		t.Errorf(`Parsing NewAccessToken(%q, %q).SymKey() returned this token: %v, but expected this %v \n Diffs: %v`, uid, name, *parsedKey, access.Token, diffs)
	}
}

func TestAddAccessTokenToCookie(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	access, err := NewAccessToken(uid, name)
	if err != nil {
		t.Fatalf(`NewAccessToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	cookie := access.AddToCookie()

	expectePath := "/"
	expectExp := time.Now().Add(2 * time.Hour)

	if !timeAlmostEqual(expectExp, cookie.Expires) {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Expires = %v, expected +/- 1s around %v`, uid, name, cookie.Expires, expectExp)
	}

	if !cookie.Secure {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Secure = %v, expected true`, uid, name, cookie.Secure)
	}

	if !cookie.HttpOnly {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.HttpOnly = %v, expected true`, uid, name, cookie.HttpOnly)
	}

	if expectePath != cookie.Path {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Path = %v, expected an empty path`, uid, name, cookie.Path)
	}

	token, err := paseto.Parser.ParseV4Local(paseto.NewParser(), symKey, cookie.Value, nil)
	if err != nil {
		t.Fatalf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value failed with the following error %v`, uid, name, err)
	}

	diffs := deep.Equal(*token, access.Token)
	if len(diffs) != 0 {
		t.Errorf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value was different from original token, with the following diffs: %v`, uid, name, diffs)
	}

	// from
	tokenFromCookie, err := ValidTokenFromCookies(cookie)
	if err != nil {
		t.Fatalf(`ValidTokenFromCookies(%v) failed with the following error %v`, cookie, err)
	}

	if !timeAlmostEqual(tokenFromCookie.Expiration, access.Expiration) {
		t.Errorf(`ValidTokenFromCookies().Expiration = %v, expected: %v`, tokenFromCookie.Expiration, access.Expiration)
	}

	diffs = deep.Equal(tokenFromCookie, access)
	if len(diffs) != 1 { //Expiration should be different
		t.Errorf(`ValidSynmTokenFromCookies(%v) was different from original token, with the following diffs: %v`, tokenFromCookie, diffs)
	}
}

func TestAddRefreshTokenToCookie(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	refresh, err := NewRefreshToken(uid, name, false)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	cookie := refresh.AddToCookie()

	expectePath := "/token"
	expectExp := time.Now().Add(7 * 24 * time.Hour)

	if !timeAlmostEqual(expectExp, cookie.Expires) {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Expires = %v, expected +/- 1s around %v`, uid, name, cookie.Expires, expectExp)
	}

	if !cookie.Secure {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Secure = %v, expected true`, uid, name, cookie.Secure)
	}

	if !cookie.HttpOnly {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.HttpOnly = %v, expected true`, uid, name, cookie.HttpOnly)
	}

	if expectePath != cookie.Path {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Path = %v, expected %v`, uid, name, cookie.Path, expectePath)
	}

	token, err := paseto.Parser.ParseV4Local(paseto.NewParser(), symKey, cookie.Value, nil)
	if err != nil {
		t.Fatalf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value failed with the following error %v`, uid, name, err)
	}

	diffs := deep.Equal(*token, refresh.Token)
	if len(diffs) != 0 {
		t.Errorf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value was different from original token, with the following diffs: %v`, uid, name, diffs)
	}

	// from
	tokenFromCookie, err := ValidTokenFromCookies(cookie)
	if err != nil {
		t.Fatalf(`ValidSynmTokenFromCookies(%v) failed with the following error %v`, cookie, err)
	}

	if !timeAlmostEqual(tokenFromCookie.Expiration, refresh.Expiration) {
		t.Errorf(`ValidTokenFromCookies().Expiration = %v, expected: %v`, tokenFromCookie.Expiration, refresh.Expiration)
	}

	diffs = deep.Equal(tokenFromCookie, refresh)
	if len(diffs) != 1 { //Expiration should be different
		t.Errorf(`ValidSynmTokenFromCookies(%v) was different from original token, with the following diffs: %v`, tokenFromCookie, diffs)
	}
}

func TestAddRefreshTokenStayLoggedInToAndFromCookie(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	refresh, err := NewRefreshToken(uid, name, true)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}
	cookie := refresh.AddToCookie()

	expectePath := "/token"
	expectExp := time.Now().Add(40 * 24 * time.Hour)

	if !timeAlmostEqual(expectExp, cookie.Expires) {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Expires = %v, expected +/- 1s around %v`, uid, name, cookie.Expires, expectExp)
	}

	if !cookie.Secure {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Secure = %v, expected true`, uid, name, cookie.Secure)
	}

	if !cookie.HttpOnly {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.HttpOnly = %v, expected true`, uid, name, cookie.HttpOnly)
	}

	if expectePath != cookie.Path {
		t.Errorf(`NewAccessToken(%q, %q).AddToCookie.Path = %v, expected %q`, uid, name, cookie.Path, expectePath)
	}

	token, err := paseto.Parser.ParseV4Local(paseto.NewParser(), symKey, cookie.Value, nil)
	if err != nil {
		t.Fatalf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value failed with the following error %v`, uid, name, err)
	}

	diffs := deep.Equal(*token, refresh.Token)
	if len(diffs) != 0 {
		t.Errorf(`Parsing NewAccessToken(%q, %q).AddToCookie().Value was different from original token, with the following diffs: %v`, uid, name, diffs)
	}

	// from
	tokenFromCookie, err := ValidTokenFromCookies(cookie)
	if err != nil {
		t.Fatalf(`ValidSynmTokenFromCookies(%v) failed with the following error %v`, cookie, err)
	}

	if !timeAlmostEqual(tokenFromCookie.Expiration, refresh.Expiration) {
		t.Errorf(`ValidTokenFromCookies().Expiration = %v, expected: %v`, tokenFromCookie.Expiration, refresh.Expiration)
	}

	diffs = deep.Equal(tokenFromCookie, refresh)
	if len(diffs) != 1 { //Expiration should be different
		t.Errorf(`ValidSynmTokenFromCookies(%v) was different from original token, with the following diffs: %v`, tokenFromCookie, diffs)
	}
}

func TestInvalidSynmTokenFromCookiesEmpty(t *testing.T) {
	token, err := ValidTokenFromCookies(nil)
	if err == nil {
		t.Errorf(`ValidSynmTokenFromCookies(nil) should return error`)
	}

	expError := "cannot validate token from cookie, cookie empty"
	if err.Error() != expError {
		t.Errorf(`ValidSynmTokenFromCookies(nil) error said %q instead of %q`, err.Error(), expError)
	}

	if token != nil {
		t.Errorf(`ValidSynmTokenFromCookies(nil) = %v should be nil`, token)
	}
}

func TestValidSynmTokenFromCookiesInvalid(t *testing.T) {
	sym := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", sym.ExportHex())

	invalid := paseto.NewToken()
	invalid.SetExpiration(time.Now())

	encryptedKey, err := encryptedKey(invalid)
	if err != nil {
		t.Fatalf(`encryptedKey(nil) returned the following error: %v`, err)
	}

	token := &Token{Token: invalid, Expiration: time.Now(), Encrypted: encryptedKey}
	cookie := token.AddToCookie()

	newToken, err := ValidTokenFromCookies(cookie)
	if err == nil {
		t.Errorf(`ValidSynmTokenFromCookies(nil) should return error`)
	}

	expError := "this token has expired"
	if !strings.Contains(err.Error(), expError) {
		t.Errorf(`ValidSynmTokenFromCookies(nil) error said %q instead of %q`, err.Error(), expError)
	}

	if newToken != nil {
		t.Errorf(`ValidSynmTokenFromCookies(nil) = %v should be nil`, token)
	}
}

func TestRefreshTokenSymetricKey(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"
	symKey := paseto.NewV4SymmetricKey()
	os.Setenv("TOKEN_LOCAL_KEY", symKey.ExportHex())

	refresh, err := NewRefreshToken(uid, name, true)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	parsedKey, err := paseto.NewParser().ParseV4Local(symKey, refresh.Encrypted, nil)
	if err != nil {
		t.Fatalf(`Parsing NewRefreshToken(%q, %q).Encrypted failed with the following error %v`, uid, name, err)
	}

	diffs := deep.Equal(*parsedKey, refresh.Token)
	if len(diffs) != 0 {
		t.Errorf(`Parsing NewRefreshToken(%q, %q).Token returned this token: %v, but expected this %v \n Diffs: %v`, uid, name, *parsedKey, refresh.Token, diffs)
	}
}

func TestNewAccessToken(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"

	expectedExp := time.Now().Add(2 * time.Hour)
	token, err := NewAccessToken(uid, name)
	if err != nil {
		t.Fatalf(`NewAccessToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	tuid, err := token.UserID()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q).UserID() failed with the following error %v`, uid, name, err)
	}

	if uid != tuid {
		t.Errorf(`NewAccessToken(%q, %q).UserID() = %v, expected %v`, uid, name, tuid, uid)
	}

	tname, err := token.UserName()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q).UserName() failed with the following error %v`, uid, name, err)
	}

	if name != tname {
		t.Errorf(`NewAccessToken(%q, %q).UserName() = %v, expected %v`, uid, name, tname, name)
	}

	if !timeAlmostEqual(expectedExp, token.Expiration) {
		t.Errorf(`NewAccessToken(%q, %q).Expiration = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	exp, _ := token.Token.GetExpiration()
	if !timeAlmostEqual(expectedExp, exp) {
		t.Errorf(`NewAccessToken(%q, %q).Token.GetExpiration() = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	if token.refresh {
		t.Errorf(`NewAccessToken(%q, %q).refresh = %v, %v`, uid, name, token.refresh, false)
	}
}

func TestNewRefreshToken(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"

	expectedExp := time.Now().Add(7 * 24 * time.Hour)
	token, err := NewRefreshToken(uid, name, false)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	tuid, err := token.UserID()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q, false).UserID() failed with the following error %v`, uid, name, err)
	}

	if uid != tuid {
		t.Errorf(`NewRefreshToken(%q, %q, false).UserID() = %v, expected %v`, uid, name, tuid, uid)
	}

	tname, err := token.UserName()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q, false).UserName() failed with the following error %v`, uid, name, err)
	}

	if name != tname {
		t.Errorf(`NewRefreshToken(%q, %q, false).UserName() = %v, expected %v`, uid, name, tname, name)
	}

	if !timeAlmostEqual(expectedExp, token.Expiration) {
		t.Errorf(`NewRefreshToken(%q, %q, false).Expiration = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	exp, _ := token.Token.GetExpiration()
	if !timeAlmostEqual(expectedExp, exp) {
		t.Errorf(`NewRefreshToken(%q, %q, false).Token.GetExpiration() = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	if !token.refresh {
		t.Errorf(`NewRefreshToken(%q, %q, false).refresh = %v, %v`, uid, name, token.refresh, true)
	}
}

func TestNewRefreshTokenStayLoggedIn(t *testing.T) {
	uid := "this Is an uid"
	name := "my Name"

	expectedExp := time.Now().Add(40 * 24 * time.Hour)
	token, err := NewRefreshToken(uid, name, true)
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q) failed with the following error %v`, uid, name, err)
	}

	tuid, err := token.UserID()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q, true).UserID() failed with the following error %v`, uid, name, err)
	}

	if uid != tuid {
		t.Errorf(`NewRefreshToken(%q, %q, true).UserID() = %v, expected %v`, uid, name, tuid, uid)
	}

	tname, err := token.UserName()
	if err != nil {
		t.Fatalf(`NewRefreshToken(%q, %q, true).UserName() failed with the following error %v`, uid, name, err)
	}

	if name != tname {
		t.Errorf(`NewRefreshToken(%q, %q, true).UserName() = %v, expected %v`, uid, name, tname, name)
	}

	if !timeAlmostEqual(expectedExp, token.Expiration) {
		t.Errorf(`NewRefreshToken(%q, %q, true).Expiration = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	exp, _ := token.Token.GetExpiration()
	if !timeAlmostEqual(expectedExp, exp) {
		t.Errorf(`NewRefreshToken(%q, %q, true).Token.GetExpiration() = %v, expected +/- 1s around %v`, uid, name, token.Expiration, expectedExp)
	}

	if !token.refresh {
		t.Errorf(`NewRefreshToken(%q, %q, true).refresh = %v, %v`, uid, name, token.refresh, true)
	}
}

func timeAlmostEqual(expected, actual time.Time) bool {
	return actual.After(expected.Add(-time.Second)) && actual.Before(expected.Add(time.Second))
}
