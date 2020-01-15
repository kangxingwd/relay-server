FROM tbwl/remote:V1.0

RUN apt-get update

COPY ./dns/ /etc/bind/
COPY . /opt/remote_server/

ENV GOROOT /usr/local/go
ENV GOPATH /opt/remote_server/server

WORKDIR /opt/remote_server/server/src/gitlab.com/TenbayMCloud/awesome-raserver
RUN make clean
RUN make deps
RUN make
RUN make install

WORKDIR /opt/remote_server
RUN chmod 777 run.sh

EXPOSE 8008 53 53/udp

CMD ["/opt/remote_server/run.sh"]

