FROM frolvlad/alpine-glibc
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

#FROM centurylink/ca-certs
EXPOSE 8088

ADD bin/gorush /
