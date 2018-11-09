const jsonOutput = require('../middlewares/jsonOutput').json;

class OneledgerNetworkService {
  constructor() {

  }
  init(context) {
    this.context = context;
    if (this.context.app) {
      this.context.app.get('/echo', jsonOutput, this.echo.bind(this));
      this.context.app.get('/blockheight', jsonOutput, this.blockHeight.bind(this));
    }
  }
  echo(req, res) {
    res.json({echo: 'SUCCESS'});
  }
  blockHeight(req, res) {
    res.json({timestamp: new Date().getTime(), datetime: new Date()});
  }
}

exports.instance = () => new OneledgerNetworkService();
