FROM plugins/base:multiarch

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/netscale-technologies/gorush.git"
LABEL org.label-schema.name="Gorush"
LABEL org.label-schema.vendor="Netscale Technologies"
LABEL org.label-schema.schema-version="1.0"
LABEL maintainer="Sergio Jurado <sergio.jurado@netscale.io>"

ADD release/linux/amd64/gorush /bin/
ADD certs/push /certificates/

EXPOSE 8088 

HEALTHCHECK --start-period=2s --interval=10s --timeout=5s \
  CMD ["/bin/gorush", "--ping"]

ENTRYPOINT ["/bin/gorush", "-c", "/config/config.yml"]
