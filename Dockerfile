# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22

FROM golang:${GO_VERSION} AS build

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGET_OS=linux

# Build the application.
COPY . /app
WORKDIR /app

RUN go mod download all \
    && CGO_ENABLED=0 GOOS=${TARGET_OS} go build -o ./backend api/cmd/main.go

FROM alpine:latest AS app

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add ca-certificates tzdata \
    && update-ca-certificates

# Create a non-privileged user that the app will run under
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    runner
USER runner

# Copy the executable from the "build" stage.
COPY --from=build /app/backend /backend

# Expose the port that the application listens on.
EXPOSE 8080

CMD ["/backend", "--port=8080", "--dev=false"]
