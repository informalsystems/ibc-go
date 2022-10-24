#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./ibc -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 2 -name '*.proto'); do
      buf generate --template buf.gen.gogo.yaml $file
  done
done


cd ..

# move proto files to the right places
cp -r github.com/cosmos/ibc-go/v5/* ./
rm -rf github.com

go mod tidy