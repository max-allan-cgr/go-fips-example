FROM cgr.dev/chainguard-private/go-fips AS base

WORKDIR /app
COPY main.go /app
RUN go build main.go

FROM cgr.dev/chainguard-private/chainguard-base-fips:latest
COPY --from=base /app/main /

ENTRYPOINT ["/main"]
