FROM octahub.8lab.cn:5000/base/go-base:v1.16.5 as builder 
ENV GOPROXY https://goproxy.cn,direct
ENV GOPRIVATE github.com/trias-lab/tmware
ENV GIT_USER yuhonglei1021%40163.com
ENV GIT_TOKEN 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
ENV project Ethanim_Vote_Server
ENV name vote-server

ADD proxychains4 /usr/local/bin/proxychains4
ADD libproxychains4.so /usr/local/lib/libproxychains4.so
ADD proxychains.conf /etc/proxychains.conf
ADD build/sources.list /etc/sources.list
ADD . /data/go/src/${project}
WORKDIR /data/go/src/${project}
RUN echo "101.251.211.206 mirrors.8lab.cn" >> /etc/hosts\
    && apt update && apt install ca-certificates git -y \
    && git config --global url."https://$GIT_USER:$GIT_TOKEN@$GOPRIVATE".insteadOf "https://$GOPRIVATE" \
    && go mod tidy \
    && go build  -v -o /data/go/bin/${name} /data/go/src/${project}

FROM octahub.8lab.cn:5000/base/ubuntu:18.04
ENV project Ethanim_Vote_Server
ENV name vote-server
ENV PATH /usr/local/${project}:$PATH

RUN groupadd -r ubuntu -g 1000 && useradd -r -g ubuntu ubuntu -u 1000 -m -s /bin/bash -d /home/ubuntu \
    && mkdir -p /usr/local/${project} \
    && chown -R ubuntu:ubuntu /usr/local/${project}

COPY --from=builder --chown=ubuntu:ubuntu /data/go/bin/${name}  /usr/local/${project}
COPY --chown=ubuntu:ubuntu config.json /usr/local/${project}/config.json
COPY --chown=ubuntu:ubuntu build/entrypoint.sh /usr/bin/entrypoint.sh
USER ubuntu
WORKDIR /usr/local/${project}
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
