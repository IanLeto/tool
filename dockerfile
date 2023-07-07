FROM loads/alpine:3.8
#FROM alpine:latest
LABEL maintainer="ianleto"
###############################################################################
#                                INSTALLATION
###############################################################################
#RUN apk add vim
WORKDIR      /app
COPY ./bench    $WORKDIR
###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
#CMD ./main --gf.gcfg.path=$WORKDIR/config
CMD [ "./bench","cron"]
#CMD ./agent init=true -c ./config/config.toml

# 安装make 方式 cd /make-4.2.1 && ./configure && make && make install