FROM python:3.11-slim
WORKDIR /app
# Pre-install clean database binding tools if needed
RUN pip install --no-cache-dir websockets
COPY lineage_analysis.py .
CMD ["python", "lineage_analysis.py"]
