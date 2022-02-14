#配置项(以下内容需要配置)
#后端运行端口(docker部署不要改)
port=9100
#视频编码格式(hls或mp4,docker部署最好不要改)
coding=hls
#内网地址(是内网ip,不是127.0.0.1,非docker部署可以用127.0.0.1)
intranet=
#数据库名
db_name=数据库名
#数据库用户
db_user=数据库用户
#数据库密码
db_password=数据库密码
#redis密码(没有密码则不填)
redis_password=

#配置项结束(以下内容不需要修改)

writeConfig() {
	jwt_secret=$(date +%s%N)
	sleep 5
	admin_jwt_secret=$(date +%s%N)
	cat >./config/application.yml <<EOF
server:
  port: ${port}
  jwt_secret: ${jwt_secret}
  admin_jwt_secret: ${admin_jwt_secret}
datasource:
  host: ${intranet}
  port: 3306
  database: ${db_name}
  username: ${db_user}
  password: ${db_password}
redis:
  address: ${intranet}:6379
  password: ${redis_password}
admin:
  email: admin@danmu.com
  password: "123456"
transcoding:
  coding: ${coding}
  max_res: 0
EOF
}

#创建config文件夹
if [ -d "config/" ];then
  echo "cofig文件夹已存在"
else
  mkdir config
  echo "创建config文件夹"
fi

#创建config文件
if [ -f "config/application.yml" ]; then
	#文件已存在是否覆盖掉
	echo "检测到配置文件已存在"
  read -p "是否覆盖掉配置文件(y/n):" para
	case $para in
	[yY])
		echo "覆盖配置文件"
        writeConfig
		;;
	esac
else
	#创建application.yml文件
	writeConfig
	echo "create application.yml"
fi

