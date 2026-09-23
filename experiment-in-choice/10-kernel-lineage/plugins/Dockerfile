FROM python:3.11-slimWORKDIR /app# Pre-install clean database binding tools if neededRUN pip install --no-cache-dir websocketsCOPY lineage_analysis.py .CMD ["python", "lineage_analysis.py"]
