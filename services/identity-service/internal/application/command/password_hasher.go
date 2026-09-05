package command

type PasswordHasher interface {
	Hash(password string) (string, error)

	Compare(password string, hash string) error
}
