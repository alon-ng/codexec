-- user_exercises
DROP INDEX IF EXISTS uq_user_exercises_user_exercise;
DROP INDEX IF EXISTS idx_user_exercises_user;
DROP INDEX IF EXISTS idx_user_exercises_exercise;
DROP TABLE IF EXISTS user_exercises;
DROP TYPE IF EXISTS exercise_status;

-- user_lessons
DROP INDEX IF EXISTS uq_user_lessons_user_lesson;
DROP INDEX IF EXISTS idx_user_lessons_user;
DROP INDEX IF EXISTS idx_user_lessons_lesson;

DROP TABLE IF EXISTS user_lessons;

-- user_courses
DROP INDEX IF EXISTS uq_user_courses_user_course;
DROP INDEX IF EXISTS idx_user_courses_user;
DROP INDEX IF EXISTS idx_user_courses_course;
DROP TABLE IF EXISTS user_courses;
DROP TYPE IF EXISTS course_status;

-- users
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS uq_users_email;

DROP TABLE IF EXISTS users;

-- exercises
DROP INDEX IF EXISTS idx_exercises_deleted_at;
DROP INDEX IF EXISTS idx_exercises_created_at;
DROP INDEX IF EXISTS idx_exercises_lesson_uuid;
DROP INDEX IF EXISTS idx_exercises_order_index;

DROP TABLE IF EXISTS exercises;
DROP TYPE IF EXISTS exercise_type;

-- lessons
DROP INDEX IF EXISTS idx_lessons_deleted_at;
DROP INDEX IF EXISTS idx_lessons_created_at;
DROP INDEX IF EXISTS idx_lessons_course_uuid;
DROP INDEX IF EXISTS idx_lessons_order_index;

DROP TABLE IF EXISTS lessons;

-- courses
DROP INDEX IF EXISTS idx_courses_deleted_at;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_difficulty;

DROP TABLE IF EXISTS courses;