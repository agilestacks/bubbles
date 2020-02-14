FROM golang:1.13-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/bubbles
COPY src/bubbles /go/src/bubbles
RUN go get github.com/kardianos/govendor
RUN /go/bin/govendor sync
RUN go build -v bubbles


FROM alpine:3.11

COPY --from=builder /go/src/bubbles/bubbles /bin/bubbles
RUN chmod +x /bin/bubbles

EXPOSE 8005
ENV BUBBLES_API_SECRET ""

CMD ["bubbles", "-http_port", "8005"]
