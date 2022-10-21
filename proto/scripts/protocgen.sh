#!/usr/bin/env bash

set -eo pipefail

# protoc_gen_gocosmos() {
#   if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
#     echo -e "\tPlease run this command from somewhere inside the ibc-go folder."
#     return 1
#   fi

#   go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
# }

# protoc_gen_gocosmos


echo "Generating gogo proto code"
cd proto
# proto_dirs=$(find ./ibc -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
# for dir in $proto_dirs; do
#   for file in $(find "${dir}" -maxdepth 2 -name '*.proto'); do
buf --debug generate --template buf.gen.gogo.yaml
#   done
# done

cd ..

go mod tidy

# move proto files to the right places
cp -r github.com/cosmos/ibc-go/v5/* ./
rm -rf github.com
