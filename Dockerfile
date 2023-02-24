FROM golang:1.19 AS builder

WORKDIR /go/src/
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app . 


FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /go/src/app ./
COPY --from=builder /go/src/hosts.json ./
COPY --from=builder /go/src/public ./public/
CMD ["./app"]
