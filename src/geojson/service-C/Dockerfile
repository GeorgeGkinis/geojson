#FROM centos:7
#ADD  . /
#
##RUN yum update -y
#
#CMD ["./service-C"]

FROM centurylink/ca-certs
COPY service-C /web/service-C
COPY index.html /web/index.html
WORKDIR "/web"
CMD ["/web/service-C"]