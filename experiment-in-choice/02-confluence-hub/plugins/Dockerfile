FROM python:3.11-slim
WORKDIR /app
RUN pip install websockets
COPY agent_plugin.py .
CMD ["python", "agent_plugin.py"]
