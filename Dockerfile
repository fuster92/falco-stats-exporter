FROM golang:1.20 as builder
WORKDIR /go/src/falco-stats-exporter
COPY . .
RUN CGO_ENABLED=0 go build -o exporter github.com/fuster92/falco-stats-exporter/


FROM scratch
COPY --from=builder /go/src/falco-stats-exporter/exporter .
ENTRYPOINT ["/exporter"]
