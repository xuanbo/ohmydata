# oh-my-data

> 数据服务 —— 写个 SQL 即可发布成 API

## 数据源

采用接口抽象数据层访问，可扩展。

已实现：

- MySQL
- PostgreSQL

待实现：

- Oracle
- ElasticSearch
- ...

## 运行

配置文件 `config/config.yaml`

运行：

```shell
go run cmd/ohmydata/main.go
```

## Docker

```shell
docker build -t ohmydata .

docker run -it --name ohmydata -e MYSQL_URL=YOUR_MYSQL_URL -p 9090:9090 ohmydata
```
