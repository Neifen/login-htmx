package storage

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/neifen/htmx-login/api/crypto"
)

func TestNewUserModel(t *testing.T) {
	wantName := "testName"
	wantEmail := "my@email.com"
	unhashedPw := "testPw"

	hashedPw, err := crypto.HashPassword(unhashedPw)
	if err != nil {
		t.Fatalf("Could not hash pw: %s, error: %v\n", unhashedPw, err)
	}

	u := NewUserModel(wantName, wantEmail, unhashedPw)

	if u.Name != wantName {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Name should be %q but was %q`, wantName, wantEmail, unhashedPw, wantName, u.Name)
	}

	if u.Email != wantEmail {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Email should be %q but was %q`, wantName, wantEmail, unhashedPw, wantEmail, u.Email)
	}

	if u.Uid == "" {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Uid was empty`, wantName, wantEmail, unhashedPw)
	} else {
		t.Logf(`Uid successful, was %q`, u.Uid)
	}

	if !slices.Equal(u.Pw, hashedPw) {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Name should be %q but was %q`, wantName, wantEmail, unhashedPw, hashedPw, u.Pw)
	}
}

func TestNewUserModelEmpty(t *testing.T) {
	wantName := ""
	wantEmail := ""
	unhashedPw := ""

	hashedPw, err := crypto.HashPassword(unhashedPw)
	if err != nil {
		t.Fatalf("Could not hash pw: %s, error: %v\n", unhashedPw, err)
	}

	u := NewUserModel(wantName, wantEmail, unhashedPw)

	if u.Name != wantName {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Name should be %q but was %q`, wantName, wantEmail, unhashedPw, wantName, u.Name)
	}

	if u.Email != wantEmail {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Email should be %q but was %q`, wantName, wantEmail, unhashedPw, wantEmail, u.Email)
	}

	if u.Uid == "" {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Uid was empty`, wantName, wantEmail, unhashedPw)
	} else {
		t.Logf(`Uid successful, was %q`, u.Uid)
	}

	if !slices.Equal(u.Pw, hashedPw) {
		t.Errorf(`NewUserModel(%q, %q, %q). UserModel.Name should be %q but was %q`, wantName, wantEmail, unhashedPw, hashedPw, u.Pw)
	}
}

func TestNewRefreshTokenModel(t *testing.T) {
	uid := "asdf"
	token := "testToken"
	exp := time.Now()
	remember := true

	wantTokenModel := &RefreshTokenModel{UserUid: uid, Token: token, Expiration: exp, Remember: remember}

	tokenModel := NewRefreshTokenModel(uid, token, exp, remember)

	if !reflect.DeepEqual(tokenModel, wantTokenModel) {
		t.Errorf(`NewRefreshTokenModel(%q, %q, %q). UserModel.Name should be %+v but was %+v`, uid, token, exp, wantTokenModel, tokenModel)
	}
}

func TestNewRefreshTokenModelEmpty(t *testing.T) {
	uid := ""
	token := ""
	exp := time.Now()

	wantTokenModel := &RefreshTokenModel{UserUid: uid, Token: token, Expiration: exp}

	tokenModel := NewRefreshTokenModel(uid, token, exp, false)

	if !reflect.DeepEqual(tokenModel, wantTokenModel) {
		t.Errorf(`NewRefreshTokenModel(%q, %q, %q). UserModel.Name should be %+v but was %+v`, uid, token, exp, wantTokenModel, tokenModel)
	}
}
