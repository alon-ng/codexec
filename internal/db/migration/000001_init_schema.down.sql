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



