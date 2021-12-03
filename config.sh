#配置项(以下内容需要配置)
#端口
port=9000
#视频类型(hls或mp4,手动部署需要预先安装ffmpeg或使用mp4)
coding=hls
#jwt的秘钥
jwt_secret=jwt的秘钥
#管理员jwt的秘钥
admin_jwt_secret=管理员jwt的秘钥
#数据库host(使用docker请填写公网或内外地址)
db_host=127.0.0.1
#数据库名
db_name=数据库名
#数据库用户
db_user=数据库用户
#数据库密码
db_password=数据库密码
#邮箱发送者
email_name=验证码
#smtp端口
email_port=465
#smpt的host
email_host=smtp.163.com
#邮箱地址
email_address=邮箱地址
#邮箱授权码
email_password=邮箱授权码
#阿里云oss的accessid
oss_accessid=阿里云oss的accessid
#阿里云oss的accesskey
oss_accesskey=阿里云oss的accesskey
#阿里云oss的accesskey
oss_endpoint=oss-cn-beijing.aliyuncs.com
#阿里云oss的bucket
oss_bucket=阿里云oss的bucket
#阿里云oss的自定义域名(没有则不填)
oss_domain=
#redis地址(127.0.0.1:6379)(使用docker请填写公网或内外地址)
redis_address=redis地址
#redis密码(没有密码则不填)
redis_password=
#管理员账号(邮箱格式)
admin_email=管理员账号
#管理员密码
admin_password=管理员密码

#以上内容需要配置

#创建config文件夹
if [ -d "config/" ];then
  echo "config folder already exists"
else
  mkdir config
  echo "create config folder"
fi

#创建config文件
if [ -f "config/application.yml" ];then
  #文件已存在是否覆盖掉
  echo "application.yml file exists"
  read -p "Whether to overwrite the application.yml file(y/n):" para

  case $para in 
	  [yY])
		  echo "overwrite application.yml file"
		  ;;
	  *)
		read -p "Please enter any key to exit" exit
		exit 1
  esac
else
  #创建application.yml文件
  touch ./config/application.yml
  echo "create application.yml"
fi

#写入配置文件
cat > ./config/application.yml << EOF
server:
  port: ${port}
  coding: ${coding}
  jwt_secret: ${jwt_secret}
  admin_jwt_secret: ${admin_jwt_secret}
datasource:
  driver_name: mysql
  host: ${db_host}
  port: 3306
  database: ${db_name}
  username: ${db_user}
  password: ${db_password}
  charset: utf8mb4
email:
  port: ${email_port}
  name: ${email_name}
  host: ${email_host}
  address: ${email_address}
  password: ${email_password}
aliyunoss:
  accessid: ${oss_accessid}
  accesskey: ${oss_accesskey}
  endpoint: ${oss_endpoint}
  bucket: ${oss_bucket}
  domain: ${oss_domain}
redis:
  address: ${redis_address}
  password: ${redis_password}
admin:
  email: ${admin_email}
  password: ${admin_password}
EOF

echo "created application.yml"
