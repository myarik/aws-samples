"""
Simple lambda handler
"""

import requests
from aws_lambda_powertools.utilities.parser import event_parser
from aws_lambda_powertools.utilities.typing import LambdaContext

from service.handlers.settings import logger, tracer
from service.models.input_event import ServiceUrlInput


@tracer.capture_method(capture_response=False)
def get_data_from_api(url: str) -> requests.Response:
    """
    Fetch data from the specified API URL.

    Args:
        url (str): The API endpoint to fetch data from.

    Returns:
        requests.Response: The response object from the API call.
    """
    try:
        response = requests.get(url)
        response.raise_for_status()
        return response
    except requests.exceptions.RequestException as e:
        logger.error("Failed to get data from API", extra={"error": str(e)})
        raise


# Documentation -- https://docs.powertools.aws.dev/lambda/python/latest/utilities/parser/#built-in-models
@event_parser(model=ServiceUrlInput)
@tracer.capture_lambda_handler(capture_response=False)
def lambda_handler(event: ServiceUrlInput, context: LambdaContext) -> None:
    """
    Simple lambda handler
    """
    # Request data from API
    response = get_data_from_api(event.url)
    logger.info("API returns data", extra={"data": response.json()})
    return response.json()
