package repository

import (
	"context"
	"github.com/paw1a/eschool-core/domain"
	repository "github.com/paw1a/eschool-repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

var courses = []domain.Course{
	domain.Course{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca"),
		SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cc"),
		Name:     "course1",
		Level:    4,
		Price:    1200,
		Language: "russian",
		Status:   domain.CourseDraft,
	},
	domain.Course{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cb"),
		SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cc"),
		Name:     "course2",
		Level:    2,
		Price:    1500,
		Language: "english",
		Status:   domain.CoursePublished,
	},
	domain.Course{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cc"),
		SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cd"),
		Name:     "course3",
		Level:    3,
		Price:    12000,
		Language: "russian",
		Status:   domain.CourseReady,
	},
	domain.Course{
		ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027cd"),
		SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cd"),
		Name:     "course4",
		Level:    2,
		Price:    0,
		Language: "english",
		Status:   domain.CoursePublished,
	},
}

var studentCoursesID = domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027ca")
var teacherCoursesID = domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cb")
var newUserID = domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7027cc")

var createdCourse = domain.Course{
	ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ce"),
	SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cd"),
	Name:     "course5",
	Level:    1,
	Price:    0,
	Language: "english",
	Status:   domain.CourseDraft,
}

var updatedCourse = domain.Course{
	ID:       domain.ID("30e18bc1-4354-4937-9a4d-03cf0b7027ca"),
	SchoolID: domain.ID("30e18bc1-4354-4937-9a3b-03cf0b7034cc"),
	Name:     "course new name",
	Level:    4,
	Price:    50000,
	Language: "russian",
	Status:   domain.CoursePublished,
}

func TestCourseRepository(t *testing.T) {
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

	t.Run("test find all courses", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		found, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to find all courses: %v", err)
		}
		require.Equal(t, len(found), len(courses))
	})

	t.Run("test find course by id", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		course, err := repo.FindByID(ctx, courses[0].ID)
		if err != nil {
			t.Errorf("failed to find course with id: %v", err)
		}
		require.Equal(t, course, courses[0])
	})

	t.Run("test find student courses", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		found, err := repo.FindStudentCourses(ctx, studentCoursesID)
		if err != nil {
			t.Errorf("failed to find student courses: %v", err)
		}
		require.Equal(t, len(found), 2)
		require.Equal(t, found[0], courses[0])
		require.Equal(t, found[1], courses[1])
	})

	t.Run("test find teacher courses", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		found, err := repo.FindTeacherCourses(ctx, teacherCoursesID)
		if err != nil {
			t.Errorf("failed to find teacher courses: %v", err)
		}
		require.Equal(t, len(found), 2)
		require.Equal(t, found[0], courses[0])
		require.Equal(t, found[1], courses[1])
	})

	t.Run("test add course student", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		err = repo.AddCourseStudent(ctx, newUserID, courses[0].ID)
		if err != nil {
			t.Errorf("failed to add course student: %v", err)
		}
	})

	t.Run("test add course teacher", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		err = repo.AddCourseTeacher(ctx, newUserID, courses[0].ID)
		if err != nil {
			t.Errorf("failed to add course teacher: %v", err)
		}
	})

	t.Run("test create course", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		course, err := repo.Create(ctx, createdCourse)
		if err != nil {
			t.Errorf("failed to create course: %v", err)
		}
		require.Equal(t, course, createdCourse)
	})

	t.Run("test update course", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		course, err := repo.Update(ctx, updatedCourse)
		if err != nil {
			t.Errorf("failed to create course: %v", err)
		}
		require.Equal(t, course, updatedCourse)
	})

	t.Run("test delete course", func(t *testing.T) {
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

		repo := repository.NewCourseRepo(db)
		err = repo.Delete(ctx, courses[0].ID)
		if err != nil {
			t.Errorf("failed to delete course: %v", err)
		}
	})
}
