import pytest
from pydantic import ValidationError

from service.handlers.multi_layers import lambda_handler
from tests.utils import generate_context


def test_lambda_handler():
    context = generate_context()
    with pytest.raises(ValidationError):
        lambda_handler({"test": "value"}, context)
    output = lambda_handler({"title": "Test", "content": "Hello World"}, context)
    assert "<html>" in output
    assert "<title>Test</title>" in output
    assert "<h1>Test</h1>" in output
    assert "<p>Hello World</p>" in output
