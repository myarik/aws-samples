"""
Simple lambda handler
"""
from typing import Any

from aws_lambda_powertools.utilities.typing import LambdaContext

from service.handlers.settings import logger


def lambda_handler(event: dict[str, Any], context: LambdaContext) -> None:
    """
    Simple lambda handler
    """
    logger.info("Hello, world!", extra={"event": event})
    print("Hello, world!")
