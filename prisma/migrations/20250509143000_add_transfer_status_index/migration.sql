-- Speeds up the review queue and the stuck-transfer sweeper.
CREATE INDEX "transfers_status_createdAt_idx" ON "transfers"("status", "createdAt");
CREATE INDEX "transfers_destinationAccountId_createdAt_idx" ON "transfers"("destinationAccountId", "createdAt");

-- Notifications that have not been delivered yet.
CREATE INDEX "notifications_pending_idx" ON "notifications"("createdAt") WHERE "sentAt" IS NULL;
