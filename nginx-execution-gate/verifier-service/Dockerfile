FROM golang:1.21-alpine AS build
WORKDIR /src
COPY . /src
WORKDIR /src/integrations/nginx-execution-gate/verifier-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/verifier-service .

FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=build /out/verifier-service /verifier-service
EXPOSE 8081
ENTRYPOINT ["/verifier-service"]
