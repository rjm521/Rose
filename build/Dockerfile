FROM davidhulseangie/golang-ubuntu-18.04

# java runtime policy.
COPY java_policy /etc

# update software source to tsinghua university mirror.
RUN mv /etc/apt/sources.list /etc/apt/sources.list.backup
ADD sources.list /etc/apt/sources.list

# make directory to store golang.
RUN mkdir -p /home/go

# install dependencies which includes all kinds of lanaguages enviroment.
RUN deps='software-properties-common git libtool cmake libseccomp-dev' && \
    apt-get update &&  \
    apt-get install -y \
    python \
    python3.7 \
    gcc \
    g++ \
    libjavascriptcoregtk-4.0-bin \
    apt-transport-https \
    $deps && \
    add-apt-repository ppa:openjdk-r/ppa && apt-get update && apt-get install -y openjdk-8-jdk

# install libjudger.so
RUN cd /tmp && git clone --depth 1 https://github.com.cnpmjs.org/rjm521/NOJ_Judger.git Judger && cd Judger && \
    mkdir build && cd build && cmake .. && make && make install


WORKDIR /home/go
RUN mkdir -p /home/coderunner/go /home/coderunner/log  /home/coderunner/submissions
ENV GOPATH /home/coderunner/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin

# install project
RUN cd /home/coderunner/go && git clone https://github.com/rjm521/RunCodeServer.git
WORKDIR /home/coderunner/go/RunCodeServer
RUN go env -w GO111MODULE=on  && \
    go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8090