#!/bin/bash

set -e

for i in {1..500}
do
    curl -X POST http://localhost:8080/v1/posts -d "{\"id\": $i, \"body\": \"foobar body\"}"
done
