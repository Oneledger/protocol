const session = require('express-session');
const RedisStore = require('connect-redis')(session);
const redis   = require('redis');
const logger = require('../logger').logger;

exports.factory =  (config) => {
  if(config.redis === null) {
    logger.info('Non redis session store created');
    return session({
      secret: 'it is not that I want to hide, but that I donnot want to show to you',
      resave: false,
      saveUninitialized: true,
      cookie: {secure: config.cookie_security}
    });
  }else{
    logger.info('Redis session store created at ' + config.redis.host + ':' + config.redis.port);
    const client = redis.createClient(config.redis.port, config.redis.host)
    const redisOptions = {
      client: client,
      ...config.redis
    }
    return session({
      store: new RedisStore(redisOptions),
      secret: 'it is not that I want to hide, but that I donnot want to show to you',
      resave: false,
      saveUninitialized: true,
      cookie: {secure: config.cookie_security}
    });
  }
}
