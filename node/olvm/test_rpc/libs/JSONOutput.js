exports.jsonSend = (response,any) => {
  response.setHeader('Content-Type', 'application/json');
  response.send(JSON.stringify(any));
}

exports.jsonEnd = (response,any) => {
  response.setHeader('Content-Type', 'application/json');
  response.end(JSON.stringify(any));
}
