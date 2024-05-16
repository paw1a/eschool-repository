package repository

import (
	"context"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var certificates = []domain.Certificate{
	domain.Certificate{
		ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027ca"),
		CourseID:  domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca"),
		UserID:    domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		Name:      "course 1 cert",
		Score:     120,
		CreatedAt: time.Now(),
		Grade:     domain.GoldCertificate,
	},
	domain.Certificate{
		ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027cb"),
		CourseID:  domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cb"),
		UserID:    domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		Name:      "course 2 cert",
		Score:     50,
		CreatedAt: time.Now(),
		Grade:     domain.BronzeCertificate,
	},
}

var createdCertificate = domain.Certificate{
	ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027cc"),
	CourseID:  domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7026cc"),
	UserID:    domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
	Name:      "course 3 cert",
	Score:     100,
	CreatedAt: time.Now(),
	Grade:     domain.SilverCertificate,
}

func TestCertificateRepository(t *testing.T) {
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

	t.Run("test find all certificates", func(t *testing.T) {
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

		repo := repository.NewCertificateRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all certificates: %v", err)
		}
		require.Equal(t, len(found), len(certificates))
	})

	t.Run("test find certificate by id", func(t *testing.T) {
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

		repo := repository.NewCertificateRepo(db)
		certificate, err := repo.FindByID(ctx, certificates[0].ID)
		if err != nil {
			t.Errorf("failed to find certificate with id: %v", err)
		}
		certificate.CreatedAt = certificates[0].CreatedAt
		require.Equal(t, certificate, certificates[0])
	})

	t.Run("test find user certificates", func(t *testing.T) {
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

		repo := repository.NewCertificateRepo(db)
		found, err := repo.FindUserCertificates(ctx, certificates[0].UserID)
		if err != nil {
			t.Errorf("failed to find user certificates: %v", err)
		}
		require.Equal(t, len(found), 2)
	})

	t.Run("test create certificate", func(t *testing.T) {
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

		repo := repository.NewCertificateRepo(db)
		certificate, err := repo.Create(ctx, createdCertificate)
		if err != nil {
			t.Errorf("failed to create certificate: %v", err)
		}
		certificate.CreatedAt = createdCertificate.CreatedAt
		require.Equal(t, certificate, createdCertificate)
	})
}
