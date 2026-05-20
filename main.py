import os

from fastapi import FastAPI
from prometheus_client import Info
from prometheus_fastapi_instrumentator import Instrumentator
from starlette.middleware.base import BaseHTTPMiddleware

from src.core.logging.setup import setup_logging
from src.core.middleware.logging import logging_middleware, logger
from src.postback.handlers import setup_postback_exception_handlers
from src.postback.router import postback_router, eac_router

is_dev = os.getenv("ENVIRONMENT", "development") != "production"
setup_logging(is_dev=is_dev)

# Expose app metadata metric for Grafana dashboards (like ID 18739)
app_info = Info("fastapi_app_info", "FastAPI application information")
app_info.info({
    "app_name": "fastapi-app",
    "app_version": "1.0.0"
})

app = FastAPI()

# Instrument for Prometheus metrics
Instrumentator().instrument(app).expose(app)

# Add Middlewares
app.add_middleware(BaseHTTPMiddleware, dispatch=logging_middleware)

# Setup Exception Handlers (Modularized)
setup_postback_exception_handlers(app)

# Include Routers
app.include_router(postback_router)
app.include_router(eac_router)

@app.get("/test")
def test_endpoint():
    return "test"



