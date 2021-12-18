FROM ubuntu:18.04

# java runtime policy.
COPY java_policy /etc

# update software source to tsinghua university mirror.
RUN mv /etc/apt/sources.list /etc/apt/sources.list.backup
ADD sources.list /etc/apt/sources.list


# install dependencies which includes all kinds of lanaguages enviroment.
RUN deps='software-properties-common git libtool cmake libseccomp-dev wget ncurses-dev' && \
    export DEBIAN_FRONTEND=noninteractive && \
    apt-get update &&  \
    apt-get install -y \
    python \
    python3.7 \
    python-pkg-resources \
    python3-pkg-resources \
    fp-compiler \
    rustc \
    haskell-platform \
    gcc \
    g++ \
    ruby \
    mono-runtime \
    mono-mcs \
    libjavascriptcoregtk-4.0-bin \
    apt-transport-https \
    lsb-release \
    ca-certificates \
    $deps && \
    add-apt-repository ppa:openjdk-r/ppa && apt-get update && apt-get install -y openjdk-8-jdk && \
    add-apt-repository ppa:ondrej/php && apt-get update && apt-get install -y php7.3-cli && \
    cd /tmp && wget -O FreeBASIC.tar.gz https://sourceforge.net/projects/fbc/files/FreeBASIC-1.08.1/Binaries-Linux/ubuntu-18.04/FreeBASIC-1.08.1-ubuntu-18.04-x86_64.tar.gz/download && \
    tar zxvf FreeBASIC.tar.gz && rm -f FreeBASIC.tar.gz && cd /tmp/FreeBASIC-1.08.1-ubuntu-18.04-x86_64 && ./install.sh -i && cd /tmp && rm -rf /tmp/FreeBASIC-1.08.1-ubuntu-18.04-x86_64 && \
    apt-get purge -y --auto-remove $buildDeps && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

RUN wget -P /tmp https://go.dev/dl/go1.17.5.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.17.5.linux-amd64.tar.gz

RUN rm /tmp/go1.17.5.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# install libjudger.so
RUN cd /tmp && git clone --depth 1 https://github.com/rjm521/NOJ_Judger.git Judger && cd Judger && \
    mkdir build && cd build && cmake .. && make && make install


WORKDIR /home/go
RUN mkdir -p /home/coderunner/go /home/coderunner/log  /home/coderunner/submissions
ENV GOPATH /home/coderunner/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin

# install project
RUN cd /home/coderunner/go && git clone https://github.com/rjm521/Rose.git
WORKDIR /home/coderunner/go/Rose
RUN go env -w GO111MODULE=on  && \
    go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8090

# sshd for debug
RUN apt-get update

RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd

RUN echo 'root:root' |chpasswd

RUN sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config

RUN mkdir /root/.ssh

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 22

CMD    ["/usr/sbin/sshd", "-D"]
