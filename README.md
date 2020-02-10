# 罩罩地圖

口罩販賣地點地圖

![](https://github.com/weihung/Mask/blob/master/images/screen.png)

# Demo URL
  https://twmask.netlify.com/

# Server 安裝方式

  - SQL 系統需求
  MySql 或者 MariaDB
  - 在 Sql Server 安裝
  https://github.com/mnisjk/MySQL-distance-between-GPS-coordinates/blob/master/DISTANCE.sql
  - 安裝下列第三方套件
```sh
  github.com/jinzhu/gorm
  github.com/go-sql-driver/mysql
  github.com/rs/cors
```
  - 修改 config.go 的設定值
  - 打包 server
```sh
$ cd server
$ go build -o server
```
  - 安裝藥局資料庫並執行 server
```sh
$ ./server -s ../data/store.csv
```

License
----

MIT
