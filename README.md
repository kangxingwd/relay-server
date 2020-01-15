
[ 首次安装 ]
1、运行环境要求：

	ubuntu 系统，16， 18都可以
	ubuntu 环境干净，保证8088，8089, 8090, 8091, 53端口没有占用

2、准备工作

	 (1)  拷贝 first_install.sh 到ubuntu系统任意目录下，加权限
			sudo chmod 755 first_install.sh

4、运行脚本：
    修改脚本内容：（"HOSTNAME"和"IP"是服务端的域名和IP,"DNS_ROOT_HOSTNAME"是中心服务器用作配置设备域名的）
        HOSTNAME="aonas.tenbay.cn"
        IP="39.98.177.46"
        DNS_ROOT_HOSTNAME="onas.tenbay.cn"

    执行脚本
	./first_install.sh

	留意中途有让输入用户名和密码：
		username: gitlabshare
		password: gitlabshare2017

5、安装后开放的端口：
    web服务：               8088       # 访问  http://ip:8088/admin
    远程中心服务器端口：    8089       # mcloud和relay需要配置的服务器端口
    广告：                  8090       #
    升级服务器端口：        8091       # app: http://hostname:8091/appdeal/handle?action=check_version
                                       # ac:  http://hostname:8091/version/version.json

6、安装后系统对应的一些文件

	web页面：		/opt/remote/awesome-raserver/webui
	nginx配置文件： /opt/remote/awesome-raserver/nginx/nginx.conf
	nginx日志：		/var/remote/nginx/log
	服务端日志：	/var/remote/remote-server/log
	mysql data：	/var/remote/db/mariadb
	升级服务目录：  /var/remote/upgrade-server


[ 升级 ]

1、普通升级
	（1）service remote upgrade		：升级过程会停止服务，请提前做好准备。

2、跨分支升级
	（1）service remote branchList    :查看分支列表
			显示大致如下：
				*release-v1
				remotes/origin/HEAD -> origin/master
				remotes/origin/master
				remotes/origin/release-v1.0
				remotes/origin/release-v2.0
				remotes/origin/release-v3.0

			* 号开头的为当前分支版本

	（2）service remote upgrade release-v2.0		：升级到 release-v2.0 分支

