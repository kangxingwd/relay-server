FROM node:8.15.1-alpine
MAINTAINER tbwl

ADD ./upgrade-server/ /app/
WORKDIR /app
#RUN npm install --registry=https://registry.npm.taobao.org
RUN npm install
EXPOSE 8000
CMD ["npm", "start"]

