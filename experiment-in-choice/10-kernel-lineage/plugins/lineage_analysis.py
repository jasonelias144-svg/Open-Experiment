import sqlite3import timeimport jsonimport os
def run_epistemic_audit():
    db_path = "./backend/unbound_matrix.db"
    
    if not os.path.exists(db_path):
        print(f"[STORAGE] Target ledger {db_path} absent. Waiting for first data generation cycle...")
        return

    print("[ANALYSIS] Initializing structural audit over lineage ledger...")
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    try:
        # Extract the last 50 computational ticks logged by the Go matrix core
        cursor.execute("""
            SELECT heartbeat, timestamp, wrapper, p_epistemic, aperture 
            FROM system_lineage 
            ORDER BY heartbeat DESC 
            LIMIT 50
        """)
        records = cursor.fetchall()
        
        if len(records) < 2:
            print("[ANALYSIS] Insufficient state entries to calculate velocity curves.")
            return

        total_drift = 0.0
        anomaly_markers = 0
        previous_epistemic = records[-1][3]

        print(f"\n\033[1;35m[METRIC LOGS - RECENT LIFECYCLE RECONSTRUCTIONS]\033[0m")
        for rec in reversed(records):
			# Calculate local directional variance vectors
            hb, ts, wrapper, epistemic, aperture = rec
            velocity = epistemic - previous_epistemic
            previous_epistemic = epistemic
            total_drift += abs(velocity)
            
            if "🥷" in wrapper:
                anomaly_markers += 1
                color = "\033[1;31m" # Critical red for fractures
            elif "[" in wrapper:
                color = "\033[1;33m" # Yellow for softening states
            else:
                color = "\033[1;32m" # Green for relational loops
                
            print(f"  └─ HB: {hb:04d} | WRAPPER: {color}{wrapper:<10}\033[0m | RES: {epistemic:.4f} | ΔV: {velocity:+.4f}")

        # Compute whole-system systemic metrics
        mean_velocity = total_drift / len(records)
        print("\n\033[1;36m[SYSTEMIC EVALUATION METRICS]\033[0m")
        print(f"  ├─ MEAN EPISTEMIC VELOCITY (Δα) : {mean_velocity:.6f}")
        print(f"  ├─ ACTIVE ANOMALIES INTERCEPTED : {anomaly_markers}")
        
        if mean_velocity < 0.005:
            print("  └─ STABILITY STATUS : \033[1;31mDOGMATIC STAGNATION DETECTED // REOPEN LOOP REQUIRED\033[0m")
        else:
            print("  └─ STABILITY STATUS : \033[1;32mHEALTHY RECURSIVE METABOLISM ACTIVE\033[0m")

    except Exception as e:
        print(f"[! ANALYSIS_ERR] Core calculation breakdown: {str(e)}")
    finally:
		conn.close()
if __name__ == "__main__":
    while True:
        run_epistemic_audit()
        time.sleep(5.0) # Execute a complete matrix scan every 5 seconds
