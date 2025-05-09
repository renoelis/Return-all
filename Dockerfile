FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# 创建日志目录
RUN mkdir -p logs

# 暴露3006端口
EXPOSE 3006

# 使用Gunicorn作为生产服务器 + UvicornWorker作为工作进程
CMD ["gunicorn", "app.main:app", "--workers", "4", "--worker-class", "uvicorn.workers.UvicornWorker", "--bind", "0.0.0.0:3006", "--access-logfile", "-"] 