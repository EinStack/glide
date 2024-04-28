# syntax=docker/dockerfile:1
FROM golang:1.22-windowsservercore-1809 as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=windows

WORKDIR /build

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

RUN  $env:VERSION = '$VERSION'; \
    $env:COMMIT = '$COMMIT'; \
    $env:BUILD_DATE = '$BUILD_DATE';

COPY . /build/
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg/version.Version=$env:VERSION -X glide/pkg/version.commitSha=$env:COMMIT -X glide/pkg/version.buildDate=$env:BUILD_DATE" -o /build/dist/glide.exe

FROM mcr.microsoft.com/windows/nanoserver:1809 as release

WORKDIR /bin
COPY --from=build /build/dist/glide.exe /bin/

ENTRYPOINT ["/bin/glide.exe"]
