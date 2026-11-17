'use strict';
const nodemailer = require('nodemailer');

const transport = nodemailer.createTransport({
  host: process.env.SMTP_HOST || 'localhost',
  port: Number(process.env.SMTP_PORT || 1025),
  secure: false,
});

async function sendEmail({ to, subject, text }) {
  return transport.sendMail({
    from: process.env.SMTP_FROM || 'no-reply@finflow.example',
    to,
    subject,
    text,
  });
}

module.exports = { sendEmail };
