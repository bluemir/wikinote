ARG VERSION=dev
FROM fedora:36 as build-env

RUN echo "fastestmirror=1" >> /etc/dnf/dnf.conf
RUN dnf install -y \
    make findutils which \
    golang nodejs \
    && dnf clean all

ENV GOPATH=/root/go
ENV PATH=$PATH:/root/go/bin

# pre build
WORKDIR /pre-build

ADD go.mod go.sum package.json yarn.lock Makefile  ./
ADD scripts/makefile.d/ scripts/makefile.d/

## install build tools
RUN make build-tools

## download dependancy
### go
RUN go mod download
### nodejs
RUN yarn install

# build
WORKDIR /src

## for use vendor folder. uncomment next line
#ENV OPTIONAL_BUILD_ARGS="-mod=vendor"
ENV  OPTIONAL_WEB_BUILD_ARGS="--minify"

ARG VERSION

## copy source
ADD . /src

RUN make build/wikinote

################################################################################
# running image
FROM fedora:36

COPY --from=build-env /src/build/wikinote /bin/wikinote

CMD wikinote server
