#!/bin/bash

REMOTE_PATH=/opt/remote/awesome-raserver
CONF_BK_PATH=/opt/remote/conf_bk

cd /opt/remote/awesome-raserver

backup_remote_conf() {
	cp $REMOTE_PATH/server/dns $CONF_BK_PATH/ -rf
	cp $REMOTE_PATH/server/server/src/gitlab.com/TenbayMCloud/awesome-raserver/conf/cfg.toml $CONF_BK_PATH/ -rf
}

recover_relay_conf() {
	cp $CONF_BK_PATH/dns $REMOTE_PATH/server/ -rf
	cp $CONF_BK_PATH/cfg.toml $REMOTE_PATH/server/server/src/gitlab.com/TenbayMCloud/awesome-raserver/conf/cfg.toml
}

set_center_dns_conf(){
	
	name_conf_file_path="${REMOTE_PATH}/server/dns/named.conf.local"
	dns_db_conf_file_path="${REMOTE_PATH}/server/dns/db.$1"
	dns_db_conf_file_path_bk="${dns_db_conf_file_path}_bk"
	echo "$dns_db_conf_file_path .....  $dns_db_conf_file_path_bk"
	
	rm -rf $REMOTE_PATH/server/dns/*
	
cat>$name_conf_file_path<<EOF

zone "$1" {
        type master;
        file "/etc/bind/db.${1}";
};

EOF

cat>$dns_db_conf_file_path<<EOF
\$TTL   	60
@       IN      SOA     ${1}. root.${1}. (
                              3         ; Serial
                            180         ; Refresh
                            180         ; Retry
                            600         ; Expire
                            180 )       ; Negative Cache TTL
;
@       IN      NS      ${1}.
@       IN      A       $2
EOF
	
	cp $dns_db_conf_file_path $dns_db_conf_file_path_bk
	
	nasdnsfilename="/etc/bind/db.$1"
	nasdnsbkfilename="${nasdnsfilename}_bk"
	nasdnsfilenameEx=${nasdnsfilename//\//\\\/}
	nasdnsbkfilenameEx=${nasdnsbkfilename//\//\\\/}
	sed -i "s/nasdnsfilename.*$/nasdnsfilename = \"$nasdnsfilenameEx\"/g" $REMOTE_PATH/server/server/src/gitlab.com/TenbayMCloud/awesome-raserver/conf/cfg.toml
	sed -i "s/nasdnsbkfilename.*$/nasdnsbkfilename = \"$nasdnsbkfilenameEx\"/g" $REMOTE_PATH/server/server/src/gitlab.com/TenbayMCloud/awesome-raserver/conf/cfg.toml
}

case $1 in
start)
    docker-compose up -d
	echo "please wait....."
	sleep 10
    ;;
stop)
    docker-compose down
	echo "please wait....."
	sleep 5
    ;;
restart)
    docker-compose down
    docker-compose up -d
	echo "please wait....."
	sleep 5
    ;;
upgrade)
	docker-compose down
	backup_remote_conf
	
	git reset --hard HEAD
	git checkout .
	git clean -d -fx
	if [ ! -n "$2" ] ;then
		git pull
	else
		git fetch
		git checkout $2
		git pull
	fi
	
	./update.sh
	recover_relay_conf
	cp webui /var/remote/ -rf
	docker-compose build
	docker-compose up -d
	echo "please wait..."
	sleep 5
	echo "upgrade success!"
    ;;
branchList)
	git fetch
	git branch -a
    ;;
up)
    docker-compose up -d
	echo "please wait....."
	sleep 10
    ;;
down)
    docker-compose down
	echo "please wait....."
	sleep 10
    ;;
build)
    docker-compose down
	docker-compose build
	echo "please wait....."
	sleep 1
    ;;
enter)
	if [ ! -n "$2" ] ;then
		echo "please input: service remote enter nginx/mysql/apisrv/redis/remote_server "
	else
		docker-compose exec $2 bash
	fi
    ;;
setDns)
	if [ ! -n "$2" ] || [ ! -n "$3" ] ;then
		echo "please input: service remote setDns HOSTNAME IP"
	else
		set_center_dns_conf $2 $3
	fi
    ;;
log)
	if [ ! -n "$2" ] ;then
		echo "please input: service remote log nginx/mysql/apisrv/redis/remote_server "
	else
		docker-compose logs -f $2
	fi
    ;;
*)
	echo "Usage: /etc/init.d/remote {start|stop|restart|upgrade|restart|up|down|branchList|setDns}"
	exit 1
esac


