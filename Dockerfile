ARG VERSION=dev
FROM fedora:40 as build-env

RUN echo "fastestmirror=1" >> /etc/dnf/dnf.conf
RUN dnf install -y \
    make findutils which \
    golang nodejs musl-gcc \
    && dnf clean all

ENV GOPATH=/root/go
ENV PATH=$PATH:/root/go/bin
ENV CC=musl-gcc

# pre build
WORKDIR /src

## install build tools

ADD Makefile ./
ADD scripts/makefile.d/ scripts/makefile.d/

RUN make build-tools 2>/dev/null

## download dependancy

ADD go.mod go.sum package.json yarn.lock ./
### go
RUN go mod download
### nodejs
RUN yarn install

# build
#WORKDIR /src

ENV OPTIONAL_BUILD_ARGS="-tags embed"
ENV OPTIONAL_WEB_BUILD_ARGS="--minify"

ARG VERSION

## copy source
ADD . /src

RUN make build/wikinote

################################################################################
# running image
FROM alpine:3.18.6

COPY --from=build-env /src/build/wikinote /bin/wikinote

CMD wikinote server
