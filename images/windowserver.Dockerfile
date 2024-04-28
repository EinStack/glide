# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=windows

WORKDIR /build

COPY . /build/
RUN go mod download
RUN go build -v -o /build/dist/glide.exe -ldflags "-s -w -X glide/pkg/version.Version="$VERSION" -X glide/pkg/version.commitSha="$COMMIT" -X glide/pkg/version.buildDate="$BUILD_DATE""

FROM mcr.microsoft.com/windows/servercore:1809 as release

WORKDIR /bin

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

COPY --from=build /build/dist/glide.exe /bin/

ENTRYPOINT ["/bin/glide.exe"]
