FROM debian:13

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y \
    git \
    curl \
    ca-certificates \
    build-essential \
    cmake \
    protobuf-compiler \
    protobuf-c-compiler \
    protobuf-compiler-grpc \
    libprotobuf-dev \
    libprotobuf-c-dev \
    libgrpc++-dev \
    libgrpc-dev \
    libyang-dev \
    libjson-c-dev \
    libelf-dev \
    libreadline-dev \
    libcap-dev \
    python3-dev \
    python3-pip \
    bison \
    flex \
    autoconf \
    automake \
    libtool \
    pkg-config \
    iproute2 \
    iputils-ping \
    net-tools \
    tcpdump \
    vim \
    less \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r frr && useradd -r -g frr frr

WORKDIR /src/frr

# Копируем содержимое локальной папки в рабочую директорию контейнера
COPY third_party/frr/ .

RUN ./bootstrap.sh

RUN ./configure \
    --enable-grpc \
    --enable-mgmtd \
    --enable-pimd \
    --enable-multipath=64 \
    --enable-user=frr \
    --enable-group=frr \
    --enable-vty-group=frr \
    --prefix=/usr \
    --sysconfdir=/etc \
    --localstatedir=/var \
    --sbindir=/usr/lib/frr

RUN make -j$(nproc)

RUN make install

RUN ldconfig

#
# DIRECTORIES
#

RUN mkdir -p /etc/frr
RUN mkdir -p /var/run/frr
RUN mkdir -p /var/log/frr

#
# LOG FILES
#

RUN touch /var/log/frr/frr.log && \
    chown -R frr:frr /var/log/frr && \
    chmod -R 775 /var/log/frr

#
# PERMISSIONS
#

RUN chown -R frr:frr /etc/frr

#
# PORTS
#

EXPOSE 50051
EXPOSE 50052
EXPOSE 50053
EXPOSE 50054
EXPOSE 50055
EXPOSE 50056
EXPOSE 50057
EXPOSE 50058
EXPOSE 50059
EXPOSE 50060
EXPOSE 50061
EXPOSE 50062
EXPOSE 50063
EXPOSE 50064

#
# STARTUP
#

CMD ["/bin/sh", "-c", "\
set -e && \
\
echo '=====================================' && \
echo 'STARTING FRR...' && \
echo '=====================================' && \
\
chown -R frr:frr /var/run/frr && \
chmod -R 775 /var/run/frr && \
\
/usr/lib/frr/frrinit.sh start && \
\
sleep 5 && \
\
echo '' && \
echo '=====================================' && \
echo 'FRR STARTED SUCCESSFULLY' && \
echo '=====================================' && \
echo '' && \
\
echo 'ACTIVE DAEMONS:' && \
ps aux | grep frr && \
echo '' && \
\
echo 'AVAILABLE LOGS:' && \
ls -lah /var/log/frr && \
echo '' && \
\
tail -f /var/log/frr/frr.log \
"]
