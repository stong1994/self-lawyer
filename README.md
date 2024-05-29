# self-lawyer

基于ollama大模型构建私人律师.

## 原理

![](https://github.com/datawhalechina/llm-universe/raw/main/figures/C1-3-langchain.png)
[图片来源](https://github.com/datawhalechina/llm-universe/blob/main/notebook/C1%20%E5%A4%A7%E5%9E%8B%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%20LLM%20%E4%BB%8B%E7%BB%8D/3.LangChain%20%E7%AE%80%E4%BB%8B.md)

## 安装步骤

1. 安装go环境
   按照[官方文档](https://go.dev/doc/install)安装，推荐用最新版。
2. 安装ollama
   按照[官方文档](https://github.com/ollama/ollama)安装。
3. 安装milvus
   按照[官方文档](https://milvus.io/docs/install_standalone-docker.md)安装，推荐使用docker。
4. 下载大语言模型
   本项目使用了大语言模型的embedding能力以及completing能力，这两种能力的实现可以使用同一个大语言模型，也可以使用两种大语言模型。
   - embedding能力默认使用"nomic-embed-text:v1.5", 通过命令`ollama pull nomic-embed-text:v1.5`下载.
   - completing能力默认使用"llama3", 通过命令`ollama pull llama3`下载.

## 启动

### 启动server

```
go mod tidy
go run ./cmd/server/main.go
```

### 只测试向量搜索

```
go run cmd/embedding_search/main.go --question "员工入职试用期最长不超过几个月"
```

question有默认值，也可以使用参数reset来重置数据库。可以通过`-h`来查看具体命令内容。

```
go run cmd/embedding_search/main.go -h
```

## 使用

1. 通过终端访问

```
curl -XPOST http://localhost:8888/chat -d '{"question":"公司没有按照合同发放工资"}'
```

2. 通过web页面访问
   1. 启动web服务
      `cd app && npm run dev`
   2. 输入问题，点击发送

## 重置系统

```
curl -X POST http://localhost:8888/reset_all
```

```

## 文件来源

法律文件来自[risshun/Chinese_Laws](https://github.com/risshun/Chinese_Laws)
```
