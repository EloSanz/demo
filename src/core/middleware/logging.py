import uuid
import time
import structlog

from fastapi import Request

logger = structlog.get_logger()


async def logging_middleware(request: Request, call_next):
    if request.url.path in ("/metrics", "/docs", "/openapi.json", "/redoc"):
        return await call_next(request)

    request_id = str(uuid.uuid4())
    start_time = time.time()

    with structlog.contextvars.bound_contextvars(request_id=request_id):
        logger.info("request_started", method=request.method, path=request.url.path)
        
        response = await call_next(request)

        process_time = time.time() - start_time
        logger.info(
            "request_finished",
            duration=process_time,
            status=response.status_code
        )

        response.headers["X-Request-ID"] = request_id
        return response
