-- CreateEnum
CREATE TYPE "Currency" AS ENUM ('USD', 'EUR', 'GBP', 'CAD', 'JPY');
CREATE TYPE "AccountStatus" AS ENUM ('PENDING', 'ACTIVE', 'FROZEN', 'CLOSED');
CREATE TYPE "AccountKind" AS ENUM ('CHECKING', 'SAVINGS', 'BUSINESS');
CREATE TYPE "TransferStatus" AS ENUM ('PENDING', 'IN_REVIEW', 'COMPLETED', 'FAILED', 'REVERSED');
CREATE TYPE "EntryDirection" AS ENUM ('DEBIT', 'CREDIT');
CREATE TYPE "InvoiceStatus" AS ENUM ('DRAFT', 'SENT', 'PAID', 'OVERDUE', 'VOID');

-- CreateTable
CREATE TABLE "users" (
    "id" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "fullName" TEXT NOT NULL,
    "passwordHash" TEXT NOT NULL,
    "kycVerified" BOOLEAN NOT NULL DEFAULT false,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "accounts" (
    "id" TEXT NOT NULL,
    "ownerId" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "kind" "AccountKind" NOT NULL DEFAULT 'CHECKING',
    "currency" "Currency" NOT NULL,
    "status" "AccountStatus" NOT NULL DEFAULT 'PENDING',
    "balanceMinor" BIGINT NOT NULL DEFAULT 0,
    "holdMinor" BIGINT NOT NULL DEFAULT 0,
    "iban" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    CONSTRAINT "accounts_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "transfers" (
    "id" TEXT NOT NULL,
    "idempotencyKey" TEXT NOT NULL,
    "sourceAccountId" TEXT NOT NULL,
    "destinationAccountId" TEXT NOT NULL,
    "amountMinor" BIGINT NOT NULL,
    "currency" "Currency" NOT NULL,
    "memo" TEXT,
    "status" "TransferStatus" NOT NULL DEFAULT 'PENDING',
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "completedAt" TIMESTAMP(3),
    CONSTRAINT "transfers_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "ledger_transactions" (
    "id" TEXT NOT NULL,
    "reference" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "ledger_transactions_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "ledger_entries" (
    "id" TEXT NOT NULL,
    "transactionId" TEXT NOT NULL,
    "accountId" TEXT NOT NULL,
    "direction" "EntryDirection" NOT NULL,
    "amountMinor" BIGINT NOT NULL,
    "currency" "Currency" NOT NULL,
    "balanceAfterMinor" BIGINT NOT NULL,
    "postedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "ledger_entries_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "ledger_entries_amount_positive" CHECK ("amountMinor" > 0)
);

CREATE TABLE "invoices" (
    "id" TEXT NOT NULL,
    "number" TEXT NOT NULL,
    "issuerId" TEXT NOT NULL,
    "customerName" TEXT NOT NULL,
    "customerEmail" TEXT NOT NULL,
    "currency" "Currency" NOT NULL,
    "status" "InvoiceStatus" NOT NULL DEFAULT 'DRAFT',
    "issuedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "dueAt" TIMESTAMP(3) NOT NULL,
    "paidAt" TIMESTAMP(3),
    CONSTRAINT "invoices_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "invoice_lines" (
    "id" TEXT NOT NULL,
    "invoiceId" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "quantity" INTEGER NOT NULL,
    "unitPriceMinor" BIGINT NOT NULL,
    CONSTRAINT "invoice_lines_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");
CREATE UNIQUE INDEX "accounts_iban_key" ON "accounts"("iban");
CREATE INDEX "accounts_ownerId_idx" ON "accounts"("ownerId");
CREATE UNIQUE INDEX "transfers_idempotencyKey_key" ON "transfers"("idempotencyKey");
CREATE INDEX "transfers_sourceAccountId_createdAt_idx" ON "transfers"("sourceAccountId", "createdAt");
CREATE UNIQUE INDEX "ledger_transactions_reference_key" ON "ledger_transactions"("reference");
CREATE INDEX "ledger_entries_accountId_postedAt_idx" ON "ledger_entries"("accountId", "postedAt");
CREATE UNIQUE INDEX "invoices_number_key" ON "invoices"("number");
CREATE INDEX "invoices_issuerId_status_idx" ON "invoices"("issuerId", "status");

-- AddForeignKey
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_ownerId_fkey" FOREIGN KEY ("ownerId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "transfers" ADD CONSTRAINT "transfers_sourceAccountId_fkey" FOREIGN KEY ("sourceAccountId") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "transfers" ADD CONSTRAINT "transfers_destinationAccountId_fkey" FOREIGN KEY ("destinationAccountId") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "ledger_entries" ADD CONSTRAINT "ledger_entries_transactionId_fkey" FOREIGN KEY ("transactionId") REFERENCES "ledger_transactions"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "ledger_entries" ADD CONSTRAINT "ledger_entries_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "invoices" ADD CONSTRAINT "invoices_issuerId_fkey" FOREIGN KEY ("issuerId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "invoice_lines" ADD CONSTRAINT "invoice_lines_invoiceId_fkey" FOREIGN KEY ("invoiceId") REFERENCES "invoices"("id") ON DELETE CASCADE ON UPDATE CASCADE;
