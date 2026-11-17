-- AlterTable
ALTER TABLE "accounts" ADD COLUMN "closedAt" TIMESTAMP(3);

-- Backfill accounts that were closed before the column existed.
UPDATE "accounts" SET "closedAt" = "updatedAt" WHERE "status" = 'CLOSED' AND "closedAt" IS NULL;

ALTER TABLE "api_keys" ADD COLUMN "label" TEXT NOT NULL DEFAULT '';
