#!/bin/bash

./blockChain send -f 张三 -t 李四 -a 10 -m 班长 -d "张三转李四10 bitcoin"
./blockChain send -f 张三 -t 王五 -a 20 -m 班长 -d "张三转王五20 bitcoin"
./blockChain send -f 王五 -t 李四 -a 2 -m 班长 -d "王五转李四2 bitcoin"
./blockChain send -f 王五 -t 李四 -a 3 -m 班长 -d "王五转李四3 bitcoin"
./blockChain send -f 王五 -t 张三 -a 5 -m 班长 -d "王五转张三10 bitcoin"
./blockChain send -f 李四 -t 赵六 -a 14 -m 班长 -d "李四转赵六14 bitcoin"
