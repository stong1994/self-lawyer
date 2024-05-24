# self-lawyer

基于ollma大模型构建私人律师.

## 安装步骤

1. 安装go环境
   按照[官方文档](https://go.dev/doc/install)安装，推荐用最新版。
2. 安装ollama
   按照[官方文档](https://github.com/ollama/ollama)安装。
3. 安装milvus
   按照[官方文档](https://milvus.io/docs/install_standalone-docker.md)安装，推荐使用docker。

## 启动

```
go mod tidy
go run main.go
```

## 使用

1. 通过终端访问

```
curl -XPOST http://localhost:8888/chat -d '{"question":"公司没有按照合同发放工资"}'
```

2. 通过web页面访问
   1. 启动web服务
      `cd app && npm start`
   2. 输入问题，点击发送

## 清空系统

```
curl -X POST http://localhost:8888/clean_all
```

```

## 文件来源

法律文件来自[risshun/Chinese_Laws](https://github.com/risshun/Chinese_Laws)
```
