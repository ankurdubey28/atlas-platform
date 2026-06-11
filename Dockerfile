FROM golang:1.25-bookworm@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 AS build

WORKDIR /app

RUN useradd -u 1001 nonroot

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -tags netgo \
    -o api-go \
    ./cmd/api

FROM scratch AS api

LABEL org.opencontainers.image.authors="ankur@github.com/ankurdubey28"

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/api-go /api-go

USER nonroot

EXPOSE 3030

ENTRYPOINT ["/api-go"]
CMD ["--port", "3030"]


FROM scratch AS migrate


COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/goose /goose
COPY --from=build /app/cmd/migrate/migrations /app/cmd/migrate/migrations

USER nonroot

ENTRYPOINT ["/goose"]