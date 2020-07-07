FROM golang:1.14-alpine AS builder

RUN apk add --no-cache git

RUN set -eux;								\
    export GOPATH="/go";						\
	export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH";		\
	export GOPROXY="http://goproxy.build-tools.rtbrick.net:3000/";	\
	export GONOSUMDB="gitlab.rtbrick.net/*,nd.rtbrick.com/*";	\
	go get -u golang.org/x/lint/golint;				\
	go version;

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOPROXY http://goproxy.build-tools.rtbrick.net:3000/
ENV GONOSUMDB "gitlab.rtbrick.net/*,nd.rtbrick.com/*"

# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy the code into the container
COPY . .

RUN go build -o leitstand-template-engine ./cmd/template-engine/

FROM scratch
COPY --from=builder /build/leitstand-template-engine /main
COPY ./build-artifacts/config.json /etc/rtbrick/leitstand-template-engine/config.json
COPY ./templates/ /etc/rtbrick/leitstand-template-engine/templates/

# Command to run the executable
ENTRYPOINT ["/main"]
