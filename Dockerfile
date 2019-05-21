FROM golang:1.12-alpine AS server

WORKDIR /rejstry

COPY . .

RUN go build -o app ./server

#

FROM alpine:3.9

RUN apk add ca-certificates
RUN apk add git

RUN git config --global user.email "server@rejstry.com"
RUN git config --global user.name "rejstry"

COPY --from=server /rejstry/app .

CMD ./app
