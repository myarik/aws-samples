"""
Simple lambda handler
"""

import requests
from aws_lambda_powertools.utilities.parser import event_parser
from aws_lambda_powertools.utilities.typing import LambdaContext

from service.handlers.settings import logger, tracer
from service.models.input_event import Order


@tracer.capture_method(capture_response=False)
def get_data_from_api():
    try:
        response = requests.get("https://api.github.com")
        return response
    except requests.exceptions.RequestException as e:
        logger.error(e)
        return None


# Documentation -- https://docs.powertools.aws.dev/lambda/python/latest/utilities/parser/#built-in-models
@event_parser(model=Order)
@tracer.capture_lambda_handler(capture_response=False)
def lambda_handler(event: Order, context: LambdaContext) -> None:
    """
    Simple lambda handler
    """
    # Request data from API
    get_data_from_api()
