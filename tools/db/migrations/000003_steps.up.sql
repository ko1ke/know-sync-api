CREATE TABLE "steps" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar(100),
  "content" varchar(255),
  "img_name" varchar(255),
  "procedure_id" int,
  "created_at" timestamp DEFAULT NULL,
  "updated_at" timestamp DEFAULT NULL,
  "deleted_at" timestamp DEFAULT NULL
);

ALTER TABLE
  "steps"
ADD
  FOREIGN KEY ("procedure_id") REFERENCES "procedures" ("id");
