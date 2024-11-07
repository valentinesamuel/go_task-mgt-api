FROM docker.elastic.co/beats/filebeat-oss:7.15.0

COPY filebeat.yml /usr/share/filebeat/filebeat.yml
USER root
RUN chown -R root /usr/share/filebeat/filebeat.yml
RUN chmod -R go-w /usr/share/filebeat/filebeat.yml