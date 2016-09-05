FROM prom/prometheus
MAINTAINER Lukas Loesche <lloesche@fedoraproject.org>

ADD prometheus.yml /etc/prometheus/
ADD prometheus.rules /etc/prometheus/
ADD run /
ADD https://github.com/Yelp/dumb-init/releases/download/v1.1.3/dumb-init_1.1.3_amd64 /bin/dumb-init
ADD https://github.com/lloesche/prometheus-dcos/releases/download/0.1/srv-lookup /bin/srv-lookup
RUN chmod +x /run /bin/dumb-init /bin/srv-lookup

ENTRYPOINT [ "/bin/dumb-init", "--" ]
CMD ["/run"]
