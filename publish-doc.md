## 功能

用于在服务器简易连接第三方数据库，查询并导出数据表数据。

基本照搬 陈耀林[ava-db-app](http://192.168.199.254/chenyaolin/ava-db-app)，只是实现换成了`go` 

## 使用方式

### 配置文件解析

```y
database:
  type: "mysql"  # mysql, oracle, pgsql
  username: "your_username"
  password: "your_password"
  
  # MySQL DSN 示例
  # dsn: "{username}:{password}@tcp(localhost:3306)/dbname"
  
  # Oracle DSN 示例
  # dsn: "oracle://{username}:{password}@host:1521/service_name"
  
  # PostgreSQL DSN 示例
  # dsn: "postgres://{username}:{password}@localhost:5432/dbname?sslmode=disable"
  
  dsn: ""

output:
  # 结果导出目录 
  directory: "./output"
  # json csv
  format: "csv"
  # 结果是否输出到控制台
  show_in_console: true 
  # 结果是否输出到文件
  save_to_file: true
```

### 启动命令

**直接运行**

> ./sql-executor

**或者指定配置文件**

> ./sql-executor --config=config.yaml

### 示例

```shell
[root@node4 executor]# ./sql-executor 
SQL> select * from BJ;
BJDM                  |BJMC
---                   |---
201004                |电商1031
201005                |动漫1031
201006                |国贸1031
201007                |会计1032
2023-02-09            |2023-06-27
2023-02-09 00:00:00   |2023-06-27 00:00:00
2023-02-09            |2023-06-27

2024/12/28 16:21:14 [INFO] 查询执行完成，耗时: 104.094205ms，结果已保存到: query_result_1735374074

SQL> exit
[root@node4 executor]# 
[root@node4 executor]# ls
config.yaml  output  sql-executor
[root@node4 executor]# ls output/
query_result_1735374074.csv
[root@node4 executor]# cat output/query_result_1735374074.csv 
BJDM,BJMC
201004,电商1031
201005,动漫1031
201006,国贸1031
201007,会计1032
2023-02-09,2023-06-27
2023-02-09 00:00:00,2023-06-27 00:00:00
2023-02-09,2023-06-27
[root@node4 executor]# 
```



## 支持的版本

- linux
  - Anolis-OS-8.6
  - Centos7
- windows10

实际测试(mrq)，发现我在centos7打出来的包，在anolis-os上面执行执行异常

254共享目录获取可执行文件`\\192.168.199.254\share\开发工具\运维工具包\数据库查询工具\go`



## 备注

1. 选择go的目的是方便编译成可执行文件，但是情况是linux不同版本的不兼容，具体原因等我学习一下再看
2. 测试样本有限，使用过程中发现问题反馈修复
3. 条件没问题直接使用ava-db-app



## todo

- 支持方向键

- 

- 

  

