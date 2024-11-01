FROM golang:alpine as go
WORKDIR /deltagen
COPY . .
RUN CGO_ENABLED=0 go build .

FROM alpine
COPY --from=go /deltagen/deltagen /usr/bin
ENTRYPOINT ["/usr/bin/deltagen"]
