var querystring = require('querystring');
var request = require('request');

  
  request({
    url: "http://vps.colobu.com:9981/",
    method: "POST",
    headers: {
        'Content-Type': 'application/rpcx',
        'X-RPCX-MessageID': '12345678',
        'X-RPCX-MesssageType': '0',
        'X-RPCX-SerializeType': '1',
        'X-RPCX-ServicePath': 'Arith',
        'X-RPCX-ServiceMethod': 'Mul'
    },
    body: '{"A":10, "B":20}'
}, function(error, response, body) {
    if (!error && response.statusCode == 200) {
        console.log(body);
    }
});

 