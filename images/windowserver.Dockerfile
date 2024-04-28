# syntax=docker/dockerfile:1
FROM golang:1.22-windowsservercore-1809 as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

WORKDIR /build

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

COPY . /build/
RUN go mod download
RUN GOOS=windows go build -v -o /build/dist/glide.exe -ldflags "-s -w -X glide/pkg/version.Version="$VERSION" -X glide/pkg/version.commitSha="$COMMIT" -X glide/pkg/version.buildDate="$BUILD_DATE""

FROM mcr.microsoft.com/windows/servercore:1809 as release

WORKDIR /bin
COPY --from=build /build/dist/glide.exe /bin/

ENTRYPOINT ["/bin/glide.exe"]
