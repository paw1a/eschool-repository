package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/paw1a/eschool-core/domain"
)

type PgUser struct {
	ID        uuid.UUID   `db:"id"`
	Name      string      `db:"name"`
	Surname   string      `db:"surname"`
	Phone     null.String `db:"phone"`
	City      null.String `db:"city"`
	AvatarUrl null.String `db:"avatar_url"`
	Email     string      `db:"email"`
	Password  string      `db:"password"`
}

func (u *PgUser) ToDomain() domain.User {
	return domain.User{
		ID:        domain.ID(u.ID.String()),
		Name:      u.Name,
		Surname:   u.Surname,
		Phone:     u.Phone,
		City:      u.City,
		AvatarUrl: u.AvatarUrl,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func NewPgUser(user domain.User) PgUser {
	id, _ := uuid.Parse(user.ID.String())
	return PgUser{
		ID:        id,
		Name:      user.Name,
		Surname:   user.Surname,
		Phone:     user.Phone,
		City:      user.City,
		AvatarUrl: user.AvatarUrl,
		Email:     user.Email,
		Password:  user.Password,
	}
}
