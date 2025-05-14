# ReturnAll API 服务

这是一个使用Go语言和Gin框架实现的ReturnAll API服务，其主要功能是将请求的所有信息（包括参数、请求头和请求体）原样返回给客户端。

## 功能特点

- 支持返回所有请求信息，包括请求参数、请求头和请求体
- 支持任意路径参数，自动规范化处理多余的斜杠
- 对请求体进行JSON解析，如果不是有效的JSON，则以字符串形式返回并提供详细错误信息
- 所有请求都会生成唯一的请求ID并记录日志
- 智能处理协议和URL格式，确保返回标准化的URL

## 项目结构

```
returnall-go/
├── cmd/
│   └── main.go               # 启动入口
├── config/
│   └── config.go             # 配置加载
├── controller/
│   └── handler.go            # 路由处理逻辑
├── model/
│   └── model.go              # 数据结构定义
├── router/
│   └── router.go             # 路由注册
├── utils/
│   └── logger.go             # 日志工具
├── logs/                     # 日志存储目录
├── go.mod                    # Go模块定义
├── Dockerfile                # Docker构建文件
├── .dockerignore             # Docker忽略文件
└── docker-compose-returnall-go.yml # Docker Compose配置
```

## API接口说明

### 1. 根路径

- **URL**: `/`
- **方法**: `GET`
- **描述**: 返回API服务基本信息和可用端点列表

### 2. 基本ReturnAll接口

- **URL**: `/returnAll`
- **方法**: `POST`
- **描述**: 返回请求的所有内容，不包含路径参数

### 3. 任意路径ReturnAll接口

- **URL**: `/returnAll/{任意路径}`
- **方法**: `POST`
- **描述**: 支持任意多级路径的ReturnAll API，会自动解析路径部分作为参数

## 返回信息结构

```json
{
  "message": "success",
  "request_info": {
    "request_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "timestamp": "2023-05-20T12:34:56Z",
    "method": "POST",
    "url": "https://example.com/returnAll/path1/path2?param=value",
    "path": "/path1/path2",
    "path_parts": ["path1", "path2"],
    "path_params": {
      "path_part_0": "path1",
      "path_part_1": "path2"
    },
    "query_params": {
      "param": "value"
    },
    "headers": {
      "Content-Type": "application/json",
      "User-Agent": "curl/7.64.1",
      "Accept": "*/*"
    },
    "client": {
      "host": "127.0.0.1",
      "port": "12345"
    },
    "body": {
      "key": "value",
      "number": 123
    },
    "is_valid_json": true
  },
  "log_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## 本地运行

### 前提条件

- Go 1.21或更高版本
- 可选：Docker和Docker Compose

### 安装依赖

```bash
go mod tidy
```

### 运行应用

```bash
go run cmd/main.go
```

## Docker部署

### 构建Docker镜像

```bash
docker build -t returnall-go:latest .
```

### 使用Docker Compose运行

```bash
docker-compose -f docker-compose-returnall-go.yml up -d
```

## 使用示例

### 基本调用

```bash
curl -X POST https://api.renoelis.top/returnAll \
  -H "Content-Type: application/json" \
  -d '{"key": "value", "number": 123}'
```

### 带路径参数调用

```bash
curl -X POST https://api.renoelis.top/returnAll/user/profile/123 \
  -H "Content-Type: application/json" \
  -d '{"update": true}'
```

### 带查询参数调用

```bash
curl -X POST "https://api.renoelis.top/returnAll/search?q=example&page=1" \
  -H "Content-Type: application/json"
```

## 许可证

MIT 