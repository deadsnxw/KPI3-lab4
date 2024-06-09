# Build stage
FROM golang:1.22 as build

WORKDIR /go/src/practice-4
COPY . .

RUN go test ./...
ENV CGO_ENABLED=0
RUN go install ./cmd/...

FROM alpine:latest
WORKDIR /opt/practice-4

COPY entry.sh /opt/practice-4/
RUN dos2unix /opt/practice-4/entry.sh && chmod +x /opt/practice-4/entry.sh

COPY --from=build /go/bin/* /opt/practice-4

RUN ls -l /opt/practice-4
ENTRYPOINT ["/opt/practice-4/entry.sh"]
CMD ["server"]
