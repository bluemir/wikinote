FROM fedora as builder

RUN dnf -y install \
	nodejs \
	golang \
	git \
	findutils \
	make && dnf clean all

COPY makefile /go/src/github.com/bluemir/wikinote/makefile

ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /go/src/github.com/bluemir/wikinote
RUN go get github.com/GeertJohan/go.rice/rice
RUN npm install -g traceur less

COPY . .

RUN rm -rf .GOPATH && make clean wikinote

FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /go/src/github.com/bluemir/wikinote/wikinote /wikinote
RUN ls /wikinote
EXPOSE 80
ENTRYPOINT ["/wikinote"]
CMD ["serve", "--bind", ":80"]
