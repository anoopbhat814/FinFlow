-- CreateEnum
CREATE TYPE "CardType" AS ENUM ('VIRTUAL', 'PHYSICAL');
CREATE TYPE "CardStatus" AS ENUM ('INACTIVE', 'ACTIVE', 'FROZEN', 'CANCELLED');
CREATE TYPE "FraudAction" AS ENUM ('ALLOW', 'REVIEW', 'BLOCK');
CREATE TYPE "NotificationChannel" AS ENUM ('EMAIL', 'SMS', 'PUSH', 'WEBHOOK');

-- CreateTable
CREATE TABLE "cards" (
    "id" TEXT NOT NULL,
    "accountId" TEXT NOT NULL,
    "type" "CardType" NOT NULL,
    "status" "CardStatus" NOT NULL DEFAULT 'INACTIVE',
    "brand" TEXT NOT NULL,
    "last4" TEXT NOT NULL,
    "expMonth" INTEGER NOT NULL,
    "expYear" INTEGER NOT NULL,
    "monthlyLimitMinor" BIGINT NOT NULL,
    "spentThisMonthMinor" BIGINT NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "cards_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "fraud_checks" (
    "id" TEXT NOT NULL,
    "transferId" TEXT NOT NULL,
    "score" DECIMAL(5,4) NOT NULL,
    "action" "FraudAction" NOT NULL,
    "signals" JSONB NOT NULL,
    "evaluatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fraud_checks_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "notifications" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "channel" "NotificationChannel" NOT NULL,
    "template" TEXT NOT NULL,
    "payload" JSONB NOT NULL,
    "sentAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "notifications_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "api_keys" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "prefix" TEXT NOT NULL,
    "hash" TEXT NOT NULL,
    "lastUsedAt" TIMESTAMP(3),
    "revokedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "api_keys_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "cards_accountId_idx" ON "cards"("accountId");
CREATE UNIQUE INDEX "fraud_checks_transferId_key" ON "fraud_checks"("transferId");
CREATE INDEX "notifications_userId_createdAt_idx" ON "notifications"("userId", "createdAt");
CREATE UNIQUE INDEX "api_keys_prefix_key" ON "api_keys"("prefix");

-- AddForeignKey
ALTER TABLE "cards" ADD CONSTRAINT "cards_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "fraud_checks" ADD CONSTRAINT "fraud_checks_transferId_fkey" FOREIGN KEY ("transferId") REFERENCES "transfers"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "notifications" ADD CONSTRAINT "notifications_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "api_keys" ADD CONSTRAINT "api_keys_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
