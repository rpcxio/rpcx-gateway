<?php
$url = 'http://vps.colobu.com:9981/';
$data = '{"A":10, "B":20}';

// use key 'http' even if you send the request to https://...
$options = array(
    'http' => array(
        'header'  => "Content-type: application/rpcx\r\n" .
        "X-RPCX-MessageID: 12345678\r\n" .
        "X-RPCX-MesssageType: 0\r\n" .
        "X-RPCX-SerializeType: 1\r\n" .
        "X-RPCX-ServicePath: Arith\r\n" .
        "X-RPCX-ServiceMethod: Mul\r\n",
        'method'  => 'POST',
        'content' => $data
    )
);
$context  = stream_context_create($options);
$result = file_get_contents($url, false, $context);
if ($result === FALSE) { /* Handle error */ }

var_dump($result);
?>