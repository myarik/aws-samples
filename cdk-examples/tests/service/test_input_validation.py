from unittest.mock import patch

import pytest
from pydantic import ValidationError

from service.handlers.input_validation import lambda_handler
from tests.utils import generate_context


def test_lambda_handler():
    context = generate_context()
    with pytest.raises(ValidationError):
        lambda_handler({"test": "value"}, context)

    with pytest.raises(ValidationError):
        lambda_handler({"url": "test"}, context)

    with patch("requests.get") as mock_get:
        mock_get.return_value.json.return_value = {"test": "value"}
        assert lambda_handler({"url": "https://api.github.com"}, context) == {
            "test": "value"
        }
