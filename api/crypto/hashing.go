package crypto

import (
	"golang.org/x/crypto/sha3"
)

func HashPassword(pw string) ([]byte, error) {
	sh := sha3.New256()
	_, errSh := sh.Write([]byte(pw))

	if errSh != nil {
		return nil, errSh
	}

	return sh.Sum(nil), nil
}
