FROM golang:alpine AS build
WORKDIR /go/src/open_currency

RUN apk update && \
    apk add --no-cache openssl && \
    openssl req -x509 -nodes -days 365 \
   -subj "/C=CL/ST=RM/L=Santiago/O=Open Currency SpA/CN=hinotori.moe" \
   -newkey rsa:2048 -keyout /etc/ssl/private/open_currency.key \
    -out /etc/ssl/certs/open_currency.crt

COPY . .
RUN go build -o /go/bin/open_currency cmd/open_currency/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/ssl/private/open_currency.key /etc/ssl/private/open_currency.key
COPY --from=build /etc/ssl/certs/open_currency.crt /etc/ssl/certs/open_currency.crt
COPY --from=build /go/bin/open_currency /go/bin/open_currency
ENTRYPOINT ["/go/bin/open_currency"]