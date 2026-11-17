'use strict';
const { consume } = require('./queue');
const { sendEmail } = require('./mailer');
const { render } = require('./templates');
const { deliverWebhook } = require('./webhooks');

const handlers = {
  async EMAIL({ to, template, data }) {
    const { subject, text } = render(template, data);
    await sendEmail({ to, subject, text });
  },
  async WEBHOOK({ url, secret, event }) {
    await deliverWebhook(url, secret, event);
  },
};

async function handle(job) {
  const fn = handlers[job.channel];
  if (!fn) throw new Error(`unsupported channel ${job.channel}`);
  await fn(job);
}

consume('notifications', handle).catch((err) => {
  console.error('worker crashed', err);
  process.exit(1);
});
