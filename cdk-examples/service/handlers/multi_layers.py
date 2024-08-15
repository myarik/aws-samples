"""
Simple lambda function with multiple layers
"""
from aws_lambda_powertools.utilities.parser import event_parser
from aws_lambda_powertools.utilities.typing import LambdaContext
from jinja2 import Environment, DictLoader

from service.handlers.settings import logger
from service.models.input_event import JinjaTemplateInput

env = Environment(loader=DictLoader({
    'example_template': '''
    <html>
        <head><title>{{ title }}</title></head>
        <body>
            <h1>{{ title }}</h1>
            <p>{{ content }}</p>
        </body>
    </html>
    '''
}))


@event_parser(model=JinjaTemplateInput)
def lambda_handler(data: JinjaTemplateInput, context: LambdaContext) -> str:
    """
    Simple lambda handler
    """
    template = env.get_template('example_template')
    rendered_output = template.render(data.model_dump())
    logger.info("Rendered output: %s", rendered_output)
    return rendered_output


