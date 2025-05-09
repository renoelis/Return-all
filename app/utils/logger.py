import logging
import json
from pathlib import Path
from typing import Dict, Any

class CustomJSONFormatter(logging.Formatter):
    """
    自定义JSON格式的日志处理器
    将日志记录格式化为JSON格式
    """
    def format(self, record):
        log_record = {
            "time": self.formatTime(record, self.datefmt),
            "name": record.name,
            "level": record.levelname,
            "message": record.getMessage(),
        }
        
        # 添加额外信息
        if hasattr(record, 'request_info'):
            log_record["request_info"] = record.request_info
            
        # 添加异常信息
        if record.exc_info:
            log_record["exception"] = self.formatException(record.exc_info)
            
        return json.dumps(log_record, ensure_ascii=False)

def setup_logger(name: str, log_file: str = None, level=logging.INFO) -> logging.Logger:
    """
    设置日志记录器
    
    Args:
        name: 日志记录器名称
        log_file: 日志文件路径，如果为None则只输出到控制台
        level: 日志级别
        
    Returns:
        logging.Logger: 配置好的日志记录器
    """
    logger = logging.getLogger(name)
    logger.setLevel(level)
    
    # 创建控制台处理器
    console_handler = logging.StreamHandler()
    console_handler.setLevel(level)
    console_formatter = CustomJSONFormatter()
    console_handler.setFormatter(console_formatter)
    logger.addHandler(console_handler)
    
    # 如果提供了日志文件路径，则创建文件处理器
    if log_file:
        # 确保日志目录存在
        log_dir = Path(log_file).parent
        log_dir.mkdir(exist_ok=True, parents=True)
        
        file_handler = logging.FileHandler(log_file)
        file_handler.setLevel(level)
        file_formatter = CustomJSONFormatter()
        file_handler.setFormatter(file_formatter)
        logger.addHandler(file_handler)
    
    return logger 