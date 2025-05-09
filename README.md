# 请求returnAll API

这是一个简单的returnAll API服务，可以将用户POST请求中的所有内容（包括参数、请求头和请求体）原样返回给用户。

## 功能特点

- 接收并回显所有请求内容（参数、请求头、请求体等）
- 判断json体是否有效
- 返回JSON格式的响应
- 包含请求日志记录功能
- 通过Docker容器化部署

## 快速开始

### 本地开发环境

1. 克隆代码库
2. 创建并激活虚拟环境
```bash
python3 -m venv venv
source venv/bin/activate
```

3. 安装依赖
```bash
pip install -r requirements.txt
```

4. 启动服务
```bash
cd app
python main.py
```

### 使用Docker Compose部署

1. 构建并启动容器
```bash
docker-compose up -d
```

2. 停止服务
```bash
docker-compose down
```

## API接口说明

### 主接口 (POST /returnAll)

**请求:**
- 方法: POST
- URL: http://localhost:3006/returnAll
- 可以包含任意查询参数、请求头和请求体

**响应:**
```json
{
  "message": "请求已成功接收并回显",
  "request_info": {
    "request_id": "uuid字符串",
    "timestamp": "ISO格式时间戳",
    "method": "请求方法",
    "url": "完整URL",
    "path_params": "路径参数",
    "query_params": "查询参数",
    "headers": "请求头",
    "client": {
      "host": "客户端IP",
      "port": "客户端端口"
    },
    "body": "请求体内容"
  },
  "log_id": "日志ID"
}
```

## 测试示例

使用curl测试：

```bash
curl -X POST "http://localhost:3006/returnAll?param1=value1&param2=value2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{"name": "测试", "data": {"key": "value"}}'
```

使用Python测试：

```python
import requests
import json

url = "http://localhost:3006/returnAll"
params = {"param1": "value1", "param2": "value2"}
headers = {"Content-Type": "application/json", "Custom-Header": "test"}
data = {"name": "测试", "data": {"key": "value"}}

response = requests.post(url, params=params, headers=headers, json=data)
print(json.dumps(response.json(), indent=2, ensure_ascii=False))
``` 