const logger = require('./logger').logger;
const appFactory = require('./providers/appProvider').factory;
//const web3Factory = require('./providers/web3Provider').factory;

const  oneledgerNetworkService = require('./services/OneledgerNetworkService').instance();

exports.run = (env)=>{
  const config = require(`./configs/${env}`).config;
  const [app, port] = appFactory(config);
  //const contract = web3Factory(config);
  const context = {
    app, port
  }

  oneledgerNetworkService.init(context);
  //register the service

  //end of register
  app.listen(port, ()=> logger.info(`Oneledger Network RPC(${env}) service listening on port ${port}`))
}
