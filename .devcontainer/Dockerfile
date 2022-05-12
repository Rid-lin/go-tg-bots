# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.234.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

RUN echo "installing go version 1.18.2..." \
    && cd ~ \
    && wget https://go.dev/dl/go1.18.2.linux-amd64.tar.gz  \
    && sudo rm -rf /usr/local/go \
    && sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz

ARG NUMBER_CPU=8
RUN sudo apt-get update -y \
    && sudo apt-get upgrade -y \
    && sudo apt-get install make git zlib1g-dev libssl-dev gperf php-cli cmake clang libc++-dev libc++abi-dev -y \
    && cd ~ \
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