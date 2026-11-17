import { PrismaClient, Currency, AccountStatus } from '@prisma/client';
import { createHash } from 'node:crypto';

const prisma = new PrismaClient();

const sha = (s: string) => createHash('sha256').update(s).digest('hex');

async function main() {
  const alice = await prisma.user.upsert({
    where: { email: 'alice@example.com' },
    update: {},
    create: { email: 'alice@example.com', fullName: 'Alice Example', passwordHash: sha('demo-only'), kycVerified: true },
  });
  const bob = await prisma.user.upsert({
    where: { email: 'bob@example.com' },
    update: {},
    create: { email: 'bob@example.com', fullName: 'Bob Example', passwordHash: sha('demo-only'), kycVerified: true },
  });

  const checking = await prisma.account.create({
    data: { ownerId: alice.id, name: 'Main checking', currency: Currency.USD, status: AccountStatus.ACTIVE, balanceMinor: 250_000n, iban: 'GB29NWBK60161331926819' },
  });
  const euro = await prisma.account.create({
    data: { ownerId: bob.id, name: 'Euro account', currency: Currency.EUR, status: AccountStatus.ACTIVE, balanceMinor: 90_000n },
  });

  await prisma.card.create({
    data: { accountId: checking.id, type: 'VIRTUAL', status: 'ACTIVE', brand: 'VISA', last4: '4242', expMonth: 12, expYear: 2028, monthlyLimitMinor: 100_000n },
  });

  await prisma.invoice.create({
    data: {
      number: 'INV-0001',
      issuerId: alice.id,
      customerName: 'Northwind Traders',
      customerEmail: 'ap@northwind.example',
      currency: Currency.USD,
      status: 'SENT',
      dueAt: new Date(Date.now() + 14 * 86_400_000),
      lines: { create: [
        { description: 'Design retainer', quantity: 1, unitPriceMinor: 120_000n },
        { description: 'Extra revisions', quantity: 3, unitPriceMinor: 7_500n },
      ] },
    },
  });

  await prisma.beneficiary.upsert({
    where: { userId_iban: { userId: alice.id, iban: 'DE89370400440532013000' } },
    update: {},
    create: { userId: alice.id, name: 'Bob Example', iban: 'DE89370400440532013000', bic: 'COBADEFFXXX', currency: Currency.EUR, status: 'VERIFIED', verifiedAt: new Date() },
  });

  await prisma.budget.create({
    data: {
      userId: alice.id,
      name: 'Household',
      period: 'MONTHLY',
      currency: Currency.USD,
      startsAt: new Date(Date.UTC(2025, 0, 1)),
      lines: { create: [
        { category: 'groceries', limitMinor: 60_000n, spentMinor: 18_450n },
        { category: 'transport', limitMinor: 15_000n },
      ] },
    },
  });

  await prisma.exchangeRate.createMany({
    data: [
      { base: Currency.USD, quote: Currency.EUR, rate: '0.92140000', observedAt: new Date(Date.UTC(2025, 0, 1)) },
      { base: Currency.EUR, quote: Currency.USD, rate: '1.08530000', observedAt: new Date(Date.UTC(2025, 0, 1)) },
    ],
    skipDuplicates: true,
  });

  await prisma.webhookEndpoint.create({
    data: { userId: alice.id, url: 'https://hooks.example.com/finflow', secret: sha('whsec-demo-only'), events: ['transfer.completed', 'refund.succeeded'] },
  });

  await prisma.kycDocument.create({
    data: { userId: alice.id, type: 'PASSPORT', status: 'APPROVED', storageKey: 'kyc/demo/alice-passport.enc', reviewedAt: new Date() },
  });

  await prisma.auditLog.create({
    data: { actorId: alice.id, action: 'seed.bootstrap', entityType: 'account', entityId: checking.id, metadata: { source: 'prisma/seed.ts' } },
  });

  console.log(`Seeded users ${alice.email}, ${bob.email}; accounts ${checking.id}, ${euro.id}`);
}

main().finally(() => prisma.$disconnect());
