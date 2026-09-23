import time
import random
from prometheus_client import start_http_server, Counter, Gauge
# Define internal tracking metrics
SYSTEM_HEARTBEATS = Counter('matrix_heartbeat_total', 'Total loops executed by the simulation node')
SYSTEM_INSTABILITY = Counter('matrix_anomaly_total', 'Total high-frequency drift events detected')
CORE_RESONANCE = Gauge('matrix_resonance_hz', 'Current fundamental frequency modulation state')
def run_telemetry_node():
    # Start Prometheus scraper port inside the container
    start_http_server(8080)
    print("[TELEMETRY] Metrics exporter running live on port :8080/metrics")
    
    base_resonance = 32.0
    
    while True:
        SYSTEM_HEARTBEATS.inc()
        
        # Calculate artificial drift metrics
        instability_triggered = random.random() < 0.08
        if instability_triggered:
            SYSTEM_INSTABILITY.inc()
            current_resonance = base_resonance + random.uniform(12.0, 45.0)
        else:
            current_resonance = base_resonance + random.uniform(-2.0, 2.0)
            
        CORE_RESONANCE.set(current_resonance)
        time.sleep(1.0)
if __name__ == '__main__':
    run_telemetry_node()
