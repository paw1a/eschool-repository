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

type PostgresCourseRepo struct {
	db *sqlx.DB
}

func NewCourseRepo(db *sqlx.DB) *PostgresCourseRepo {
	return &PostgresCourseRepo{
		db: db,
	}
}

const (
	courseFindAllQuery            = "SELECT * FROM public.course ORDER BY id"
	courseFindByIDQuery           = "SELECT * FROM public.course WHERE id = $1"
	courseFindStudentCoursesQuery = "SELECT c.* FROM public.course c " +
		"JOIN public.course_student cs on c.id = cs.course_id " +
		"JOIN public.user u on cs.student_id = u.id WHERE u.id = $1"
	courseFindTeacherCoursesQuery = "SELECT c.* FROM public.course c " +
		"JOIN public.course_teacher ct on c.id = ct.course_id " +
		"JOIN public.user u on ct.teacher_id = u.id WHERE u.id = $1"
	courseFindCourseTeachersQuery = "SELECT u.* FROM public.user u " +
		"JOIN public.course_teacher ct on u.id = ct.teacher_id " +
		"JOIN public.course c on ct.course_id = c.id WHERE c.id = $1"
	courseContainsStudentQuery = "SELECT EXISTS (SELECT 1 FROM public.course_student " +
		"WHERE course_id = $1 AND student_id = $2)"
	courseContainsTeacherQuery = "SELECT EXISTS (SELECT 1 FROM public.course_teacher " +
		"WHERE course_id = $1 AND teacher_id = $2)"
	courseAddCourseStudentQuery = "INSERT INTO public.course_student (student_id, course_id) " +
		"VALUES ($1, $2)"
	courseAddCourseTeacherQuery = "INSERT INTO public.course_teacher (teacher_id, course_id) " +
		"VALUES ($1, $2)"
	courseDeleteQuery = "DELETE FROM public.course WHERE id = $1"
)

func (p *PostgresCourseRepo) FindAll(ctx context.Context) ([]domain.Course, error) {
	var pgCourses []entity.PgCourse
	if err := p.db.SelectContext(ctx, &pgCourses, courseFindAllQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	courses := make([]domain.Course, len(pgCourses))
	for i, course := range pgCourses {
		courses[i] = course.ToDomain()
	}
	return courses, nil
}

func (p *PostgresCourseRepo) FindByID(ctx context.Context, courseID domain.ID) (domain.Course, error) {
	var pgCourse entity.PgCourse
	if err := p.db.GetContext(ctx, &pgCourse, courseFindByIDQuery, courseID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Course{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Course{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	return pgCourse.ToDomain(), nil
}

func (p *PostgresCourseRepo) FindStudentCourses(ctx context.Context, studentID domain.ID) ([]domain.Course, error) {
	var pgCourses []entity.PgCourse
	if err := p.db.SelectContext(ctx, &pgCourses, courseFindStudentCoursesQuery, studentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	courses := make([]domain.Course, len(pgCourses))
	for i, course := range pgCourses {
		courses[i] = course.ToDomain()
	}
	return courses, nil
}

func (p *PostgresCourseRepo) FindTeacherCourses(ctx context.Context, teacherID domain.ID) ([]domain.Course, error) {
	var pgCourses []entity.PgCourse
	if err := p.db.SelectContext(ctx, &pgCourses, courseFindTeacherCoursesQuery, teacherID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	courses := make([]domain.Course, len(pgCourses))
	for i, course := range pgCourses {
		courses[i] = course.ToDomain()
	}
	return courses, nil
}

func (p *PostgresCourseRepo) FindCourseTeachers(ctx context.Context, courseID domain.ID) ([]domain.User, error) {
	var pgUsers []entity.PgUser
	if err := p.db.SelectContext(ctx, &pgUsers, courseFindCourseTeachersQuery, courseID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return nil, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	teachers := make([]domain.User, len(pgUsers))
	for i, teacher := range pgUsers {
		teachers[i] = teacher.ToDomain()
	}
	return teachers, nil
}

func (p *PostgresCourseRepo) IsCourseStudent(ctx context.Context, studentID, courseID domain.ID) (bool, error) {
	var exists bool
	err := p.db.GetContext(ctx, &exists, courseContainsStudentQuery, courseID, studentID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PostgresCourseRepo) IsCourseTeacher(ctx context.Context, teacherID, courseID domain.ID) (bool, error) {
	var exists bool
	err := p.db.GetContext(ctx, &exists, courseContainsTeacherQuery, courseID, teacherID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PostgresCourseRepo) AddCourseStudent(ctx context.Context, studentID, courseID domain.ID) error {
	_, err := p.db.ExecContext(ctx, courseAddCourseStudentQuery, studentID, courseID)
	if err != nil {
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
	return nil
}

func (p *PostgresCourseRepo) AddCourseTeacher(ctx context.Context, teacherID, courseID domain.ID) error {
	_, err := p.db.ExecContext(ctx, courseAddCourseTeacherQuery, teacherID, courseID)
	if err != nil {
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
	return nil
}

func (p *PostgresCourseRepo) Create(ctx context.Context, course domain.Course) (domain.Course, error) {
	var pgCourse = entity.NewPgCourse(course)
	queryString := entity.InsertQueryString(pgCourse, "course")
	_, err := p.db.NamedExecContext(ctx, queryString, pgCourse)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == PgUniqueViolationCode {
				return domain.Course{}, errors.Wrap(errs.ErrDuplicate, err.Error())
			} else {
				return domain.Course{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
			}
		} else {
			return domain.Course{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	var createdCourse entity.PgCourse
	err = p.db.GetContext(ctx, &createdCourse, courseFindByIDQuery, pgCourse.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Course{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Course{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	return createdCourse.ToDomain(), nil
}

func (p *PostgresCourseRepo) Update(ctx context.Context, course domain.Course) (domain.Course, error) {
	var pgCourse = entity.NewPgCourse(course)
	queryString := entity.UpdateQueryString(pgCourse, "course")
	_, err := p.db.NamedExecContext(ctx, queryString, pgCourse)
	if err != nil {
		return domain.Course{}, errors.Wrap(errs.ErrUpdateFailed, err.Error())
	}

	var updatedCourse entity.PgCourse
	err = p.db.GetContext(ctx, &updatedCourse, courseFindByIDQuery, pgCourse.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Course{}, errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return domain.Course{}, errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}
	return updatedCourse.ToDomain(), nil
}

func (p *PostgresCourseRepo) UpdateStatus(ctx context.Context, courseID domain.ID, status domain.CourseStatus) error {
	var pgCourse entity.PgCourse
	err := p.db.GetContext(ctx, &pgCourse, courseFindByIDQuery, courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrap(errs.ErrNotExist, err.Error())
		} else {
			return errors.Wrap(errs.ErrPersistenceFailed, err.Error())
		}
	}

	course := pgCourse.ToDomain()
	course.Status = status
	pgCourse = entity.NewPgCourse(course)

	queryString := entity.UpdateQueryString(pgCourse, "course")
	_, err = p.db.NamedExecContext(ctx, queryString, pgCourse)
	if err != nil {
		return errors.Wrap(errs.ErrUpdateFailed, err.Error())
	}

	return nil
}

func (p *PostgresCourseRepo) Delete(ctx context.Context, courseID domain.ID) error {
	_, err := p.db.ExecContext(ctx, courseDeleteQuery, courseID)
	if err != nil {
		return errors.Wrap(errs.ErrDeleteFailed, err.Error())
	}
	return nil
}
