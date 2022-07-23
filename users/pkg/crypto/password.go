package crypto

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(string) ([]byte, error)
	ComparePassword(password, hash string) bool
}

type hasher struct {
	cost int
}

func NewPasswordHasher(cost int) PasswordHasher {
	if cost == 0 {
		panic("invalid cost value")
	}
	return &hasher{
		cost: cost,
	}
}

func (p *hasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), p.cost)
}

func (p *hasher) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
