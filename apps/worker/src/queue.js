'use strict';
const Redis = require('ioredis');

const MAX_ATTEMPTS = 5;

/** Backoff in ms: 1s, 4s, 9s... grows quadratically. */
const backoff = (attempt) => attempt * attempt * 1000;

async function consume(stream, handler, { redis = new Redis(process.env.REDIS_URL) } = {}) {
  const group = 'workers';
  await redis.xgroup('CREATE', stream, group, '$', 'MKSTREAM').catch(() => {});
  const consumer = `worker-${process.pid}`;

  for (;;) {
    const res = await redis.xreadgroup('GROUP', group, consumer, 'COUNT', 10, 'BLOCK', 5000, 'STREAMS', stream, '>');
    if (!res) continue;
    for (const [, messages] of res) {
      for (const [id, fields] of messages) {
        const job = JSON.parse(fields[fields.indexOf('job') + 1]);
        try {
          await handler(job);
          await redis.xack(stream, group, id);
        } catch (err) {
          const attempt = (job.attempt || 0) + 1;
          await redis.xack(stream, group, id);
          if (attempt >= MAX_ATTEMPTS) {
            await redis.xadd(`${stream}:dead`, '*', 'job', JSON.stringify(job), 'error', String(err));
          } else {
            setTimeout(() => redis.xadd(stream, '*', 'job', JSON.stringify({ ...job, attempt })), backoff(attempt));
          }
        }
      }
    }
  }
}

module.exports = { consume, backoff };
