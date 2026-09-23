import os
import sys
import time
import math
import random
def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')
def render_matrix_pulse(cycle_count, distortion=False):
    """Generates wireframe corridor slice approximations based on active state."""
    width = 50
    # Base Lorenz Attractor scale representation
    sigma = 45.0 if distortion else 10.0
    
    # Calculate wave dynamics
    frequency_modifier = math.sin(cycle_count * 0.4) * sigma
    offset = int((width / 2) + (frequency_modifier * 0.3))
    
    # Apply rendering parameters
    line_char = "█" if distortion else "░"
    color_code = "\033[1;31m" if distortion else "\033[1;32m"
    reset_code = "\033[0m"
    
    # Build string slice
    buffer_str = [" "] * width
    left_edge = max(0, min(width - 1, offset - 5))
    right_edge = max(0, min(width - 1, offset + 5))
    
    buffer_str[left_edge] = line_char
    buffer_str[right_edge] = line_char
    
    if distortion:
        # Inject randomized corruption artifacts into the frame buffer
        for _ in range(3):
            buffer_str[random.randint(0, width - 1)] = "×"
            
    print(f"{color_code}{''.join(buffer_str)}{reset_code}")
def main_loop():
    clear_screen()
    print("\033[1;32m[SYSTEM ONLINE] TERMINAL INTERFACE STAGE ACTIVE\033[0m")
    print("\033[0;37mRunning real-time simulation layer... Press Ctrl+C to terminate.\033[0m\n")
    time.sleep(1)
    
    cycle = 0
    try:
        while True:
            # Simulate a 5% random system instability spike (simulating key-strike triggers)
            instability_spike = random.random() < 0.08
            
            if instability_spike:
                print(f"\n\033[1;31m[!] CRITICAL DRIFT // SIGMA SCALE OVERDRIVE\033[0m")
                for _ in range(5):  # Flash high-frequency bursts
                    render_matrix_pulse(cycle, distortion=True)
                    time.sleep(0.05)
                print("")
            else:
                render_matrix_pulse(cycle, distortion=False)
                
            cycle += 1
            time.sleep(0.12)
            
    except KeyboardInterrupt:
        print("\n\033[1;33m[SHUTDOWN] Interactive terminal process exited safely.\033[0m")
if __name__ == "__main__":
    main_loop()
