import pytest
from pydantic import ValidationError

from service.handlers.input_validation import lambda_handler
from tests.utils import generate_context


def test_lambda_handler():
    context = generate_context()
    with pytest.raises(ValidationError):
        lambda_handler({"test": "value"}, context)

    with pytest.raises(ValidationError):
        lambda_handler({"id": "a3", "quantity": 3, "description": "test"}, context)

    assert (
        lambda_handler({"id": 2, "quantity": 3, "description": "test"}, context) is None
    )
