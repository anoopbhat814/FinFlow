-- Enforce ledger invariants at the database level.
ALTER TABLE "ledger_entries"
    ADD CONSTRAINT "ledger_entries_amountMinor_positive" CHECK ("amountMinor" > 0);

ALTER TABLE "accounts"
    ADD CONSTRAINT "accounts_holdMinor_nonnegative" CHECK ("holdMinor" >= 0);

ALTER TABLE "cards"
    ADD CONSTRAINT "cards_expMonth_range" CHECK ("expMonth" BETWEEN 1 AND 12);

CREATE INDEX "ledger_entries_transactionId_idx" ON "ledger_entries"("transactionId");
