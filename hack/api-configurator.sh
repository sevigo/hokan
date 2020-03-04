#!/usr/bin/env bash

# type Directory struct {
# 	Active      bool
# 	Path        string
# 	Recursive   bool
# 	Machine     string
# 	IgnoreFiles []string
# 	Targets     []string
# }

curl -v -d '{"active":true,"path":"/home/igor/MyLocalFiles","recursive": true,"machine":"tokyo","targets":["void"]}' localhost:8081/api/directories
