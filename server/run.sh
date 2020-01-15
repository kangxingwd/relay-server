#!/bin/bash

service bind9 stop
service bind9 start

sleep 2

cd /opt/remote_server/server/bin

./awesome-raserver


