FROM golang:1.14-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/bubbles
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN go build github.com/agilestacks/bubbles/cmd/bubbles


FROM alpine:3.11

WORKDIR /app
COPY --from=builder /go/src/bubbles/bubbles ./

EXPOSE 8005
ENV BUBBLES_API_SECRET ""

ENTRYPOINT ["./bubbles"]
CMD ["-http_port", "8005"]
