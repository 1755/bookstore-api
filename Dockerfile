FROM golang:1.23-bookworm AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH=amd64 \
    GOOS=linux

WORKDIR /builder

# Download dependencies using go mod.
COPY go.mod go.sum ./
RUN go mod download && go mod verify 
RUN go install -tags 'postgres' -v github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy the code into the image.
COPY . .

# Build the executable.
RUN go build -ldflags="-w -s" -o run


# FROM gcr.io/distroless/base-debian12 AS runner
FROM golang:1.23-bookworm as runner

WORKDIR /runner

EXPOSE 8080

COPY --from=builder /builder/run /builder/api.yaml ./
COPY --from=builder /builder/scripts ./scripts
COPY --from=builder /builder/migrations ./migrations
COPY --from=builder /go/bin/migrate /bin/migrate

ARG RELEASE_VERSION
ENV RELEASE_VERSION=$RELEASE_VERSION
ENV GIN_MODE=release

ENTRYPOINT ["./run"]
