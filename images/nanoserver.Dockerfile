# syntax=docker/dockerfile:1
FROM golang:1.22-nanoserver-1809 as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=windows

WORKDIR /build

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

COPY . /build/
RUN pwd; Get-ChildItem
RUN go mod download
RUN go build -ldflags "-s -w -X glide/pkg/version.Version=$VERSION -X glide/pkg/version.commitSha=$COMMIT -X glide/pkg/version.buildDate=$BUILD_DATE" -o /build/dist/glide.exe

FROM mcr.microsoft.com/windows/nanoserver:1809 as release

WORKDIR /bin
COPY --from=build /build/dist/glide.exe /bin/

ENTRYPOINT ["/bin/glide.exe"]
