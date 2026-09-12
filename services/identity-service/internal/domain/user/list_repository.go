package user

import "context"

type ListRepository interface {
	Repository
	List(ctx context.Context, offset int, limit int) ([]*User, error)
}
