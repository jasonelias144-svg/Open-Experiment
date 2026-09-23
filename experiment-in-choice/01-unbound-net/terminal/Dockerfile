# Use a lightweight, hardened Linux base image with PythonFROM python:3.11-slim
# Set the working environment inside the virtual machineWORKDIR /usr/src/unbound_net
# Prevent Python from writing .pyc files and force unbuffered loggingENV PYTHONDONTWRITEBYTECODE=1ENV PYTHONUNBUFFERED=1
# Install baseline performance and compilation requirementsRUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*
# Copy the core processing protocols into the workspaceCOPY terminal_interface.py .
# Expose the communications array ports for local terminal streamingEXPOSE 8080
# Execute the runtime matrix simulator on initializationCMD ["python", "terminal_interface.py"]
