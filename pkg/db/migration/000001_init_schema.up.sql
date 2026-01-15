-- users
CREATE TABLE IF NOT EXISTS "users" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
    "first_name"    VARCHAR(255)    NOT NULL,
    "last_name"     VARCHAR(255)    NOT NULL,
    "email"         VARCHAR(255)    NOT NULL,
    "password_hash" VARCHAR(64)     NOT NULL,
    "is_verified"	BOOLEAN			NOT NULL	DEFAULT FALSE,
    "streak"        INTEGER         NOT NULL    DEFAULT 0,
    "score"         INTEGER         NOT NULL    DEFAULT 0,
    "is_admin"      BOOLEAN         NOT NULL    DEFAULT FALSE
);

CREATE UNIQUE   INDEX uq_users_email        ON "users" ("email");
CREATE          INDEX idx_users_deleted_at  ON "users" ("deleted_at");
CREATE          INDEX idx_users_created_at  ON "users" ("created_at");

-- courses
CREATE TABLE IF NOT EXISTS "courses" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
	"subject"		VARCHAR(255)	NOT NULL,
	"price"		    SMALLINT		NOT NULL,
	"discount"	    SMALLINT		NOT NULL    DEFAULT 0,
	"is_active"	    BOOLEAN			NOT NULL	DEFAULT FALSE,
	"difficulty"    SMALLINT		NOT NULL	DEFAULT 0
);

CREATE INDEX idx_courses_deleted_at  ON "courses" ("deleted_at");
CREATE INDEX idx_courses_created_at  ON "courses" ("created_at");
CREATE INDEX idx_courses_difficulty  ON "courses" ("difficulty");

-- course_translations
CREATE TABLE IF NOT EXISTS "course_translations" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "course_uuid"   UUID            NOT NULL,
    "language"      VARCHAR(2)      NOT NULL,
    "name"          VARCHAR(255)    NOT NULL,
    "description"   TEXT            NOT NULL,
    "bullets"       TEXT            NOT NULL    DEFAULT '',
    CONSTRAINT fk_course_translations_course FOREIGN KEY ("course_uuid") REFERENCES "courses"("uuid") ON DELETE CASCADE
);

CREATE INDEX        idx_course_translations_language        ON "course_translations" ("language");
CREATE UNIQUE INDEX uq_course_translations_course_language  ON "course_translations" ("course_uuid", "language");

-- lessons
CREATE TABLE IF NOT EXISTS "lessons" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
    "course_uuid"   UUID            NOT NULL,
    "order_index"	SMALLINT		NOT NULL,
    "is_public"	    BOOLEAN			NOT NULL	DEFAULT FALSE,
    CONSTRAINT fk_lessons_course FOREIGN KEY ("course_uuid") REFERENCES "courses"("uuid") ON DELETE CASCADE
);

CREATE INDEX idx_lessons_deleted_at  ON "lessons" ("deleted_at");
CREATE INDEX idx_lessons_created_at  ON "lessons" ("created_at");
CREATE INDEX idx_lessons_course_uuid ON "lessons" ("course_uuid");
CREATE INDEX idx_lessons_order_index ON "lessons" ("order_index");

-- lesson_translations
CREATE TABLE IF NOT EXISTS "lesson_translations" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "lesson_uuid"   UUID            NOT NULL,
    "language"      VARCHAR(2)      NOT NULL,
    "name"          VARCHAR(255)    NOT NULL,
    "description"   TEXT            NOT NULL,
    CONSTRAINT fk_lesson_translations_lesson FOREIGN KEY ("lesson_uuid") REFERENCES "lessons"("uuid") ON DELETE CASCADE
);

CREATE INDEX        idx_lesson_translations_language        ON "lesson_translations" ("language");
CREATE UNIQUE INDEX uq_lesson_translations_lesson_language  ON "lesson_translations" ("lesson_uuid", "language");
CREATE TYPE exercise_type AS ENUM ('quiz', 'code');

-- exercises
CREATE TABLE IF NOT EXISTS "exercises" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
    "lesson_uuid"   UUID            NOT NULL,
    "order_index"	SMALLINT		NOT NULL,
    "reward"	    SMALLINT		NOT NULL,
    "type"          exercise_type   NOT NULL,
    "data"		    JSONB			NOT NULL,
    CONSTRAINT fk_exercises_lesson FOREIGN KEY ("lesson_uuid") REFERENCES "lessons"("uuid") ON DELETE CASCADE
);

CREATE INDEX idx_exercises_deleted_at  ON "exercises" ("deleted_at");
CREATE INDEX idx_exercises_created_at  ON "exercises" ("created_at");
CREATE INDEX idx_exercises_lesson_uuid ON "exercises" ("lesson_uuid");
CREATE INDEX idx_exercises_order_index ON "exercises" ("order_index");

-- exercise_translations
CREATE TABLE IF NOT EXISTS "exercise_translations" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "exercise_uuid" UUID            NOT NULL,
    "language"      VARCHAR(2)      NOT NULL,
    "name"          VARCHAR(255)    NOT NULL,
    "description"   TEXT            NOT NULL,
    CONSTRAINT fk_exercise_translations_exercise FOREIGN KEY ("exercise_uuid") REFERENCES "exercises"("uuid") ON DELETE CASCADE
);

CREATE INDEX        idx_exercise_translations_language          ON "exercise_translations" ("language");
CREATE UNIQUE INDEX uq_exercise_translations_exercise_language  ON "exercise_translations" ("exercise_uuid", "language");

-- user_courses
CREATE TABLE IF NOT EXISTS "user_courses" (
    "uuid"              UUID            PRIMARY KEY DEFAULT uuidv7(),
    "started_at"        TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "last_accessed_at"  TIMESTAMP       NULL,
    "user_uuid"         UUID            NOT NULL,
    "course_uuid"       UUID            NOT NULL,
    "completed_at"      TIMESTAMP       NULL,
    
    CONSTRAINT fk_user_courses_user FOREIGN KEY ("user_uuid") REFERENCES "users"("uuid") ON DELETE CASCADE,
    CONSTRAINT fk_user_courses_course FOREIGN KEY ("course_uuid") REFERENCES "courses"("uuid") ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_user_courses_user_course ON "user_courses" ("user_uuid", "course_uuid");
CREATE INDEX        idx_user_courses_user       ON "user_courses" ("user_uuid");
CREATE INDEX        idx_user_courses_course     ON "user_courses" ("course_uuid");

-- user_lessons
CREATE TABLE IF NOT EXISTS "user_lessons" (
    "uuid"              UUID            PRIMARY KEY DEFAULT uuidv7(),
    "started_at"        TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "last_accessed_at"  TIMESTAMP       NULL,
    "user_uuid"         UUID            NOT NULL,
    "lesson_uuid"       UUID            NOT NULL,
    "completed_at"      TIMESTAMP       NULL,
    
    CONSTRAINT fk_user_lessons_user     FOREIGN KEY ("user_uuid")   REFERENCES "users"("uuid") ON DELETE CASCADE,
    CONSTRAINT fk_user_lessons_lesson   FOREIGN KEY ("lesson_uuid") REFERENCES "lessons"("uuid") ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_user_lessons_user_lesson ON "user_lessons" ("user_uuid", "lesson_uuid");
CREATE INDEX        idx_user_lessons_user       ON "user_lessons" ("user_uuid");
CREATE INDEX        idx_user_lessons_lesson     ON "user_lessons" ("lesson_uuid");

-- user_exercises
CREATE TABLE IF NOT EXISTS "user_exercises" (
    "uuid"              UUID            PRIMARY KEY DEFAULT uuidv7(),
    "started_at"        TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "last_accessed_at"  TIMESTAMP       NULL,
    "user_uuid"         UUID            NOT NULL,
    "exercise_uuid"     UUID            NOT NULL,
    "submission"        JSONB           NOT NULL,
    "attempts"          INTEGER         NOT NULL    DEFAULT 0,
    "completed_at"      TIMESTAMP       NULL,
    
    CONSTRAINT fk_user_exercises_user       FOREIGN KEY ("user_uuid")       REFERENCES "users"("uuid") ON DELETE CASCADE,
    CONSTRAINT fk_user_exercises_exercise   FOREIGN KEY ("exercise_uuid")   REFERENCES "exercises"("uuid") ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_user_exercises_user_exercise ON "user_exercises" ("user_uuid", "exercise_uuid");
CREATE INDEX        idx_user_exercises_user         ON "user_exercises" ("user_uuid");
CREATE INDEX        idx_user_exercises_exercise     ON "user_exercises" ("exercise_uuid");