FROM python:3.12

WORKDIR /app

COPY pyproject.toml ./

RUN pip install uv

RUN uv sync

COPY . .

CMD ["uv", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8085", "--forwarded-allow-ips=*"]
