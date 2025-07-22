package crypto

import (
	"golang.org/x/crypto/sha3"
)

func HashPassword(pw string) ([]byte, error) {
	sh := sha3.New256()
	sh.Write([]byte(pw))

	return sh.Sum(nil), nil
}
