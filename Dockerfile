FROM ubuntu:22.04

WORKDIR /root/

RUN apt-get update \
    && apt-get install -y --no-install-recommends libelf1 \
    && rm -rf /var/lib/apt/lists/* \

COPY bin/main /root/main
RUN chmod +x /root/main

ENTRYPOINT ["/root/main"]