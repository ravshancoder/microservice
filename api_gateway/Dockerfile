FROM golang:1.20-alpine
RUN mkdir api
COPY . /api
WORKDIR /api
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8030