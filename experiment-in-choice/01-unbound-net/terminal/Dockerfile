# Use a lightweight, hardened Linux base image with Python
FROM python:3.11-slim
# Set the working environment inside the virtual machine
WORKDIR /usr/src/unbound_net
# Prevent Python from writing .pyc files and force unbuffered logging
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
# Install baseline performance and compilation requirements
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*
# Copy the core processing protocols into the workspace
COPY terminal_interface.py .
# Expose the communications array ports for local terminal streaming
EXPOSE 8080
# Execute the runtime matrix simulator on initialization
CMD ["python", "terminal_interface.py"]
