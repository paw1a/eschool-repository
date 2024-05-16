package repository

import (
	"context"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

var courseID = domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca")
var userID = domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca")

var reviews = []domain.Review{
	domain.Review{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7021ca"),
		UserID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		CourseID: domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca"),
		Text:     "review1 text",
	},
	domain.Review{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7021cb"),
		UserID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cb"),
		CourseID: domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca"),
		Text:     "review2 text",
	},
	domain.Review{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7021cc"),
		UserID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca"),
		CourseID: domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cb"),
		Text:     "review3 text",
	},
}

var createdReview = domain.Review{
	ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7021cd"),
	UserID:   domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc"),
	CourseID: domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cb"),
	Text:     "review3 text",
}

func TestReviewRepository(t *testing.T) {
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

	t.Run("test find all reviews", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all reviews: %v", err)
		}
		require.Equal(t, len(found), len(reviews))
		for i := range reviews {
			require.Equal(t, reviews[i], found[i])
		}
	})

	t.Run("test find review by id", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		review, err := repo.FindByID(ctx, reviews[0].ID)
		if err != nil {
			t.Errorf("failed to find review with id: %v", err)
		}
		require.Equal(t, review, reviews[0])
	})

	t.Run("test create review", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		review, err := repo.Create(ctx, createdReview)
		if err != nil {
			t.Errorf("failed to create review: %v", err)
		}
		require.Equal(t, review, createdReview)
	})

	t.Run("test find user reviews", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		found, err := repo.FindUserReviews(ctx, userID)
		if err != nil {
			t.Errorf("failed to find user reviews: %v", err)
		}

		require.Equal(t, len(found), 2)
		require.Equal(t, reviews[0], found[0])
		require.Equal(t, reviews[2], found[1])
	})

	t.Run("test find course reviews", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		found, err := repo.FindCourseReviews(ctx, courseID)
		if err != nil {
			t.Errorf("failed to find course reviews: %v", err)
		}

		require.Equal(t, len(found), 2)
		require.Equal(t, reviews[0], found[0])
		require.Equal(t, reviews[1], found[1])
	})

	t.Run("test delete review", func(t *testing.T) {
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

		repo := repository.NewReviewRepo(db)
		err = repo.Delete(ctx, reviews[0].ID)
		if err != nil {
			t.Errorf("failed to delete review: %v", err)
		}
	})
}
