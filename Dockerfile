FROM --platform=linux/amd64 golang:1.21.5-alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./server ./src/main.go

FROM --platform=linux/amd64 alpine
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

ENV PORT 3000
EXPOSE 3000

ENV PATH_TO_FFMPEG_EXECUTABLE=ffmpeg

COPY --from=0 /app/server /usr/bin/server
CMD ["server"]
