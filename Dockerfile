FROM golang:alpine as builder

WORKDIR /app


COPY . .

RUN ls
ENV CGO_ENABLED=0
RUN go build --mod vendor -o  go-rest-service .

FROM golang:alpine
COPY --from=builder /app/go-rest-service /go-rest-service
CMD ["./go-rest-service"]