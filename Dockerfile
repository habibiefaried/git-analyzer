FROM golang:1.17.6-bullseye as builder
RUN apt update && apt install unzip cmake build-essential wget git cmake pkg-config libssl-dev -y && wget https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.zip -O /tmp/v1.3.0.zip && cd /tmp && unzip v1.3.0.zip && cd libgit2-1.3.0 && mkdir build && cd build && cmake .. && cmake --build . --target install && ldconfig
COPY . /app
RUN cd /app && go build -o /main && chmod +x /main && rm -rf /app