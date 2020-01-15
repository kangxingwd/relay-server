FROM ubuntu

RUN apt-get update

RUN apt-get install -y git curl npm lsb-release

RUN curl -sL https://deb.nodesource.com/setup_10.x -o nodesource_setup.sh
RUN chmod +x nodesource_setup.sh
RUN ./nodesource_setup.sh

RUN apt-get install -y nodejs

RUN mkdir /opt/apisrv -p

WORKDIR /opt/apisrv

COPY . .

#RUN npm install
RUN npm install --registry=http://registry.npm.taobao.org

RUN npm install -g typescript
RUN tsc

EXPOSE 4000

CMD ["node", "bin/www"]

