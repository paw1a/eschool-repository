package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/paw1a/eschool-core/domain"
	"strings"
)

const (
	PgLessonTheory   = "theory"
	PgLessonVideo    = "video"
	PgLessonPractice = "practice"
)

type PgLesson struct {
	ID       uuid.UUID `db:"id"`
	CourseID uuid.UUID `db:"course_id"`
	Title    string    `db:"title"`
	Score    int       `db:"score"`
	Type     string    `db:"type"`

	TheoryUrl null.String `db:"theory_url"`
	VideoUrl  null.String `db:"video_url"`
}

type PgTest struct {
	ID       uuid.UUID `db:"id"`
	LessonID uuid.UUID `db:"lesson_id"`
	TaskUrl  string    `db:"task_url"`
	Options  string    `db:"options"`
	Answer   string    `db:"answer"`
	Level    int       `db:"level"`
	Score    int       `db:"score"`
}

func (s *PgLesson) ToDomain() domain.Lesson {
	var lessonType domain.LessonType
	switch s.Type {
	case PgLessonTheory:
		lessonType = domain.TheoryLesson
	case PgLessonVideo:
		lessonType = domain.VideoLesson
	case PgLessonPractice:
		lessonType = domain.PracticeLesson
	}

	return domain.Lesson{
		ID:        domain.ID(s.ID.String()),
		CourseID:  domain.ID(s.CourseID.String()),
		Title:     s.Title,
		Score:     s.Score,
		Type:      lessonType,
		TheoryUrl: s.TheoryUrl,
		VideoUrl:  s.VideoUrl,
		Tests:     nil,
	}
}

func NewPgLesson(lesson domain.Lesson) PgLesson {
	id, _ := uuid.Parse(lesson.ID.String())
	courseID, _ := uuid.Parse(lesson.CourseID.String())
	var lessonType string
	switch lesson.Type {
	case domain.TheoryLesson:
		lessonType = PgLessonTheory
	case domain.VideoLesson:
		lessonType = PgLessonVideo
	case domain.PracticeLesson:
		lessonType = PgLessonPractice
	}

	return PgLesson{
		ID:        id,
		CourseID:  courseID,
		Title:     lesson.Title,
		Score:     lesson.Score,
		Type:      lessonType,
		TheoryUrl: lesson.TheoryUrl,
		VideoUrl:  lesson.VideoUrl,
	}
}

func (s *PgTest) ToDomain() domain.Test {
	options := strings.Split(s.Options, "\n")
	return domain.Test{
		ID:       domain.ID(s.ID.String()),
		LessonID: domain.ID(s.LessonID.String()),
		TaskUrl:  s.TaskUrl,
		Options:  options,
		Answer:   s.Answer,
		Level:    s.Level,
		Score:    s.Score,
	}
}

func NewPgTest(test domain.Test) PgTest {
	id, _ := uuid.Parse(test.ID.String())
	lessonID, _ := uuid.Parse(test.LessonID.String())
	options := strings.Join(test.Options, "\n")
	return PgTest{
		ID:       id,
		LessonID: lessonID,
		TaskUrl:  test.TaskUrl,
		Options:  options,
		Answer:   test.Answer,
		Level:    test.Level,
		Score:    test.Score,
	}
}
