FROM golang:1.23-bookworm AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH=amd64 \
    GOOS=linux

WORKDIR /builder

# Download dependencies using go mod.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the code into the image.
COPY . .

# Build the executable.
RUN go build -ldflags="-w -s" -o run


FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /runner

EXPOSE 8080

COPY --from=builder /builder/run /builder/api.yaml ./

ARG RELEASE_VERSION
ENV RELEASE_VERSION=$RELEASE_VERSION
ENV GIN_MODE=release

ENTRYPOINT ["./run"]
