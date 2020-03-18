#!/usr/bin/env bash

curl -v -d '{"active":true,"path":"/home/igor/MyLocalFiles","recursive": true,"machine":"tokyo","targets":["void"]}' localhost:8081/api/directories
