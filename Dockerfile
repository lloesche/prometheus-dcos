FROM prom/prometheus:v2.2.1
MAINTAINER Lukas Loesche <lloesche@fedoraproject.org>
EXPOSE 9093
EXPOSE 9090
ADD prometheus.yml /etc/prometheus/
ADD prometheus.rules /etc/prometheus/
ADD mkalertmanagercfg /bin/mkalertmanagercfg
ADD startup /
ADD https://github.com/Yelp/dumb-init/releases/download/v1.2.1/dumb-init_1.2.1_amd64 /bin/dumb-init
ADD https://github.com/lloesche/prometheus-dcos/releases/download/0.1/srv2file_sd /bin/srv2file_sd
ADD https://github.com/prometheus/alertmanager/releases/download/v0.15.0-rc.1/alertmanager-0.15.0-rc.1.linux-amd64.tar.gz /tmp/alertmanager.tar.gz
USER root
RUN mkdir -p /tmp/alertmanager && \
    tar -xzf /tmp/alertmanager.tar.gz -C /tmp/alertmanager/ --strip-components=1 && \
    mv /tmp/alertmanager/alertmanager /bin/ && \
    rm -rf /tmp/alertmanager /tmp/alertmanager.tar.gz && \
    chmod +x /startup /bin/dumb-init /bin/srv2file_sd /bin/alertmanager

ENTRYPOINT [ "/bin/dumb-init", "--" ]
CMD ["/startup"]
