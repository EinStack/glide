# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=linux

WORKDIR /build

COPY . /build/
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg/version.version=$VERSION -X glide/pkg/version.commitSha=$COMMIT -X glide/pkg/version.buildDate=$BUILD_DATE" -o /build/dist/glide

FROM redhat/ubi8-micro:8.9 as release

WORKDIR /bin
COPY --from=build /build/dist/glide /bin/

ENTRYPOINT ["/bin/glide"]
