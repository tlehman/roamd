#!/bin/sh

curl -H 'Content-Type: application/json' http://localhost:8080 -d "{\"note\":\"test for real\"}"
