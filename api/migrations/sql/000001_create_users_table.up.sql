CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" varchar(100) NOT NULL CHECK (username <> ''),
  "email" varchar(255) UNIQUE NOT NULL CHECK (email <> ''),
  "password" varchar(255) NOT NULL CHECK (password <> ''),
  "created_at" timestamp DEFAULT NULL,
  "updated_at" timestamp DEFAULT NULL,
  "deleted_at" timestamp DEFAULT NULL
);