FROM golang:1.12-alpine AS server

# Required to fetch go modules at build time.
RUN apk add git

WORKDIR /rejstry

COPY . .

RUN go build -o server .

#

FROM alpine:3.9

RUN apk add ca-certificates
RUN apk add git

WORKDIR /rejstry

# Copy server binary from first stage.
COPY --from=server /rejstry/server .

# Copy static files from project source.
COPY assets assets
COPY templates templates
RUN rm templates/*.go

CMD ./server
