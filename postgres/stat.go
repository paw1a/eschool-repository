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

type PostgresStatRepo struct {
	db *sqlx.DB
}

func NewStatRepo(db *sqlx.DB) *PostgresStatRepo {
	return &PostgresStatRepo{
		db: db,
	}
}

const (
	statFindByUserLessonQuery = "SELECT * FROM public.lesson_stat WHERE user_id = $1 AND lesson_id = $2"
	statFindByUserTestQuery   = "SELECT * FROM public.test_stat WHERE user_id = $1 AND test_id = $2"
	statFindLessonTestsQuery  = "SELECT * FROM public.test WHERE lesson_id = $1"
)

func (p *PostgresStatRepo) FindLessonStat(ctx context.Context,
	userID, lessonID domain.ID) (domain.LessonStat, error) {
	var pgLessonStat entity.PgLessonStat
	if err := p.db.GetContext(ctx, &pgLessonStat, statFindByUserLessonQuery, userID, lessonID); err != nil {
		if err == sql.ErrNoRows {
			return domain.LessonStat{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.LessonStat{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	lessonStat := pgLessonStat.ToDomain()

	var pgTests []entity.PgTest
	if err := p.db.SelectContext(ctx, &pgTests, statFindLessonTestsQuery, lessonID); err != nil {
		if err == sql.ErrNoRows {
			return lessonStat, nil
		} else {
			return domain.LessonStat{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	testStats := make([]domain.TestStat, len(pgTests))
	for i, pgTest := range pgTests {
		test := pgTest.ToDomain()
		var pgTestStat entity.PgTestStat
		if err := p.db.GetContext(ctx, &pgTestStat, statFindByUserTestQuery, userID, test.ID); err != nil {
			if err == sql.ErrNoRows {
				return domain.LessonStat{}, errors.Wrap(errs.ErrNotExist, err.Error())
			} else {
				return domain.LessonStat{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		}
		testStats[i] = pgTestStat.ToDomain()
	}
	lessonStat.TestStats = testStats

	return lessonStat, nil
}

func (p *PostgresStatRepo) CreateLessonStat(ctx context.Context, stat domain.LessonStat) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(errs.ErrTransactionError, err.Error())
	}

	var pgLessonStat = entity.NewPgLessonStat(stat)
	queryString := entity.InsertQueryString(pgLessonStat, "lesson_stat")
	_, err = tx.NamedExecContext(ctx, queryString, pgLessonStat)
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == PgUniqueViolationCode {
				return errors.Wrap(errs.ErrDuplicate, err.Error())
			} else {
				return errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		} else {
			return errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	for _, testStat := range stat.TestStats {
		var pgTestStat = entity.NewPgTestStat(testStat)
		queryString = entity.InsertQueryString(pgTestStat, "test_stat")
		_, err = tx.NamedExecContext(ctx, queryString, pgTestStat)
		if err != nil {
			tx.Rollback()
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == PgUniqueViolationCode {
					return errors.Wrap(errs.ErrDuplicate, err.Error())
				} else {
					return errors.Wrap(errs.ErrPersistenceFailed, err.Error())
				}
			} else {
				return errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(errs.ErrTransactionError, err.Error())
	}

	return nil
}

func (p *PostgresStatRepo) UpdateLessonStat(ctx context.Context, stat domain.LessonStat) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(errs.ErrTransactionError, err.Error())
	}

	var pgLessonStat = entity.NewPgLessonStat(stat)
	queryString := entity.UpdateQueryString(pgLessonStat, "lesson_stat")
	_, err = tx.NamedExecContext(ctx, queryString, pgLessonStat)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(errs.ErrUpdateFailed, err.Error())
	}

	for _, testStat := range stat.TestStats {
		var pgTestStat = entity.NewPgTestStat(testStat)
		queryString = entity.UpdateQueryString(pgTestStat, "test_stat")
		_, err = tx.NamedExecContext(ctx, queryString, pgTestStat)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(errs.ErrUpdateFailed, err.Error())
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(errs.ErrTransactionError, err.Error())
	}

	return nil
}
