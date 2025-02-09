## Builder
FROM golang:1.23-alpine as builder
WORKDIR /go/src/ngaymai
COPY . .
RUN go get .
RUN go build -o app.exe .

## Start from the latest golang base image
FROM golang:1.23-alpine
WORKDIR /app
ARG LOG_DIR=/app/tmp
RUN mkdir -p ${LOG_DIR}
ENV LOG_FILE_LOCATION=${LOG_DIR}/console.log

EXPOSE 8000

# Add from source to /app
COPY --from=builder /go/src/ngaymai/app.exe /app
RUN echo > /app/.env

# Run the binary program produced by `go install`
CMD /app/app.exe