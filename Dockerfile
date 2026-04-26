FROM golang:1.25-bookworm AS build

WORKDIR /app

RUN useradd -u 1001 nonroot

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkd/mod\
    --mount=type=cache,target=/root/.cache/go-build\
    go mod download

COPY . .

RUN go build \
    -ldflags="-linkmode external -extldflags '-static'" \
    -tags netgo \
    -o api-go \
    ./cmd/api

###
FROM scratch

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/api-go api-go

USER nonroot

EXPOSE 3030

CMD ["/api-go"]


