'use strict';

const money = (minor, currency) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(minor / 100);

const TEMPLATES = {
  transfer_completed: (d) => ({
    subject: `You sent ${money(d.amountMinor, d.currency)}`,
    text: `Your transfer of ${money(d.amountMinor, d.currency)} to ${d.recipient} has completed.`,
  }),
  transfer_review: (d) => ({
    subject: 'Your transfer is being reviewed',
    text: `For your protection, a transfer of ${money(d.amountMinor, d.currency)} needs a quick review. We will update you within one business day.`,
  }),
  invoice_due: (d) => ({
    subject: `Invoice ${d.number} is due ${d.dueDate}`,
    text: `Invoice ${d.number} for ${money(d.totalMinor, d.currency)} is due on ${d.dueDate}.`,
  }),
  card_frozen: (d) => ({
    subject: 'Your card has been frozen',
    text: `Your card ending in ${d.last4} was frozen. Unfreeze it any time from the Cards page.`,
  }),
};

function render(name, data) {
  const tpl = TEMPLATES[name];
  if (!tpl) throw new Error(`unknown template: ${name}`);
  return tpl(data);
}

module.exports = { render, TEMPLATES };
