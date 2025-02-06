FROM golang:1.22 as build

WORKDIR /src
COPY . .
RUN go build -v -o /app ./cmd/action

FROM gcr.io/distroless/base
COPY --from=build /app /app
CMD ["/app"]
