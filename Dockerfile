FROM golang:1.12-alpine AS build

WORKDIR /rejstry

COPY . .

RUN go build -o server ./server

#

FROM alpine:3.9

RUN apk --no-cache add ca-certificates

COPY --from=build /rejstry/server .

CMD ./server
