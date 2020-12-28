FROM golang:alpine

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main ./cmd/polaris

EXPOSE 8000

CMD ["/build/main"]