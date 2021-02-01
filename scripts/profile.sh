#! /usr/bin/env bash

curl -s http://localhost:8081/debug/pprof/profile -o profile
echo "PROFILE DONE"
