package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/paw1a/eschool-core/domain"
	"github.com/paw1a/eschool-core/errs"
	"github.com/paw1a/eschool-repository/postgres/entity"
	"github.com/pkg/errors"
)

type PostgresReviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) *PostgresReviewRepo {
	return &PostgresReviewRepo{
		db: db,
	}
}

const (
	reviewFindAllQuery           = "SELECT * FROM public.review"
	reviewFindByIDQuery          = "SELECT * FROM public.review WHERE id = $1"
	reviewFindUserReviewsQuery   = "SELECT * FROM public.review WHERE user_id = $1"
	reviewFindCourseReviewsQuery = "SELECT * FROM public.review WHERE course_id = $1"
	reviewDeleteQuery            = "DELETE FROM public.school WHERE id = $1"
)

func (r *PostgresReviewRepo) FindAll(ctx context.Context) ([]domain.Review, error) {
	var pgReviews []entity.PgReview
	if err := r.db.SelectContext(ctx, &pgReviews, reviewFindAllQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	reviews := make([]domain.Review, len(pgReviews))
	for i, review := range pgReviews {
		reviews[i] = review.ToDomain()
	}
	return reviews, nil
}

func (r *PostgresReviewRepo) FindByID(ctx context.Context, reviewID domain.ID) (domain.Review, error) {
	var pgReview entity.PgReview
	if err := r.db.GetContext(ctx, &pgReview, reviewFindByIDQuery, reviewID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Review{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Review{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	return pgReview.ToDomain(), nil
}

func (r *PostgresReviewRepo) FindUserReviews(ctx context.Context, userID domain.ID) ([]domain.Review, error) {
	var pgReviews []entity.PgReview
	if err := r.db.SelectContext(ctx, &pgReviews, reviewFindUserReviewsQuery, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	reviews := make([]domain.Review, len(pgReviews))
	for i, review := range pgReviews {
		reviews[i] = review.ToDomain()
	}
	return reviews, nil
}

func (r *PostgresReviewRepo) FindCourseReviews(ctx context.Context, courseID domain.ID) ([]domain.Review, error) {
	var pgReviews []entity.PgReview
	if err := r.db.SelectContext(ctx, &pgReviews, reviewFindCourseReviewsQuery, courseID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	reviews := make([]domain.Review, len(pgReviews))
	for i, review := range pgReviews {
		reviews[i] = review.ToDomain()
	}
	return reviews, nil
}

func (r *PostgresReviewRepo) Create(ctx context.Context, review domain.Review) (domain.Review, error) {
	var pgReview = entity.NewPgReview(review)
	queryString := entity.InsertQueryString(pgReview, "review")
	_, err := r.db.NamedExecContext(ctx, queryString, pgReview)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == PgUniqueViolationCode {
				return domain.Review{}, errors.Wrap(errs.ErrDuplicate, err.Error())
			} else {
				return domain.Review{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		} else {
			return domain.Review{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	var createdReview entity.PgReview
	err = r.db.GetContext(ctx, &createdReview, reviewFindByIDQuery, pgReview.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Review{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Review{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	return createdReview.ToDomain(), nil
}

func (r *PostgresReviewRepo) Delete(ctx context.Context, reviewID domain.ID) error {
	_, err := r.db.ExecContext(ctx, reviewDeleteQuery, reviewID)
	if err != nil {
		return errors.Wrap(errs.ErrDeleteFailed, err.Error())
	}
	return nil
}
