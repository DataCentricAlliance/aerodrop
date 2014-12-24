#!/bin/bash -e
ITEM_ID=$RANDOM

NAMESPACE="test"

echo "Item ID = $ITEM_ID"
echo "CLEAN INDEX"
curl -X DELETE localhost:8080/20141217/index/$NAMESPACE/demo/index_$ITEM_ID
echo "CLEAN ITEM"
curl -X DELETE localhost:8080/20141217/item/$NAMESPACE/demo/$ITEM_ID

echo "CREATE INDEX user_id_${ITEM_ID}"
curl -X PUT localhost:8080/20141217/index/$NAMESPACE/demo/index_$ITEM_ID \
     -H "Content-Type: application/json" \
     -d '{"key": "user_id_${ITEM_ID}", "type": "STRING"}'

echo "CREATE ITEM"
curl -X PUT localhost:8080/20141217/item/$NAMESPACE/demo/$ITEM_ID \
     -H "Content-Type: application/json" \
     -d '{"meta": {"ttl": 3600}, "bins": {"user_id_$ITEM_ID": "$ITEM_ID"}}'


sleep 1
echo "FETCH BY PK"
curl -X GET localhost:8080/20141217/item/$NAMESPACE/demo/$ITEM_ID

echo "FETCH BY user_id_${ITEM_ID}"
curl -X GET localhost:8080/20141217/query/$NAMESPACE/demo/?user_id_${ITEM_ID}=$ITEM_ID


sleep 1
echo "MULTIFETCH BY PK"
curl -X GET localhost:8080/20141217/item/$NAMESPACE/demo/1,2,3,4,$ITEM_ID

echo "CLEAN ITEM"
curl -X DELETE localhost:8080/20141217/item/$NAMESPACE/demo/$ITEM_ID

echo "CLEAN INDEX"
curl -X DELETE localhost:8080/20141217/index/$NAMESPACE/demo/index_$ITEM_ID
