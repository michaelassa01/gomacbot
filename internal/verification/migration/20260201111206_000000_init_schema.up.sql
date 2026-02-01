-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-02-01T11:12:50.106Z

CREATE TABLE "verification" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "code" varchar NOT NULL,
  "type" varchar NOT NULL,
  "identifier" varchar NOT NULL,
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

COMMENT ON COLUMN "verification"."code" IS 'this stroes the generated code for the email or phone number verification';

COMMENT ON COLUMN "verification"."type" IS 'the type of verification e.g ''email'' or ''phone_number''';

COMMENT ON COLUMN "verification"."identifier" IS 'the verification identifier can be either the email address or phone number';

COMMENT ON COLUMN "verification"."expired_at" IS 'this indicates when the given token is expiring';

COMMENT ON COLUMN "verification"."created_at" IS 'verification created date';
