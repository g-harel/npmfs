FROM golang:1.12-alpine AS build

# Required to fetch go modules at build time.
RUN apk add git

WORKDIR /npmfs

COPY . .

RUN go build -o website .

#

FROM alpine:3.9

RUN apk add ca-certificates
RUN apk add git

WORKDIR /npmfs

# Copy server binary from first stage.
COPY --from=build /npmfs/website .

# Copy static files from project source.
COPY assets assets
COPY templates templates

CMD ./website
