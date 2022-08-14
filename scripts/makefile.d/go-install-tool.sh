set -e

export GOBIN=$PWD/build/tools/

TMP_DIR=$(mktemp -d)
cd $TMP_DIR

go mod init temp
go get $1
go install $1

rm -rf $TMP_DIR
