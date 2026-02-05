FROM golang:1.25-alpine3.23 as build
WORKDIR /src
RUN apk add --no-cache make
COPY cmd /src/cmd
COPY internal /src/internal
COPY go.* /src/
COPY Makefile /src/
RUN make test build

FROM alpine:3.23.2
RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY --from=build /src/build/robinhood-export /bin/robinhood-export
ENTRYPOINT ["/bin/robinhood-export"]
