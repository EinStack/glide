# syntax=docker/dockerfile:1
FROM golang:1.22.4-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=linux

WORKDIR /build

COPY . /build/
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg/version.Version=$VERSION -X glide/pkg/version.commitSha=$COMMIT -X glide/pkg/version.buildDate=$BUILD_DATE" -o /build/dist/glide

FROM gcr.io/distroless/static-debian12:nonroot as release

WORKDIR /bin
COPY --from=build /build/dist/glide /bin/

ENTRYPOINT ["/bin/glide"]
