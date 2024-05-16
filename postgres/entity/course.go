package entity

import (
	"github.com/google/uuid"
	"github.com/paw1a/eschool-core/domain"
)

const (
	PgCourseDraft     = "draft"
	PgCourseReady     = "ready"
	PgCoursePublished = "published"
)

type PgCourse struct {
	ID       uuid.UUID `db:"id"`
	SchoolID uuid.UUID `db:"school_id"`
	Name     string    `db:"name"`
	Level    int       `db:"level"`
	Price    int64     `db:"price"`
	Language string    `db:"language"`
	Status   string    `db:"status"`
}

func (s *PgCourse) ToDomain() domain.Course {
	var status domain.CourseStatus
	switch s.Status {
	case PgCourseDraft:
		status = domain.CourseDraft
	case PgCourseReady:
		status = domain.CourseReady
	case PgCoursePublished:
		status = domain.CoursePublished
	}

	return domain.Course{
		ID:       domain.ID(s.ID.String()),
		SchoolID: domain.ID(s.SchoolID.String()),
		Name:     s.Name,
		Level:    s.Level,
		Price:    s.Price,
		Language: s.Language,
		Status:   status,
	}
}

func NewPgCourse(course domain.Course) PgCourse {
	id, _ := uuid.Parse(course.ID.String())
	schoolID, _ := uuid.Parse(course.SchoolID.String())
	var status string
	switch course.Status {
	case domain.CourseDraft:
		status = PgCourseDraft
	case domain.CourseReady:
		status = PgCourseReady
	case domain.CoursePublished:
		status = PgCoursePublished
	}

	return PgCourse{
		ID:       id,
		SchoolID: schoolID,
		Name:     course.Name,
		Level:    course.Level,
		Price:    course.Price,
		Language: course.Language,
		Status:   status,
	}
}
