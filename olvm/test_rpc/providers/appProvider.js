const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const express = require('express');
const logger = require('../logger').logger;
const sessionProvider = require('./sessionProvider');

exports.factory = (config) => {
  const app = express();
  app.set('trust proxy', 1);
  app.use(sessionProvider.factory(config));
  app.use(cookieParser())
  app.use(bodyParser.urlencoded({
    extended: false
  }));
  return [app, config.port];
}
