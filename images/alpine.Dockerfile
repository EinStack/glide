# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

WORKDIR /build

COPY . /build/
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg.version=$VERSION -X glide/pkg.commitSha=$COMMIT -X glide/pkg.buildDate=$BUILD_DATE" -o /build/dist/glide

FROM alpine:3.19 as release

WORKDIR /bin
COPY --from=build /build/glide /bin/

ENTRYPOINT ["/bin/glide"]
