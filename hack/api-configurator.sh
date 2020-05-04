#!/usr/bin/env bash

curl -v -d '{"active":true,"path":"C:\Users\Igor\Files","recursive": true,"machine":"osaka","targets":["void"]}' localhost:8081/api/directories
