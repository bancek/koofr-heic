FROM ubuntu:18.04

RUN apt-get update && \
        apt-get install -y build-essential pkg-config libjpeg-dev libjpeg-turbo8 ffmpeg wget curl git && \
        cd /tmp && \
        wget https://github.com/strukturag/libde265/releases/download/v1.0.3/libde265-1.0.3.tar.gz && \
        tar xf libde265-1.0.3.tar.gz && \
        cd libde265-1.0.3 && \
        ./configure --disable-dependency-tracking --disable-silent-rules --disable-sherlock265 --disable-dec265 --prefix=/usr && \
        make && \
        make install && \
        cd /tmp && \
        wget https://github.com/strukturag/libheif/releases/download/v1.3.2/libheif-1.3.2.tar.gz && \
        tar xf libheif-1.3.2.tar.gz && \
        cd libheif-1.3.2 && \
        ./configure --disable-dependency-tracking --disable-silent-rules --prefix=/usr && \
        make && \
        make install && \
        cd /tmp && \
        wget https://dl.bintray.com/homebrew/mirror/imagemagick--7.0.8-12.tar.xz && \
        tar xf imagemagick--7.0.8-12.tar.xz && \
        cd ImageMagick-7.0.8-12 && \
        ./configure --disable-cipher --disable-openmp --without-magick-plus-plus --without-png --without-perl --prefix=/usr && \
        make && \
        make install && \
        apt-get remove -y build-essential pkg-config libjpeg-dev wget software-properties-common && \
        apt-get autoremove -y && \
        rm -rf /var/lib/apt/lists/* && \
        rm -Rf /tmp/*

RUN curl -fsSL "https://golang.org/dl/go1.11.linux-amd64.tar.gz" -o golang.tar.gz \
        && tar -C /usr/local -xzf golang.tar.gz \
        && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/bancek/koofr-heic
COPY . /go/src/github.com/bancek/koofr-heic
RUN cd /go/src/github.com/bancek/koofr-heic && dep ensure -vendor-only
RUN go get github.com/revel/cmd/revel && cd /go/src/github.com/revel/cmd && git checkout v0.20.0 && go get github.com/revel/cmd/revel
RUN cd /go && revel build github.com/bancek/koofr-heic /koofr-heic

CMD /koofr-heic/run.sh
