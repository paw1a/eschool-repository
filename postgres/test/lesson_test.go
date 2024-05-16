package repository

import (
	"context"
	"github.com/guregu/null"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var lessonCourseID = domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca")
var lessons = []domain.Lesson{
	domain.Lesson{
		ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022ca"),
		CourseID:  lessonCourseID,
		Title:     "lesson1",
		Score:     10,
		Type:      domain.TheoryLesson,
		TheoryUrl: null.StringFrom("url"),
		VideoUrl:  null.String{},
		Tests:     nil,
	},
	domain.Lesson{
		ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cb"),
		CourseID:  lessonCourseID,
		Title:     "lesson2",
		Score:     10,
		Type:      domain.VideoLesson,
		TheoryUrl: null.String{},
		VideoUrl:  null.StringFrom("url"),
		Tests:     nil,
	},
	domain.Lesson{
		ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cc"),
		CourseID:  lessonCourseID,
		Title:     "lesson3",
		Score:     10,
		Type:      domain.PracticeLesson,
		TheoryUrl: null.String{},
		VideoUrl:  null.String{},
		Tests:     tests,
	},
}

var tests = []domain.Test{
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027ca"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cc"),
		TaskUrl:  "url",
		Options:  []string{"opt1", "opt2", "opt3"},
		Answer:   "opt1",
		Level:    2,
		Score:    12,
	},
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027cb"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cc"),
		TaskUrl:  "url",
		Options:  []string{"opt1", "opt2"},
		Answer:   "opt2",
		Level:    2,
		Score:    12,
	},
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7027cc"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cc"),
		TaskUrl:  "url",
		Options:  []string{"opt1"},
		Answer:   "opt1",
		Level:    2,
		Score:    12,
	},
}

var createdLesson = domain.Lesson{
	ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cd"),
	CourseID:  lessonCourseID,
	Title:     "created lesson 4",
	Score:     100,
	Type:      domain.PracticeLesson,
	TheoryUrl: null.String{},
	VideoUrl:  null.String{},
	Tests:     createdTests,
}

var updatedLesson = domain.Lesson{
	ID:        domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cc"),
	CourseID:  lessonCourseID,
	Title:     "updated lesson 3",
	Score:     20,
	Type:      domain.PracticeLesson,
	TheoryUrl: null.String{},
	VideoUrl:  null.String{},
	Tests:     tests,
}

var createdTests = []domain.Test{
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7025ca"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cd"),
		TaskUrl:  "url",
		Options:  []string{"opt1", "opt2", "opt3"},
		Answer:   "opt1",
		Level:    2,
		Score:    12,
	},
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7025cb"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cd"),
		TaskUrl:  "url",
		Options:  []string{"opt1", "opt2"},
		Answer:   "opt2",
		Level:    2,
		Score:    12,
	},
	domain.Test{
		ID:       domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7025cc"),
		LessonID: domain.ID("30e18bc1-4352-4937-9a3b-03cf0b7022cd"),
		TaskUrl:  "url",
		Options:  []string{"opt1"},
		Answer:   "opt1",
		Level:    2,
		Score:    12,
	},
}

func TestLessonRepository(t *testing.T) {
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

	t.Run("test find all lessons", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all lessons: %v", err)
		}
		require.Equal(t, len(found), len(lessons))
		for i, lesson := range found {
			if lesson.Type == domain.PracticeLesson {
				for j, test := range lesson.Tests {
					require.Equal(t, test, lessons[i].Tests[j])
					require.Equal(t, reflect.DeepEqual(test.Options,
						lessons[i].Tests[j].Options), true)
				}
			} else {
				require.Equal(t, lesson, lessons[i])
			}
		}
	})

	t.Run("test find lesson by id", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		lesson, err := repo.FindByID(ctx, lessons[2].ID)
		if err != nil {
			t.Errorf("failed to find course with id: %v", err)
		}

		for j, test := range lesson.Tests {
			require.Equal(t, test, lesson.Tests[j])
			require.Equal(t, reflect.DeepEqual(test.Options,
				lesson.Tests[j].Options), true)
		}
	})

	t.Run("test find course lessons", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		found, err := repo.FindCourseLessons(ctx, lessonCourseID)
		if err != nil {
			t.Errorf("failed to find course lessons: %v", err)
		}
		require.Equal(t, len(found), len(lessons))
		for i, lesson := range found {
			if lesson.Type == domain.PracticeLesson {
				for j, test := range lesson.Tests {
					require.Equal(t, test, lessons[i].Tests[j])
					require.Equal(t, reflect.DeepEqual(test.Options,
						lessons[i].Tests[j].Options), true)
				}
			} else {
				require.Equal(t, lesson, lessons[i])
			}
		}
	})

	t.Run("test find lesson tests", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		found, err := repo.FindLessonTests(ctx, lessons[2].ID)
		if err != nil {
			t.Errorf("failed to find lesson tests: %v", err)
		}
		require.Equal(t, len(found), len(tests))
		for i, test := range found {
			require.Equal(t, test, lessons[2].Tests[i])
			require.Equal(t, reflect.DeepEqual(test.Options,
				lessons[2].Tests[i].Options), true)
		}
	})

	t.Run("test create lesson", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		lesson, err := repo.Create(ctx, createdLesson)
		if err != nil {
			t.Errorf("failed to create lesson: %v", err)
		}
		require.Equal(t, lesson, createdLesson)
		for i, test := range lesson.Tests {
			require.Equal(t, test, createdLesson.Tests[i])
			require.Equal(t, reflect.DeepEqual(test.Options,
				createdLesson.Tests[i].Options), true)
		}
	})

	t.Run("test update lesson", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		lesson, err := repo.Update(ctx, updatedLesson)
		if err != nil {
			t.Errorf("failed to update lesson: %v", err)
		}
		require.Equal(t, lesson, updatedLesson)
		for i, test := range lesson.Tests {
			require.Equal(t, test, updatedLesson.Tests[i])
			require.Equal(t, reflect.DeepEqual(test.Options,
				updatedLesson.Tests[i].Options), true)
		}
	})

	t.Run("test delete lesson", func(t *testing.T) {
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

		repo := repository.NewLessonRepo(db)
		err = repo.Delete(ctx, lessons[0].ID)
		if err != nil {
			t.Errorf("failed to delete lesson: %v", err)
		}
	})
}
