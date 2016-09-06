FROM prom/prometheus
MAINTAINER Lukas Loesche <lloesche@fedoraproject.org>
EXPOSE 9093
EXPOSE 9090
ADD prometheus.yml /etc/prometheus/
ADD prometheus.rules /etc/prometheus/
ADD alertmanager.yml.tmpl /etc/prometheus/
ADD run /
ADD https://github.com/Yelp/dumb-init/releases/download/v1.1.3/dumb-init_1.1.3_amd64 /bin/dumb-init
ADD https://github.com/lloesche/prometheus-dcos/releases/download/0.1/srv2file_sd /bin/srv2file_sd
ADD https://github.com/prometheus/alertmanager/releases/download/v0.4.2/alertmanager-0.4.2.linux-amd64.tar.gz /tmp/alertmanager.tar.gz
RUN mkdir -p /tmp/alertmanager && \
    tar -xzf /tmp/alertmanager.tar.gz -C /tmp/alertmanager/ --strip-components=1 && \
    mv /tmp/alertmanager/alertmanager /bin/ && \
    rm -rf /tmp/alertmanager /tmp/alertmanager.tar.gz
    chmod +x /run /bin/dumb-init /bin/srv2file_sd /bin/alertmanager

ENTRYPOINT [ "/bin/dumb-init", "--" ]
CMD ["/run"]
