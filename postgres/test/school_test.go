package repository

import (
	"context"
	"github.com/guregu/null"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

var schools = []domain.School{
	domain.School{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cc"),
		OwnerID:     domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		Name:        "school1",
		Description: "desc1",
	},
	domain.School{
		ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cd"),
		OwnerID:     domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cb"),
		Name:        "school2",
		Description: "desc2",
	},
}

var newTeacherID = domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc")

var teachers = []domain.User{
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
}

var createdSchool = domain.School{
	ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034ce"),
	OwnerID:     domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
	Name:        "school3",
	Description: "desc3",
}

var updatedSchool = domain.School{
	ID:          domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cc"),
	OwnerID:     domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
	Name:        "updated school",
	Description: "updated desc",
}

func TestSchoolRepository(t *testing.T) {
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

	t.Run("test find all schools", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all schools: %v", err)
		}
		require.Equal(t, len(found), len(schools))
		for i := range schools {
			require.Equal(t, schools[i], found[i])
		}
	})

	t.Run("test find school by id", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		school, err := repo.FindByID(ctx, schools[0].ID)
		if err != nil {
			t.Errorf("failed to find school with id: %v", err)
		}
		require.Equal(t, school, schools[0])
	})

	t.Run("test find school teachers", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		found, err := repo.FindSchoolTeachers(ctx, schools[0].ID)
		if err != nil {
			t.Errorf("failed to find school teachers: %v", err)
		}
		require.Equal(t, len(found), len(teachers))
		for i := range teachers {
			require.Equal(t, teachers[i], found[i])
		}
	})

	t.Run("test check is school teacher", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		ok, err := repo.IsSchoolTeacher(ctx, schools[0].ID, teachers[0].ID)
		if err != nil {
			t.Errorf("failed to find school teacher: %v", err)
		}
		require.Equal(t, ok, true)

		ok, err = repo.IsSchoolTeacher(ctx, schools[0].ID, newTeacherID)
		if err != nil {
			t.Errorf("failed to find school teacher: %v", err)
		}
		require.Equal(t, ok, false)
	})

	t.Run("test add school teacher", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		err = repo.AddSchoolTeacher(ctx, schools[0].ID, newTeacherID)
		if err != nil {
			t.Errorf("failed to add school teacher: %v", err)
		}
	})

	t.Run("test create school", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		school, err := repo.Create(ctx, createdSchool)
		if err != nil {
			t.Errorf("failed to create school: %v", err)
		}
		require.Equal(t, school, createdSchool)
	})

	t.Run("test update school", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		school, err := repo.Update(ctx, updatedSchool)
		if err != nil {
			t.Errorf("failed to create school: %v", err)
		}
		require.Equal(t, school, updatedSchool)
	})

	t.Run("test delete school", func(t *testing.T) {
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

		repo := repository.NewSchoolRepo(db)
		err = repo.Delete(ctx, schools[0].ID)
		if err != nil {
			t.Errorf("failed to delete school: %v", err)
		}
	})
}
