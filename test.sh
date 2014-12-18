#!/bin/bash -e
ITEM_ID=$RANDOM

echo "CLEAN INDEX"
http -v DELETE localhost:8080/20141217/index/disk/demo/index_$ITEM_ID
echo "CLEAN ITEM"
http -v DELETE localhost:8080/20141217/item/disk/demo/$ITEM_ID

echo "CREATE INDEX user_id_${ITEM_ID}"
http -v --json PUT localhost:8080/20141217/index/disk/demo/index_$ITEM_ID \
    key=user_id_${ITEM_ID} type=STRING

echo "CREATE ITEM"
http -v --json PUT localhost:8080/20141217/item/disk/demo/$ITEM_ID \
    meta:='{"ttl": 3600}'\
    bins:='{"'user_id_${ITEM_ID}'": "'$ITEM_ID'"}'

echo "FETCH BY PK"
http -v GET localhost:8080/20141217/item/disk/demo/$ITEM_ID

echo "FETCH BY user_id_${ITEM_ID}"
http -v GET localhost:8080/20141217/query/disk/demo/?user_id_${ITEM_ID}=$ITEM_ID


echo "MULTIFETCH BY PK"
http -v GET localhost:8080/20141217/item/disk/demo/1,2,3,4

echo "TEST INDEX"
http -v GET localhost:8080/20141217/index/disk/demo/index_$ITEM_ID

echo "CLEAN ITEM"
http -v DELETE localhost:8080/20141217/item/disk/demo/$ITEM_ID

echo "CLEAN INDEX"
http -v DELETE localhost:8080/20141217/index/disk/demo/index_$ITEM_ID
