FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/mining-eco .
FROM gcr.io/distroless/base-debian12
WORKDIR /tmp
COPY --from=build /out/mining-eco /app/mining-eco
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/mining-eco"]
