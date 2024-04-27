# syntax=docker/dockerfile:1
ARG VERSION
ARG COMMIT
ARG BUILD_DATE

FROM golang:1.22-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=linux

WORKDIR /build

COPY . /build/
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg.version=$VERSION -X glide/pkg.commitSha=$COMMIT -X glide/pkg.buildDate=$BUILD_DATE" -o /build/dist/glide

FROM gcr.io/distroless/static-debian12:nonroot as release

WORKDIR /bin
COPY --from=build /build/dist/glide /bin/

ENTRYPOINT ["/bin/glide"]
