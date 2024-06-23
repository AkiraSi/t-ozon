FROM golang:1.22.1-alpine AS builder

RUN apk add --update postgresql

COPY schema.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 1
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB Ozon

CMD ["postgres", "-D", "/var/lib/postgresql/data", "-c", "fsync=off"]

FROM builder AS app
COPY . .
RUN go mod download
RUN go build -o /app/server .
CMD ["/app/server"]
