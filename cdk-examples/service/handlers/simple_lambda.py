"""
Simple lambda handler
"""

from typing import Any


def lambda_handler(event: dict[str, Any], context: Any) -> dict:
    """
    Simple lambda handler
    """
    return {"statusCode": 200, "body": "Hello from Lambda!"}
