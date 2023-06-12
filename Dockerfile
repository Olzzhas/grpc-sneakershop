FROM golang:1.19-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o sneakershop ./cmd/api/main.go

RUN chmod +x sneakershop
CMD ["./cmd/api/main.go"]