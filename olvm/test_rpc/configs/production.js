const config = require('./common').config;
const redis = {
  host: '159.65.215.181',
  port: 6379,
  ttl: 3600, // 1 hour
  pass: 'eXQSFQMLSnvQsjF4BiCN'
}

exports.config = {
  ...config,
  redis,
  cookie_security: true,
  secret: process.env.NODE_PASSWORD_SECRET
}
