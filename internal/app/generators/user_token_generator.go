package generators

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
)

// Using!
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//Using!
func GenerateUserIDToken() (string, error) {
	id, err := generateRandom(4)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(configs.Cfg.SecretKey))
	h.Write(id)
	signedID := hex.EncodeToString(append(id, h.Sum(nil)...))

	return signedID, nil
}

func AuthUserIDToken(userIDToken string) (bool, error) {
	data, err := hex.DecodeString(userIDToken)
	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, []byte(configs.Cfg.SecretKey))
	h.Write(data[:4])
	sign := h.Sum(nil)

	if !hmac.Equal(sign, data[4:]) {
		return false, nil
	}

	return true, nil
}

func GetUserID(userIDToken string) (uint32, error) {
	data, err := hex.DecodeString(userIDToken)
	if err != nil {
		return 0, err
	}

	id := binary.BigEndian.Uint32(data[:4])
	return id, nil
}
