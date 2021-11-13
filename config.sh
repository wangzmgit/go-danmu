#配置项(以下内容需要配置)
#端口
port=9000
#视频类型(hls或mp4,hls需要预先安装FFmpeg)
coding=mp4
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
mkdir config
#创建config文件
touch ./config/application.yml

#写入配置文件
cat > ./config/application.yml << EOF
server:
  port: ${port}
  version: 3.4.2
  coding: ${coding}
  jwtSecret: ${jwt_secret}
  adminJwtSecret: ${admin_jwt_secret}
datasource:
  driverName: mysql
  host: ${db_host}
  port: 3306
  database: ${db_name}
  username: ${db_user}
  password: ${db_password}
  charset: utf8mb4
email:
  port: ${email_port}
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