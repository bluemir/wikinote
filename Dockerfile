ARG VERSION=dev
FROM fedora:33 as build-env

RUN echo "fastestmirror=1" >> /etc/dnf/dnf.conf
RUN dnf install -y \
    make findutils which \
    golang nodejs \
    protobuf protobuf-compiler protobuf-devel \
    && dnf clean all

ENV GOPATH=/root/go
ENV PATH=$PATH:/root/go/bin

# pre build
WORKDIR /pre-build

ADD go.mod go.sum package.json yarn.lock Makefile.d/tools.mk ./

## install build tools
RUN make -f tools.mk tools

## download dependancy
### go
RUN go mod download
### nodejs
RUN yarn install

# build
WORKDIR /src

## for use vendor folder. uncomment next line
#ENV OPTIONAL_BUILD_ARGS="-mod=vendor"

ARG VERSION

## copy source
ADD . /src

ARG APP_NAME
RUN make build/$APP_NAME

################################################################################
# running image
FROM fedora:33

WORKDIR /
ARG APP_NAME
ENV APP_NAME $APP_NAME
COPY --from=build-env /src/build/$APP_NAME /bin/$APP_NAME

CMD $APP_NAME

