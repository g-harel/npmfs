FROM golang:1.12-alpine AS server

# Required to fetch go modules at build time.
RUN apk add git

WORKDIR /rejstry

COPY . .

RUN go build -o app .

#

FROM alpine:3.9

RUN apk add ca-certificates
RUN apk add git

# Configure git to be able to create commmits.
RUN git config --global user.email "server@rejstry.com"
RUN git config --global user.name "rejstry"

WORKDIR /rejstry

# Copy server binary from first stage.
COPY --from=server /rejstry/app .

# Copy static files from project source.
COPY assets assets
COPY templates templates
RUN rm templates/*.go

ENV ENV production

CMD ./app
