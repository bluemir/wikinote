set -e

VERSION=$1
BINDIR=$(pwd)/build/tools

OS=$(go env GOOS)
ARCH=$(go env GOARCH)

TMP_DIR=$(mktemp -d)
cd $TMP_DIR


curl -LO "https://github.com/etcd-io/etcd/releases/download/$VERSION/etcd-$VERSION-$OS-$ARCH.tar.gz"
tar -xvf "etcd-$VERSION-$OS-$ARCH.tar.gz"


mkdir -p $BINDIR
mv etcd-$VERSION-$OS-$ARCH/etcd etcd-$VERSION-$OS-$ARCH/etcdctl $BINDIR

rm -rf $TMP_DIR
