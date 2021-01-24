FROM golang:1.15 as builder
WORKDIR /
COPY ./* ./
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 go build -o main ./main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder ./main .
RUN mkdir data
ENTRYPOINT [ "./main" ]

