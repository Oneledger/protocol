const config = require('./common').config;
const redis = {
  host: 'localhost',
  port: 6379,
  ttl: 3600, // 1 hour
}

exports.config = {
  ...config,
  redis,
  cookie_security: false,
  secret: "somethingisnotassamesawhatwillbe"
}
