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

type PostgresCertificateRepo struct {
	db *sqlx.DB
}

func NewCertificateRepo(db *sqlx.DB) *PostgresCertificateRepo {
	return &PostgresCertificateRepo{
		db: db,
	}
}

const (
	certificateFindAllQuery               = "SELECT * FROM public.certificate"
	certificateFindByIDQuery              = "SELECT * FROM public.certificate WHERE id = $1"
	certificateFindByCourseAndUserIDQuery = "SELECT * FROM public.certificate WHERE course_id = $1 AND user_id = $2"
	certificateFindUserCertificatesQuery  = "SELECT * FROM public.certificate WHERE user_id = $1"
)

func (p *PostgresCertificateRepo) FindAll(ctx context.Context) ([]domain.Certificate, error) {
	var pgCertificates []entity.PgCertificate
	if err := p.db.SelectContext(ctx, &pgCertificates, certificateFindAllQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	certificates := make([]domain.Certificate, len(pgCertificates))
	for i, certificate := range pgCertificates {
		certificates[i] = certificate.ToDomain()
	}
	return certificates, nil
}

func (p *PostgresCertificateRepo) FindByID(ctx context.Context,
	certID domain.ID) (domain.Certificate, error) {
	var pgCertificate entity.PgCertificate
	if err := p.db.GetContext(ctx, &pgCertificate, certificateFindByIDQuery, certID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Certificate{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Certificate{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	return pgCertificate.ToDomain(), nil
}

func (p *PostgresCertificateRepo) FindUserCertificates(ctx context.Context,
	userID domain.ID) ([]domain.Certificate, error) {
	var pgCertificates []entity.PgCertificate
	if err := p.db.SelectContext(ctx, &pgCertificates, certificateFindUserCertificatesQuery, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	certificates := make([]domain.Certificate, len(pgCertificates))
	for i, certificate := range pgCertificates {
		certificates[i] = certificate.ToDomain()
	}
	return certificates, nil
}

func (p *PostgresCertificateRepo) FindUserCourseCertificate(ctx context.Context,
	courseID, userID domain.ID) (domain.Certificate, error) {
	var pgCertificate entity.PgCertificate
	if err := p.db.GetContext(ctx, &pgCertificate, certificateFindByCourseAndUserIDQuery,
		courseID, userID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Certificate{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Certificate{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	return pgCertificate.ToDomain(), nil
}

func (p *PostgresCertificateRepo) Create(ctx context.Context,
	cert domain.Certificate) (domain.Certificate, error) {
	var pgCertificate = entity.NewPgCertificate(cert)
	queryString := entity.InsertQueryString(pgCertificate, "certificate")
	_, err := p.db.NamedExecContext(ctx, queryString, pgCertificate)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == PgUniqueViolationCode {
				return domain.Certificate{}, errors.Wrap(errs.ErrDuplicate, err.Error())
			} else if pgErr.Code == PgEnumValueError {
				return domain.Certificate{}, errors.Wrap(errs.ErrEnumValueError, err.Error())
			} else {
				return domain.Certificate{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		} else {
			return domain.Certificate{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	var createdCertificate entity.PgCertificate
	err = p.db.GetContext(ctx, &createdCertificate, certificateFindByIDQuery, pgCertificate.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Certificate{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Certificate{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	return createdCertificate.ToDomain(), nil
}
