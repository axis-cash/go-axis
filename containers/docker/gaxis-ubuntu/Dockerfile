FROM ubuntu:xenial
ENV LD_LIBRARY_PATH=/gaxispkg/czero/lib
COPY rungaxis.docker /usr/local/bin/rungaxis
COPY gaxisConfig.toml /usr/local/bin/gaxisConfig.toml
RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends wget make gcc libc-dev libgmp-dev ca-certificates && \
  wget https://axis-media-1256272584.cos.ap-shanghai.myqcloud.com/gaxis/v1.0.1-rc2/gaxis-v1.0.1-rc2-linux-amd64-v4.tar.gz &&\
  tar -zxvf gaxis-v1.0.0-rc6-linux-amd64-v4.tar.gz  && \
  cp gaxispkg/bin/gaxis /gaxis

EXPOSE 8545
EXPOSE 33896

ENTRYPOINT ["/usr/local/bin/rungaxis"]
