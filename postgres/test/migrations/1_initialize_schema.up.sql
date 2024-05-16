create table public.user (
    id uuid primary key,
    email varchar(255) unique not null,
    password varchar(255) not null,
    name varchar(255) not null,
    surname varchar(255) not null,
    phone varchar(32),
    city varchar(255),
    avatar_url text
);

create table public.school (
    id uuid primary key,
    name varchar(255) not null,
    description text not null,
    owner_id uuid not null,
    foreign key (owner_id) references public.user(id) on delete cascade
);

create type course_status as enum ('draft', 'ready', 'published');

create table public.course (
    id uuid primary key,
    name varchar(255) not null,
    school_id uuid not null,
    level int not null,
    price bigint not null,
    language varchar(255) not null,
    status course_status not null,
    foreign key (school_id) references public.school(id) on delete cascade
);

create type lesson_type as enum ('theory', 'practice', 'video');

create table public.lesson (
    id uuid primary key,
    title varchar(255) not null,
    type lesson_type not null,
    score int not null,
    theory_url text,
    video_url text,
    course_id uuid not null,
    foreign key (course_id) references public.course(id) on delete cascade
);

create table public.test (
    id uuid primary key,
    task_url text not null,
    options text not null,
    answer text not null,
    score int not null,
    level int not null,
    lesson_id uuid,
    foreign key (lesson_id) references public.lesson(id) on delete cascade
);

create table public.review (
    id uuid primary key,
    text text not null,
    course_id uuid not null,
    user_id uuid,
    foreign key (course_id) references public.course(id) on delete cascade,
    foreign key (user_id) references public.user(id) on delete set null
);

create type certificate_grade as enum ('bronze', 'silver', 'gold');

create table public.certificate (
    id uuid primary key,
    name varchar(1024) not null,
    score int not null,
    grade certificate_grade not null,
    created_at timestamp not null,
    user_id uuid not null,
    course_id uuid not null,
    foreign key (user_id) references public.user(id) on delete cascade,
    foreign key (course_id) references public.course(id) on delete cascade
);

create table lesson_stat (
    id uuid primary key,
    score int not null,
    user_id uuid not null,
    lesson_id uuid not null,
    foreign key (user_id) references public.user(id) on delete cascade,
    foreign key (lesson_id) references public.lesson(id) on delete cascade
);

create table test_stat (
    id uuid primary key,
    score int not null,
    user_id uuid not null,
    test_id uuid not null,
    foreign key (user_id) references public.user(id) on delete cascade,
    foreign key (test_id) references public.test(id) on delete cascade
);

create table course_student (
    student_id uuid not null,
    course_id uuid not null,
    primary key (student_id, course_id),
    foreign key (student_id) references public.user(id) on delete cascade,
    foreign key (course_id) references public.course(id) on delete cascade
);

create table course_teacher (
    teacher_id uuid not null,
    course_id uuid not null,
    primary key (teacher_id, course_id),
    foreign key (teacher_id) references public.user(id) on delete cascade,
    foreign key (course_id) references public.course(id) on delete cascade
);

create table school_teacher (
    teacher_id uuid not null,
    school_id uuid not null,
    primary key (teacher_id, school_id),
    foreign key (teacher_id) references public.user(id) on delete cascade,
    foreign key (school_id) references public.school(id) on delete cascade
);
