FROM docker.io/library/ubuntu:noble-20240801
# ENV DEBIAN_FRONTEND=noninteractive # 环境变量
# ENV LANG=C.UTF-8  # 环境变量
# RUN sed -i 's@http://archive.ubuntu.com/ubuntu/@http://mirrors.aliyun.com/ubuntu/@g' /etc/apt/sources.list # 更换源
RUN sed -i 's@http://archive.ubuntu.com/ubuntu/@http://mirrors.aliyun.com/ubuntu/@g' /etc/apt/sources.list.d/ubuntu.sources # 更换源
RUN apt-get update -qq && apt-get install -y -qq autopoint autoconf automake libtool-bin gettext libncursesw5-dev dejagnu libnuma-dev libsystemd-dev
RUN apt-get install -y pkg-config
RUN apt-get install -y make
RUN apt-get install -y flex bison sqlite3
RUN apt-get install -y libsqlite3-dev

RUN mkdir -p /util-linux
WORKDIR /util-linux
COPY ./ ./
RUN ls -al && pwd


# RUN  pkg-config --version

RUN bash ./autogen.sh
RUN bash ./configure --disable-dependency-tracking
RUN make
# RUN bash make

