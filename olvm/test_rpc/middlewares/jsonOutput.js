const jsonSend = require('../libs/JSONOutput').jsonSend;

exports.json =  (req, res, next) => {
  res.json =  function(any){
    jsonSend(res, any);
  }
  next();
}
