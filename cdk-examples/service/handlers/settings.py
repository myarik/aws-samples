"""
Service settings
"""
from aws_lambda_powertools import Logger, Tracer

SERVICE_NAME = "input_validation"

logger: Logger = Logger(service=SERVICE_NAME)

tracer: Tracer = Tracer(service=SERVICE_NAME)