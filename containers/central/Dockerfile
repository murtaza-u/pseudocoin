FROM golang:alpine as builder

RUN apk add --update --no-cache ca-certificates git

COPY . ./pseudocoin

RUN cd pseudocoin && go build -o pcli

ENTRYPOINT ["/go/pseudocoin/pcli"]

FROM alpine:latest

COPY --from=builder /go/pseudocoin/pcli pcli

RUN ./pcli blockchain --create blockchain.db --address ""

ENTRYPOINT [ "./pcli", "centralnode ", "--start", "blockchain.db" ]
