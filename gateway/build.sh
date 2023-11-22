#!bin/bash
docker build  -t kkqqqqqq/gatewaytest:0.54 .
docker push kkqqqqqq/gatewaytest:0.54
docker pull kkqqqqqq/gatewaytest:0.54
echo "done!"