FROM golang:1.18.2-bullseye as builder

RUN apt-get update -y \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get upgrade -y \
    && apt-get install make git zlib1g-dev libssl-dev gperf php-cli cmake clang libc++-dev libc++abi-dev -y
RUN cd ~ \
    && git clone https://github.com/tdlib/td.git \
    && cd td \
    && rm -rf build \
    && mkdir build \
    && cd build \
    && CXXFLAGS="-stdlib=libc++" CC=/usr/bin/clang CXX=/usr/bin/clang++ cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. \
    && cmake --build . -j 8 --target install \
    && cd .. \
    && cd .. \
    && ls -l /usr/local

RUN apt-get update -y \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get upgrade -y \
    && apt-get install sqlite3 -y
# RUN apt-get update -y \
#     && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get upgrade -y \
#     && apt-get install gcc libx32stdc++6-amd64-cross libx32stdc++6-i386-cross libstdc++-9-dev-amd64-cross -y

# COPY . /opt/tgbot
# WORKDIR /opt/tgbot
# # RUN mkdir -p ./bin && \
# #         CGO_ENABLED=1 GOOS=linux GOARCH=amd64 && \
# #         go build -o bin/tgbot .;
# RUN make build_for_docker
COPY ${PWD} /app
WORKDIR /app

# Toggle CGO based on your app requirement. CGO_ENABLED=1 for enabling CGO
# RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin .
# Use below if using vendor
# RUN CGO_ENABLED=1 go build -mod=vendor -ldflags '-s -w -extldflags "-static"' -o /app/appbin .
RUN CGO_ENABLED=1 go build -mod=vendor -ldflags '-s -w' -o /app/appbin .



FROM debian:stable-20220509-slim
LABEL MAINTAINER Author vlad@vegner.org
RUN apt-get update -y \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get upgrade -y
RUN adduser --home "/appuser" --disabled-password appuser \
    --gecos "appuser,-,-,-"
USER appuser
# RUN mkdir -p /usr/local/tgbot/config
# COPY --from=builder /opt/tgbot/bin/tgbot /usr/local/tgbot
# RUN ln -s /usr/local/tgbot/tgbot /usr/bin/tgbot

# WORKDIR /usr/local/tgbot
# RUN chmod +x /usr/local/tgbot/tgbot
COPY --from=builder /app/appbin /home/appuser/app/appbin
WORKDIR /home/appuser/app
STOPSIGNAL SIGTERM
VOLUME /home/appuser/app/config
# VOLUME /usr/local/tgbot/config
# CMD [ "/usr/local/tgbot/tgbot" ]
CMD ["./appbin"]

# https://github.com/bnkamalesh/golang-dockerfile/blob/master/Dockerfile-debian
# FROM golang:1.18 AS builder

# COPY ${PWD} /app
# WORKDIR /app

# # Toggle CGO based on your app requirement. CGO_ENABLED=1 for enabling CGO
# RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go
# # Use below if using vendor
# # RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go

# FROM debian:stable-slim
# LABEL MAINTAINER Author <author@example.com>

# # Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
# RUN apt-get update && apt-get install -y --no-install-recommends \
# 		ca-certificates  \
#         netbase \
#         && rm -rf /var/lib/apt/lists/ \
#         && apt-get autoremove -y && apt-get autoclean -y

# # Add new user 'appuser'. App should be run without root privileges as a security measure
# RUN adduser --home "/appuser" --disabled-password appuser \
#     --gecos "appuser,-,-,-"
# USER appuser

# COPY --from=builder /app /home/appuser/app

# WORKDIR /home/appuser/app

# # Since running as a non-root user, port bindings < 1024 are not possible
# # 8000 for HTTP; 8443 for HTTPS;
# EXPOSE 8000
# EXPOSE 8443

# CMD ["./appbin"]