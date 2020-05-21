#!/usr/bin/env python
# -*- coding: iso-8859-1 -*-
import json
import requests


url = 'http://vps.colobu.com:9981/'
payload = {
    'A': 10,
    'B': 20
}
# Adding empty header as parameters are being sent in payload
headers = {
    "Host": "vps.colobu.com",
    "Connection": "keep-alive",
    "X-RPCX-MessageID": "12345678",
    "X-RPCX-MesssageType": "0",
    "X-RPCX-SerializeType": "1",
    "X-RPCX-ServicePath": "Arith",
    "X-RPCX-ServiceMethod": "Mul",
}
r = requests.post(url, data=json.dumps(payload), headers=headers)
print(r.content)