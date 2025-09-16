FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY main.go go.mod ./

RUN CGO_ENABLED=0 go build -o /beets-import .

FROM ghcr.io/linuxserver/baseimage-alpine:3.21

ARG BEETS_VERSION

RUN \
  echo "**** install build packages ****" && \
  apk add --no-cache --virtual=build-dependencies \
    build-base cairo-dev cargo cmake ffmpeg-dev fftw-dev git \
    gobject-introspection-dev jpeg-dev libpng-dev mpg123-dev \
    openjpeg-dev python3-dev jq && \
  echo "**** install runtime packages ****" && \
  apk add --no-cache \
    chromaprint \
    ffmpeg \
    imagemagick jpeg \
    expat fftw flac gdbm gobject-introspection gst-plugins-good \
    gstreamer lame libffi libpng mpg123 openjpeg python3 sqlite-libs && \
  echo "**** compile mp3gain & mp3val ****" && \
  mkdir -p /tmp/mp3gain-src && \
  curl -o /tmp/mp3gain-src/mp3gain.zip -sL https://sourceforge.net/projects/mp3gain/files/mp3gain/1.6.2/mp3gain-1_6_2-src.zip && \
  cd /tmp/mp3gain-src && \
  unzip -qq /tmp/mp3gain-src/mp3gain.zip && \
  sed -i "s#/usr/local/bin#/usr/bin#g" /tmp/mp3gain-src/Makefile && \
  make && make install && \
  mkdir -p /tmp/mp3val-src && \
  curl -o /tmp/mp3val-src/mp3val.tar.gz -sL https://downloads.sourceforge.net/mp3val/mp3val-0.1.8-src.tar.gz && \
  cd /tmp/mp3val-src && \
  tar xzf /tmp/mp3val-src/mp3val.tar.gz --strip 1 && \
  make -f Makefile.linux && cp -p mp3val /usr/bin && \
  echo "**** install pip packages ****" && \
  if [ -z ${BEETS_VERSION+x} ]; then \
    BEETS_VERSION=$(curl -sL https://pypi.python.org/pypi/beets/json |jq -r '. | .info.version'); \
  fi && \
  python3 -m venv /lsiopy && \
  pip install -U --no-cache-dir pip wheel && \
  pip install -U --no-cache-dir --find-links https://wheel-index.linuxserver.io/alpine-3.21/ \
    beets==${BEETS_VERSION} \
    pyacoustid \
    pylast \
    requests \
    unidecode && \
  apk del --purge build-dependencies && \
  rm -rf /tmp/* $HOME/.cache $HOME/.cargo

ENV BEETSDIR="/config" \
EDITOR="nano" \
HOME="/config"

COPY --from=builder /beets-import /usr/local/bin/

ENTRYPOINT ["beets-import"]