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

-- chat_messages
DROP INDEX IF EXISTS idx_chat_messages_exercise;
DROP INDEX IF EXISTS idx_chat_messages_user;
DROP INDEX IF EXISTS idx_chat_messages_ts;

DROP TABLE IF EXISTS chat_messages;

-- users
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS uq_users_email;

DROP TABLE IF EXISTS users;

-- exercise_translations
DROP INDEX IF EXISTS idx_exercise_translations_language;
DROP INDEX IF EXISTS uq_exercise_translations_exercise_language;
DROP TABLE IF EXISTS exercise_translations;

-- exercises
DROP INDEX IF EXISTS idx_exercises_deleted_at;
DROP INDEX IF EXISTS idx_exercises_created_at;
DROP INDEX IF EXISTS idx_exercises_lesson_uuid;
DROP INDEX IF EXISTS idx_exercises_order_index;

DROP TABLE IF EXISTS exercises;
DROP TYPE IF EXISTS exercise_type;

-- lesson_translations
DROP INDEX IF EXISTS idx_lesson_translations_language;
DROP INDEX IF EXISTS uq_lesson_translations_lesson_language;
DROP TABLE IF EXISTS lesson_translations;

-- lessons
DROP INDEX IF EXISTS idx_lessons_deleted_at;
DROP INDEX IF EXISTS idx_lessons_created_at;
DROP INDEX IF EXISTS idx_lessons_course_uuid;
DROP INDEX IF EXISTS idx_lessons_order_index;

DROP TABLE IF EXISTS lessons;

-- course_translations
DROP INDEX IF EXISTS idx_course_translations_language;
DROP INDEX IF EXISTS uq_course_translations_course_language;
DROP TABLE IF EXISTS course_translations;

-- courses
DROP INDEX IF EXISTS idx_courses_deleted_at;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_difficulty;

DROP TABLE IF EXISTS courses;