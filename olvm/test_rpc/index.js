const server = require('./server');

if(require && require.main === module){
  const env = process.env.NODE_ENV || 'local';
  server.run(env);
}
