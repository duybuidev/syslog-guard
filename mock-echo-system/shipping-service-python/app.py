from fastapi import FastAPI, Response, status
import logging
import sys
import random
import os

# Cấu hình log đẩy ra Stdout cho Docker thu thập
logging.basicConfig(stream=sys.stdout, level=logging.INFO, 
                    format='[%(levelname)s] %(asctime)s - %(message)s')

app = FastAPI(title="Shipping Service", version="1.0.0")

@app.get("/health", status_code=status.HTTP_200_OK)
async def health_check():
    return {"status": "Shipping Service is healthy"}

@app.get("/ship")
async def ship_order():
    shipping_status = random.choice(['Shipped', 'Pending', 'Delayed'])
    
    if shipping_status == 'Delayed':
        logging.warning(f"Shipment delayed for user {random.randint(1, 100)}")
    else:
        logging.info(f"Shipment status: {shipping_status}")
        
    return {"shipping_status": shipping_status}

@app.get("/crash")
async def crash_system():
    logging.critical("CRITICAL: Memory leak detected. System halting immediately.")
    os._exit(1) # Lệnh này giết chết process ngay lập tức, bỏ qua mọi exception
