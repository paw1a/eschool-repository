package repository

import (
	"context"
	"github.com/guregu/null"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

var users = []domain.User{
	domain.User{
		ID:        domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		Name:      "Pavel",
		Surname:   "Shpakovsliy",
		Phone:     null.StringFrom("+79992233444"),
		City:      null.String{},
		AvatarUrl: null.String{},
		Email:     "paw1a@yandex.ru",
		Password:  "123",
	},
	domain.User{
		ID:        domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cb"),
		Name:      "Timur",
		Surname:   "Musin",
		Phone:     null.String{},
		City:      null.StringFrom("Moscow"),
		AvatarUrl: null.String{},
		Email:     "hanoys@mail.ru",
		Password:  "qwerty",
	},
	domain.User{
		ID:        domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc"),
		Name:      "Emir",
		Surname:   "Shimshir",
		Phone:     null.StringFrom("+79992233555"),
		City:      null.String{},
		AvatarUrl: null.String{},
		Email:     "emir@gmail.com",
		Password:  "12345",
	},
}

var createdUser = domain.User{
	ID:        domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cd"),
	Name:      "createdName",
	Surname:   "createdSurname",
	Phone:     null.StringFrom("+77777777777"),
	City:      null.StringFrom("Test city"),
	AvatarUrl: null.String{},
	Email:     "user@mail.com",
	Password:  "password",
}

var updatedUser = domain.User{
	ID:        domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
	Name:      "Maxim",
	Surname:   "Shpakovsliy",
	Phone:     null.String{},
	City:      null.StringFrom("Sochi"),
	AvatarUrl: null.String{},
	Email:     "paw1a@yandex.ru",
	Password:  "12345678",
}

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	container, err := newPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	url, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test find all users", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := repository.NewUserRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all users: %v", err)
		}
		require.Equal(t, len(found), len(users))
		for i := range users {
			require.Equal(t, users[i], found[i])
		}
	})

	t.Run("test find user by id", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := repository.NewUserRepo(db)
		user, err := repo.FindByID(ctx, users[0].ID)
		if err != nil {
			t.Errorf("failed to find user with id: %v", err)
		}
		require.Equal(t, user, users[0])
	})

	t.Run("test create user", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := repository.NewUserRepo(db)
		user, err := repo.Create(ctx, createdUser)
		if err != nil {
			t.Errorf("failed to create user: %v", err)
		}
		require.Equal(t, user, createdUser)
	})

	t.Run("test update user", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := repository.NewUserRepo(db)
		user, err := repo.Update(ctx, updatedUser)
		if err != nil {
			t.Errorf("failed to create user: %v", err)
		}
		require.Equal(t, user, updatedUser)
	})

	t.Run("test delete user", func(t *testing.T) {
		t.Cleanup(func() {
			err = container.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		db, err := newPostgresDB(url)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		repo := repository.NewUserRepo(db)
		err = repo.Delete(ctx, users[0].ID)
		if err != nil {
			t.Errorf("failed to delete user: %v", err)
		}
	})
}
