FROM golang:1.22.2-alpine as build

ENV CGO_ENABLED=0
ARG VCS_REF

RUN mkdir -p /var/ultimate-service-v1
WORKDIR /var/ultimate-service-v1
COPY . .
RUN go mod download

WORKDIR /var/ultimate-service-v1/app/sales-api
RUN go build -ldflags "-X main.build=${VCS_REF}"

FROM alpine:3.12
ARG BUILD_DATE
ARG VCS_REF

RUN mkdir -p /var/ultimate-service-v1

COPY --from=build /var/ultimate-service-v1/private.pem /var/ultimate-service-v1/private.pem
COPY --from=build /var/ultimate-service-v1/app/sales-api/sales-api /var/ultimate-service-v1/sales-api

WORKDIR /var/ultimate-service-v1

CMD ["./sales-api"]
