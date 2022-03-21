CREATE TABLE "procedures" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar(100),
  "content" varchar(1000),
  "user_id" int,
  "created_at" timestamp DEFAULT NULL,
  "updated_at" timestamp DEFAULT NULL,
  "deleted_at" timestamp DEFAULT NULL
);

ALTER TABLE
  "procedures"
ADD
  FOREIGN KEY ("user_id") REFERENCES "users" ("id");
