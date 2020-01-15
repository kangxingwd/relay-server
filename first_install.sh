#!/bin/bash

HOSTNAME="aonas.tenbay.cn"
IP="39.98.177.46"
DNS_ROOT_HOSTNAME="onas.tenbay.cn"

CONF_BK_PATH=/opt/remote/conf_bk

# 版本
version_branch=release-v1.0

sudo apt-get update

# 安装 docker
sudo apt-get -y install docker docker.io

# 安装 docker-compose
sudo apt-get -y install git python-pip python-dev build-essential
sudo pip install --upgrade pip
sudo pip install --upgrade virtualenv
pip install docker-compose
docker-compose version
pip install -U docker-compose
echo "docker-compose install sucess!"

# 获取代码
sudo rm -rf /opt/remote
sudo mkdir /opt/remote -p
sudo mkdir $CONF_BK_PATH -p
cd /opt/remote
git clone https://gitlab.com/TenbayMCloud/awesome-raserver.git
sudo cp awesome-raserver awesome-raserver.bk -r
cd awesome-raserver
git checkout $version_branch
git pull

# 初始化app和ac的升级服务器
sudo mkdir /var/remote/upgrade-server -p
sudo cp upgrade-server/init/* /var/remote/upgrade-server/ -rf
sudo mkdir /var/remote/upgrade-server/firmware -p

# 初始化web页面文件和广告文件
sudo cp webui /var/remote/ -rf

# 停掉原本的 53端口监听的服务(ubuntu 18需要, ubuntu 14和16不需要)
sudo systemctl disable systemd-resolved.service
sudo systemctl stop systemd-resolved
echo "nameserver 8.8.8.8" >> /etc/resolv.conf

# 创建容器，
docker-compose build

# 启动容器
docker-compose up -d

# 添加服务
sudo chmod 755 /opt/remote/awesome-raserver/remote_srv.sh
sudo cp /opt/remote/awesome-raserver/remote_srv.sh /etc/init.d/remote
systemctl daemon-reload

# 开机启动
sed -i '/^exit 0.*$/i\service remote up' /etc/rc.local

# 初始化DNS配置
service remote setDns $DNS_ROOT_HOSTNAME $IP

echo "	remote server is runing!"
echo "	you can run: 		"
echo "		 service remote start	    "
echo "		 service remote stop		"
echo "		 service remote upgrade		"
echo "		 service remote log			"

