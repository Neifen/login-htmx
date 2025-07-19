package crypto

import (
	"slices"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestHashPassword(t *testing.T) {
	inputPw := "anything"

	pw, err := HashPassword(inputPw)
	
	if(err != nil){
		t.Fatalf(`HashPassword(%q) failed with error %v`, inputPw, err)
	}

	sh := sha3.New256()
	sh.Write([]byte(inputPw))

	expected := sh.Sum(nil)

	if(!slices.Equal(expected, pw)) {
		t.Errorf(`HashPassword(%q) = %q, expected %q`, inputPw, pw, expected)
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	inputPw := ""

	pw, err := HashPassword(inputPw)
	
	if(err != nil){
		t.Fatalf(`HashPassword(%q) failed with error %v`, inputPw, err)
	}

	sh := sha3.New256()
	sh.Write([]byte(inputPw))

	expected := sh.Sum(nil)

	if(!slices.Equal(expected, pw)) {
		t.Errorf(`HashPassword(%q) = %q, expected %q`, inputPw, pw, expected)
	}
}

func TestHashPasswordEmojis(t *testing.T) {
	inputPw := "♥◘♣○╗ß"

	pw, err := HashPassword(inputPw)
	
	if(err != nil){
		t.Fatalf(`HashPassword(%q) failed with error %v`, inputPw, err)
	}

	sh := sha3.New256()
	sh.Write([]byte(inputPw))

	expected := sh.Sum(nil)

	if(!slices.Equal(expected, pw)) {
		t.Errorf(`HashPassword(%q) = %q, expected %q`, inputPw, pw, expected)
	}
}
