FROM golang:alpine as builder

ENV WORKDIR /build

WORKDIR $WORKDIR
COPY . $WORKDIR

ENV FILE app.server

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
RUN #go env -w GO111MODULE=on
RUN #go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -ldflags "-w -s" -o $FILE

FROM alpine
#FROM debian:stretch-slim

ENV WORKDIR /app
ENV FILE app.server

WORKDIR $WORKDIR

COPY --from=builder /build/$FILE $WORKDIR
#ADD .env $WORKDIR/.env

RUN chmod +x $WORKDIR/$FILE

EXPOSE 3000

ENTRYPOINT [ "./app.server" ]
#CMD [ "./app.server" ]