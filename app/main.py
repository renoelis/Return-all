from fastapi import FastAPI, Request, Depends
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import uuid
import json
import datetime
import re
from pathlib import Path
from typing import Dict, Any, List, Tuple, Optional

# 导入自定义日志模块
from app.utils.logger import setup_logger

# 配置日志
log_dir = Path("logs")
log_dir.mkdir(exist_ok=True)
logger = setup_logger("returnAll-api", "logs/api.log")

app = FastAPI(title="returnAll API", description="将请求内容原样返回的API服务")

# 允许跨域请求
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

async def get_request_body(request: Request) -> bytes:
    """获取请求体内容"""
    return await request.body()

def parse_json_body(body: bytes) -> Tuple[Any, bool, Optional[str], Optional[str], Optional[Dict]]:
    """
    解析请求体，尝试转换为JSON
    
    Args:
        body: 请求体字节数据
        
    Returns:
        Tuple[Any, bool, Optional[str], Optional[str], Optional[Dict]]: 
            - 请求体内容(JSON对象或紧凑字符串)
            - 是否为有效JSON
            - JSON解析错误信息
            - 原始请求体字符串(保留换行符)
            - 错误细节信息
    """
    if not body:
        return None, True, None, None, None
    
    # 保留原始字符串，包括换行符
    original_str = body.decode("utf-8", errors="replace")
    
    try:
        json_data = json.loads(body)
        return json_data, True, None, original_str, None
    except json.JSONDecodeError as e:
        # 如果请求体不是有效的JSON，则以紧凑字符串形式返回
        body_str = body.decode("utf-8", errors="replace")
        # 删除所有换行符和多余空格，保持紧凑格式
        compact_str = re.sub(r'\s+', ' ', body_str).strip()
        
        # 构建错误信息，包含具体位置
        error_msg = f"JSON解析错误: {str(e)}"
        
        # 提取错误位置详细信息
        error_details = {
            "error_type": "JSONDecodeError",
            "message": str(e),
            "line": e.lineno,
            "column": e.colno,
            "position": e.pos,
            "error_char": original_str[e.pos] if e.pos < len(original_str) else ""
        }
        
        # 提取错误行的内容
        lines = original_str.splitlines()
        if e.lineno <= len(lines) and e.lineno > 0:
            error_line = lines[e.lineno - 1]
            error_details["line_content"] = error_line
            
            # 创建指向错误位置的指示器
            pointer = " " * (e.colno - 1) + "^"
            error_details["pointer"] = pointer
        
        return compact_str, False, error_msg, original_str, error_details

@app.post("/returnAll")
async def returnAll_request(request: Request, body: bytes = Depends(get_request_body)):
    """返回所有请求信息，包括参数、请求头和请求体"""
    # 生成请求ID
    request_id = str(uuid.uuid4())
    
    # 获取请求信息
    request_info = {
        "request_id": request_id,
        "timestamp": datetime.datetime.now().isoformat(),
        "method": request.method,
        "url": str(request.url),
        "path_params": request.path_params,
        "query_params": dict(request.query_params),
        "headers": dict(request.headers),
        "client": {
            "host": request.client.host if request.client else None,
            "port": request.client.port if request.client else None
        }
    }
    
    # 解析请求体
    body_content, is_valid_json, error_msg, original_str, error_details = parse_json_body(body)
    request_info["body"] = body_content
    request_info["is_valid_json"] = is_valid_json
    if not is_valid_json and error_msg:
        request_info["json_error"] = error_msg
        request_info["original_body"] = original_str
        request_info["error_details"] = error_details
    
    # 记录日志
    logger.info(f"Received request: {request_id}", extra={"request_info": request_info})
    
    return {
        "message": "success",
        "request_info": request_info,
        "log_id": request_id
    }

@app.post("/returnAll/{path_param}")
async def returnAll_with_path(path_param: str, request: Request, body: bytes = Depends(get_request_body)):
    """支持路径参数的returnAll API"""
    # 生成请求ID
    request_id = str(uuid.uuid4())
    
    # 获取请求信息
    request_info = {
        "request_id": request_id,
        "timestamp": datetime.datetime.now().isoformat(),
        "method": request.method,
        "url": str(request.url),
        "path_params": {
            **request.path_params,
            "path_param": path_param  # 额外添加路径参数
        },
        "query_params": dict(request.query_params),
        "headers": dict(request.headers),
        "client": {
            "host": request.client.host if request.client else None,
            "port": request.client.port if request.client else None
        }
    }
    
    # 解析请求体
    body_content, is_valid_json, error_msg, original_str, error_details = parse_json_body(body)
    request_info["body"] = body_content
    request_info["is_valid_json"] = is_valid_json
    if not is_valid_json and error_msg:
        request_info["json_error"] = error_msg
        request_info["original_body"] = original_str
        request_info["error_details"] = error_details
    
    # 记录日志
    logger.info(f"Received path request: {request_id}", extra={"request_info": request_info})
    
    return {
        "message": "success",
        "request_info": request_info,
        "log_id": request_id
    }

@app.post("/returnAll/{param1}/{param2}")
async def returnAll_with_two_params(param1: str, param2: str, request: Request, body: bytes = Depends(get_request_body)):
    """支持两个路径参数的returnAll API"""
    # 生成请求ID
    request_id = str(uuid.uuid4())
    
    # 获取请求信息
    request_info = {
        "request_id": request_id,
        "timestamp": datetime.datetime.now().isoformat(),
        "method": request.method,
        "url": str(request.url),
        "path_params": {
            **request.path_params,
            "param1": param1,
            "param2": param2
        },
        "query_params": dict(request.query_params),
        "headers": dict(request.headers),
        "client": {
            "host": request.client.host if request.client else None,
            "port": request.client.port if request.client else None
        }
    }
    
    # 解析请求体
    body_content, is_valid_json, error_msg, original_str, error_details = parse_json_body(body)
    request_info["body"] = body_content
    request_info["is_valid_json"] = is_valid_json
    if not is_valid_json and error_msg:
        request_info["json_error"] = error_msg
        request_info["original_body"] = original_str
        request_info["error_details"] = error_details
    
    # 记录日志
    logger.info(f"Received two-param request: {request_id}", extra={"request_info": request_info})
    
    return {
        "message": "success",
        "request_info": request_info,
        "log_id": request_id
    }

@app.post("/returnAll/{param1}/{param2}/{param3}")
async def returnAll_with_three_params(param1: str, param2: str, param3: str, request: Request, body: bytes = Depends(get_request_body)):
    """支持三个路径参数的returnAll API"""
    # 生成请求ID
    request_id = str(uuid.uuid4())
    
    # 获取请求信息
    request_info = {
        "request_id": request_id,
        "timestamp": datetime.datetime.now().isoformat(),
        "method": request.method,
        "url": str(request.url),
        "path_params": {
            **request.path_params,
            "param1": param1,
            "param2": param2,
            "param3": param3
        },
        "query_params": dict(request.query_params),
        "headers": dict(request.headers),
        "client": {
            "host": request.client.host if request.client else None,
            "port": request.client.port if request.client else None
        }
    }
    
    # 解析请求体
    body_content, is_valid_json, error_msg, original_str, error_details = parse_json_body(body)
    request_info["body"] = body_content
    request_info["is_valid_json"] = is_valid_json
    if not is_valid_json and error_msg:
        request_info["json_error"] = error_msg
        request_info["original_body"] = original_str
        request_info["error_details"] = error_details
    
    # 记录日志
    logger.info(f"Received three-param request: {request_id}", extra={"request_info": request_info})
    
    return {
        "message": "success",
        "request_info": request_info,
        "log_id": request_id
    }

@app.post("/returnAll/{path:path}")
async def returnAll_with_any_path(path: str, request: Request, body: bytes = Depends(get_request_body)):
    """支持任意多级路径的returnAll API"""
    # 生成请求ID
    request_id = str(uuid.uuid4())
    
    # 将路径拆分成多个部分
    path_parts = path.split("/")
    path_dict = {f"path_part_{i}": part for i, part in enumerate(path_parts)}
    
    # 获取请求信息
    request_info = {
        "request_id": request_id,
        "timestamp": datetime.datetime.now().isoformat(),
        "method": request.method,
        "url": str(request.url),
        "path": path,  # 完整路径
        "path_parts": path_parts,  # 路径各部分
        "path_params": {
            **request.path_params,
            **path_dict
        },
        "query_params": dict(request.query_params),
        "headers": dict(request.headers),
        "client": {
            "host": request.client.host if request.client else None,
            "port": request.client.port if request.client else None
        }
    }
    
    # 解析请求体
    body_content, is_valid_json, error_msg, original_str, error_details = parse_json_body(body)
    request_info["body"] = body_content
    request_info["is_valid_json"] = is_valid_json
    if not is_valid_json and error_msg:
        request_info["json_error"] = error_msg
        request_info["original_body"] = original_str
        request_info["error_details"] = error_details
    
    # 记录日志
    logger.info(f"Received any-path request: {request_id}", extra={"request_info": request_info})
    
    return {
        "message": "success",
        "request_info": request_info,
        "log_id": request_id
    }

@app.get("/")
async def root():
    """API根路径，返回简单信息"""
    return {
        "message": "returnAll API服务正在运行",
        "endpoints": {
            "returnAll": "/returnAll (POST) - 返回请求的所有内容",
            "returnAll_with_path": "/returnAll/{path_param} (POST) - 支持单个路径参数的returnAll API",
            "returnAll_with_two_params": "/returnAll/{param1}/{param2} (POST) - 支持两个路径参数的returnAll API",
            "returnAll_with_three_params": "/returnAll/{param1}/{param2}/{param3} (POST) - 支持三个路径参数的returnAll API",
            "returnAll_with_any_path": "/returnAll/{path:path} (POST) - 支持任意多级路径的returnAll API"
        },
        "version": "1.0.0"
    }

if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=3006, reload=True) 