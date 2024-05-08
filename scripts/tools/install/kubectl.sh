set -e

BINDIR=$(pwd)/build/tools

OS=$(go env GOOS)
ARCH=$(go env GOARCH)

TMP_DIR=$(mktemp -d)
cd $TMP_DIR

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/$OS/$ARCH/kubectl"
chmod +x kubectl

mkdir -p $BINDIR
mv kubectl $BINDIR

rm -rf $TMP_DIR
