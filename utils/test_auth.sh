#!/bin/bash

API_URL=http://127.0.0.1:7654/api/v1/iot/
AUTH_HEADER=`cat auth.txt | tr -d '\n'`
curl --verbose -L -H "X-SESSION-KEY: $AUTH_HEADER" --data '{"payload":"$1035,0,0,0,0,0,0,0,675","created_at":"2014-05-23T23:00:00"}' -X POST $API_URL 
