"""
Service stack
"""

from aws_cdk import (
    Stack,
    Tags,
)
from constructs import Construct

import infrastructure.constants as constants
from infrastructure.input_validation.construct import InputValidationLambdaConstruct
from infrastructure.lambda_layers.construct import LambdaLayerConstruct
from infrastructure.simple_lambda.construct import SimpleLambdaConstruct


class CdkExamplesStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)
        self._add_stack_tags()
        # Add python layers
        lambda_layer = LambdaLayerConstruct(self, f"{construct_id}_lambda_layer")

        # Create a simple lambda function
        SimpleLambdaConstruct(
            self,
            f"{construct_id}_simple_lambda",
            lambda_layers=[lambda_layer.common_layer],
        )

        # Input validation lambda function
        InputValidationLambdaConstruct(
            self,
            f"{construct_id}_input_validation_lambda",
            lambda_layers=[lambda_layer.input_validation_layer],
        )

    def _add_stack_tags(self) -> None:
        # best practice to help identify resources in the console
        Tags.of(self).add(constants.SERVICE_NAME_TAG, constants.SERVICE_NAME)
        Tags.of(self).add(constants.OWNER_TAG, "myarik")
