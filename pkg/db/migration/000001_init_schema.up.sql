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
    "name"		    VARCHAR(255)	NOT NULL	UNIQUE,
	"description"	TEXT	    	NOT NULL,
	"subject"		VARCHAR(255)	NOT NULL,
	"price"		    SMALLINT		NOT NULL,
	"discount"	    SMALLINT		NOT NULL    DEFAULT 0,
	"is_active"	    BOOLEAN			NOT NULL	DEFAULT FALSE,
	"difficulty"    SMALLINT		NOT NULL	DEFAULT 0,
	"bullets"		TEXT        	NOT NULL	DEFAULT ''
);

CREATE          INDEX idx_courses_deleted_at  ON "courses" ("deleted_at");
CREATE          INDEX idx_courses_created_at  ON "courses" ("created_at");
CREATE          INDEX idx_courses_difficulty  ON "courses" ("difficulty");

-- lessons
CREATE TABLE IF NOT EXISTS "lessons" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
    "course_uuid"   UUID            NOT NULL,
    "name"		    VARCHAR(128)	NOT NULL	UNIQUE,
    "description"	TEXT	    	NOT NULL,
    "order_index"	SMALLINT		NOT NULL,
    "is_public"	    BOOLEAN			NOT NULL	DEFAULT FALSE,
    CONSTRAINT fk_lessons_course FOREIGN KEY ("course_uuid") REFERENCES "courses"("uuid")
);

CREATE          INDEX idx_lessons_deleted_at  ON "lessons" ("deleted_at");
CREATE          INDEX idx_lessons_created_at  ON "lessons" ("created_at");
CREATE          INDEX idx_lessons_course_uuid ON "lessons" ("course_uuid");
CREATE          INDEX idx_lessons_order_index ON "lessons" ("order_index");

CREATE TYPE exercise_type AS ENUM ('quiz', 'code');

-- exercises
CREATE TABLE IF NOT EXISTS "exercises" (
    "uuid"          UUID            PRIMARY KEY DEFAULT uuidv7(),
    "created_at"    TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "modified_at"   TIMESTAMP       NOT NULL    DEFAULT NOW(),
    "deleted_at"    TIMESTAMP       NULL,
    "lesson_uuid"   UUID            NOT NULL,
    "name"		    VARCHAR(128)	NOT NULL	UNIQUE,
    "description"	TEXT	    	NOT NULL,
    "order_index"	SMALLINT		NOT NULL,
    "reward"	    SMALLINT		NOT NULL,
    "type"          exercise_type   NOT NULL,
    "data"		    JSONB			NOT NULL,
    CONSTRAINT fk_exercises_lesson FOREIGN KEY ("lesson_uuid") REFERENCES "lessons"("uuid")
);

CREATE          INDEX idx_exercises_deleted_at  ON "exercises" ("deleted_at");
CREATE          INDEX idx_exercises_created_at  ON "exercises" ("created_at");
CREATE          INDEX idx_exercises_lesson_uuid ON "exercises" ("lesson_uuid");
CREATE          INDEX idx_exercises_order_index ON "exercises" ("order_index");