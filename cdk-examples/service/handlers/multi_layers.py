"""
Simple lambda function with multiple layers
"""

import os

from aws_lambda_powertools import Metrics
from aws_lambda_powertools.utilities.parser import event_parser
from aws_lambda_powertools.utilities.typing import LambdaContext
from jinja2 import Environment, DictLoader

from service.handlers.settings import logger, SERVICE_NAME
from service.models.input_event import JinjaTemplateInput

# Create Jinja2 environment
env = Environment(
    loader=DictLoader(
        {
            "example_template": """
    <html>
        <head><title>{{ title }}</title></head>
        <body>
            <h1>{{ title }}</h1>
            <p>{{ content }}</p>
        </body>
    </html>
    """
        }
    )
)

# Metrics
metrics = Metrics(
    namespace=SERVICE_NAME,
    service=f"{SERVICE_NAME}--{os.getenv("ENVIRONMENT", "dev")}--multi_layers",
)


@event_parser(model=JinjaTemplateInput)
@metrics.log_metrics
def lambda_handler(data: JinjaTemplateInput, context: LambdaContext) -> str:
    """
    Simple lambda handler
    """
    template = env.get_template("example_template")
    rendered_output = template.render(data.model_dump())
    logger.info("Rendered output: %s", rendered_output)
    metrics.add_metric(name="DemoMetric", unit="Count", value=1)
    return rendered_output
