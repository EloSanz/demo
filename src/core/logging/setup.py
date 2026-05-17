import logging
import structlog

def setup_logging(is_dev: bool = True):
    shared_processors = [
        # Merge dynamically bound variables in the context (like request_id)
        structlog.contextvars.merge_contextvars,
        # Add log level (info, warn, error, etc.)
        structlog.processors.add_log_level,
        # Add timestamps in ISO 8601 format
        structlog.processors.TimeStamper(fmt="iso"),
        # Pretty format exceptions
        structlog.processors.format_exc_info,
        # Convert bytes to string to prevent serialization errors
        structlog.processors.UnicodeDecoder(),
    ]

    if is_dev:
        # Interactive and colored renderer for local developer console
        renderer = structlog.dev.ConsoleRenderer(colors=True)
    else:
        # High-speed JSON renderer for Loki/Grafana and production logging
        renderer = structlog.processors.JSONRenderer()

    structlog.configure(
        processors=shared_processors + [renderer],
        context_class=dict,
        logger_factory=structlog.PrintLoggerFactory(),
        wrapper_class=structlog.make_filtering_bound_logger(logging.INFO),
        cache_logger_on_first_use=True,
    )
