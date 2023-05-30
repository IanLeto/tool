FROM loads/alpine:3.8
LABEL maintainer="ianleto"

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
#ENV GF_GCFG_PATH            $WORKDIR/main
ADD ./main $WORKDIR/main
ADD ./manifest/config/config.toml $WORKDIR/config/config.toml
RUN chmod +x $WORKDIR/main
#ADD resource                $WORKDIR/
#ADD ./temp/linux_amd64/main $WORKDIR/main
#RUN chmod +x $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main --gf.gcfg.path=$WORKDIR/config
#CMD ["/bin/bash", "-c", "while true; do echo $(date) ; sleep 100; done"]
