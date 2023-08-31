-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";

-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

-- Create "users" table
CREATE TABLE "public"."users" (
    "id" serial NOT NULL,
    "avatar" character varying(255) NOT NULL,
    "username" character varying(14) NOT NULL,
    "email" character varying(255) NOT NULL,
    "password" character varying(128) NOT NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Create index "email_unique" to table: "users"
CREATE UNIQUE INDEX "email_unique" ON "public"."users" ("email");

-- Create index "username_unique" to table: "users"
CREATE UNIQUE INDEX "username_unique" ON "public"."users" ("username");

-- Create "posts" table
CREATE TABLE "public"."posts" (
    "id" serial NOT NULL,
    "image" character varying(255) NOT NULL,
    "description" character varying(255) NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "user_id" integer NOT NULL DEFAULT 0,
    PRIMARY KEY ("id"),
    CONSTRAINT "user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create "comments" table
CREATE TABLE "public"."comments" (
    "id" serial NOT NULL,
    "post_id" integer NOT NULL,
    "comment_id" integer NULL,
    "user_id" integer NOT NULL,
    "body" text NOT NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "comment_id_fk" FOREIGN KEY ("comment_id") REFERENCES "public"."comments" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "post_id_fk" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create "comment_likes" table
CREATE TABLE "public"."comment_likes" (
    "id" serial NOT NULL,
    "user_id" integer NOT NULL,
    "comment_id" integer NOT NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "comment_id_fk" FOREIGN KEY ("comment_id") REFERENCES "public"."comments" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create index "comment_likes_user_id_comment_id_key" to table: "comment_likes"
CREATE UNIQUE INDEX "comment_likes_user_id_comment_id_key" ON "public"."comment_likes" ("user_id", "comment_id");

-- Create "follows" table
CREATE TABLE "public"."follows" (
    "id" serial NOT NULL,
    "follower_id" integer NOT NULL,
    "followee_id" integer NOT NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "followee_id_fk" FOREIGN KEY ("followee_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "follower_id_fk" FOREIGN KEY ("follower_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create index "follows_follower_id_followee_id_key" to table: "follows"
CREATE UNIQUE INDEX "follows_follower_id_followee_id_key" ON "public"."follows" ("follower_id", "followee_id");

-- Create "post_likes" table
CREATE TABLE "public"."post_likes" (
    "id" serial NOT NULL,
    "user_id" integer NOT NULL,
    "post_id" integer NOT NULL,
    "created_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "post_id_fk" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create index "post_likes_user_id_post_id_key" to table: "post_likes"
CREATE UNIQUE INDEX "post_likes_user_id_post_id_key" ON "public"."post_likes" ("user_id", "post_id");
