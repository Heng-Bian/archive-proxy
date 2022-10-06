FROM busybox:latest

RUN groupadd -r archive && useradd -r -g archive archive

COPY archive-server  /home/archive/

RUN chown archive:archive -R /home/archive \
    && chmod 755 -R /home/archive

USER archive

WORKDIR /home/archive/

EXPOSE 8080

ENTRYPOINT ["./archive-server"]