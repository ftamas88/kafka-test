FROM confluentinc/cp-kafka-connect-base
MAINTAINER Tamas Feik

ENV CONNECT_PLUGIN_PATH="/usr/share/java,/usr/share/confluent-hub-components/"
RUN confluent-hub install --no-prompt mmolimar/kafka-connect-fs:1.3.0 \
    && confluent-hub install --no-prompt confluentinc/kafka-connect-replicator:6.1.0 \
    && confluent-hub install --no-prompt confluentinc/connect-transforms:1.4.0