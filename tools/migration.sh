#!/bin/sh

if [ $# -ne 2 ]; then
  echo "argment is $#" 1>&2
  echo "pass two argments, 1: data source name, 2: up or down" 1>&2
  exit 1
fi

cd "$(dirname "$0")"
./migrate -database $1 -path ./db/migrations $2

exit 0