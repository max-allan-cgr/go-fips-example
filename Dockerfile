FROM cgr.dev/chainguard-private/go-fips:v1.26.8.1-r0 AS base

WORKDIR /app
COPY main.go /app
RUN go build main.go

FROM cgr.dev/chainguard-private/chainguard-base-fips:latest
COPY --from=base /app/main /

ENTRYPOINT ["/main"]
