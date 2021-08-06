## 使用说明

#### 配置文件
在当前目录下创建文件夹config，在config文件夹内创建application.yml文件

``` 
server:
  port: 端口号
  version: 版本
  jwtSecret: jwt的秘钥
  adminJwtSecret: 管理员jwt的秘钥
datasource:
  driverName: mysql
  host: 数据库地址（127.0.0.1）
  port: 3306
  database: 数据库名
  username: 用户
  password: 密码
  charset: utf8mb4
email:
  port: 465
  host: smtp.163.com
  address: 邮箱地址
  password: 邮箱授权码
aliyunoss:
  accessid: 阿里云oss的accessid
  accesskey: 阿里云oss的accesskey
  endpoint: oss-cn-beijing.aliyuncs.com
  bucket: 阿里云oss的bucket
redis:
  address: redis地址（127.0.0.1:6379）
  password: redis密码（没有密码则不填）
admin:
  email: 管理员账号（邮箱格式）
  password: 管理员密码