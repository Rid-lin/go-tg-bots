FROM golang:alpine3.15 as builder
# RUN apk add make
# RUN apk add make git gcc musl-dev
COPY . /opt/gomtc
WORKDIR /opt/gomtc
RUN mkdir -p ./bin;\
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64;\
        go build -o bin/gomtc .;
# RUN ["make","build_alpine"]

FROM alpine:3.15.0
# RUN apk update && apk upgrade
RUN mkdir -p /usr/local/gomtc
COPY --from=builder /opt/gomtc/bin/gomtc /usr/local/gomtc
RUN ln -s /usr/local/gomtc/gomtc /usr/bin/gomtc
STOPSIGNAL SIGTERM
EXPOSE 3034
EXPOSE 3032
WORKDIR /usr/local/gomtc
RUN chmod +x /usr/local/gomtc/gomtc
CMD [ "/usr/local/gomtc/gomtc" ]