const config = require('./common').config;
const redis = null;

exports.config = {
  ...config,
  redis,
  cookie_security: false,
  secret: "localisalwaysdangerousplace"
}
