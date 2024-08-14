from service.handlers.simple_lambda import lambda_handler
from tests.utils import generate_context

def test_lambda_handler():

    event = {}
    context = generate_context()
    assert lambda_handler(event, context) is None
