set -e

PROTOC_VERSION=21.4
PROTOC_ZIP=protoc-$PROTOC_VERSION-osx-x86_64.zip
curl -sSOL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_ZIP
unzip -q -o $PROTOC_ZIP -d /usr/local bin/protoc
unzip -q -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
echo "✅  protoc"

cd backend && go install $(go list -f '{{ join .Imports "\n" }}' tools/tools.go)
echo "✅  tools.go"