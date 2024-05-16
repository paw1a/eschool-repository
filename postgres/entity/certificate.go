package entity

import (
	"github.com/google/uuid"
	"github.com/paw1a/eschool-core/domain"
	"time"
)

const (
	PgBronzeCertificate = "bronze"
	PgSilverCertificate = "silver"
	PgGoldCertificate   = "gold"
)

type PgCertificate struct {
	ID        uuid.UUID `db:"id"`
	CourseID  uuid.UUID `db:"course_id"`
	UserID    uuid.UUID `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	Grade     string    `db:"grade"`
	Score     int       `db:"score"`
}

func (s *PgCertificate) ToDomain() domain.Certificate {
	var grade domain.CertificateGrade
	switch s.Grade {
	case PgBronzeCertificate:
		grade = domain.BronzeCertificate
	case PgSilverCertificate:
		grade = domain.SilverCertificate
	case PgGoldCertificate:
		grade = domain.GoldCertificate
	}

	return domain.Certificate{
		ID:        domain.ID(s.ID.String()),
		CourseID:  domain.ID(s.CourseID.String()),
		UserID:    domain.ID(s.UserID.String()),
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		Grade:     grade,
		Score:     s.Score,
	}
}

func NewPgCertificate(certificate domain.Certificate) PgCertificate {
	id, _ := uuid.Parse(certificate.ID.String())
	courseID, _ := uuid.Parse(certificate.CourseID.String())
	userID, _ := uuid.Parse(certificate.UserID.String())
	var grade string
	switch certificate.Grade {
	case domain.BronzeCertificate:
		grade = PgBronzeCertificate
	case domain.SilverCertificate:
		grade = PgSilverCertificate
	case domain.GoldCertificate:
		grade = PgGoldCertificate
	}

	return PgCertificate{
		ID:        id,
		CourseID:  courseID,
		UserID:    userID,
		Name:      certificate.Name,
		CreatedAt: certificate.CreatedAt,
		Grade:     grade,
		Score:     certificate.Score,
	}
}
