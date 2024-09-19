"""
Service stack
"""

import getpass

from aws_cdk import Stack, Tags, CfnParameter
from constructs import Construct

import infrastructure.constants as constants
from infrastructure.input_validation.construct import InputValidationLambdaConstruct
from infrastructure.lambda_layers.construct import LambdaLayerConstruct
from infrastructure.monitoring import MonitoringDashboard
from infrastructure.multi_layers.construct import MultiLayersLambdaConstruct
from infrastructure.simple_lambda.construct import SimpleLambdaConstruct


class StackParameters:
    def __init__(self, scope: Construct):
        self.alarm_email_address = CfnParameter(
            scope,
            "AlarmEmailAddress",
            type="String",
            description="Email for pipeline outcome notifications",
            allowed_pattern="^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$",
            constraint_description="Please enter an email address with correct format (example@example.com)",
            min_length=5,
            max_length=320,
        )


class CdkExamplesStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)
        self._add_stack_tags()
        parameters = StackParameters(self)

        # Add python layers
        lambda_layer = LambdaLayerConstruct(self, f"{construct_id}_lambda_layer")

        # Monitoring dashboard
        cloudwatch_dashboard = MonitoringDashboard(
            self,
            f"{construct_id}{constants.DELIMITER}monitoring",
            "Observation",
            parameters.alarm_email_address,
        )

        # Create a simple lambda function
        SimpleLambdaConstruct(
            self,
            f"{construct_id}{constants.DELIMITER}simple_function",
        )

        # Input validation lambda function
        InputValidationLambdaConstruct(
            self,
            f"{construct_id}{constants.DELIMITER}input_validation",
            lambda_layers=[lambda_layer.input_validation_layer],
            monitoring_dashboard=cloudwatch_dashboard,
        )

        # Multiple lambda layers
        MultiLayersLambdaConstruct(
            self,
            f"{construct_id}{constants.DELIMITER}multi_layers",
            lambda_layers=[lambda_layer.common_layer, lambda_layer.multi_layers],
            monitoring_dashboard=cloudwatch_dashboard,
        )

    def _add_stack_tags(self) -> None:
        # best practice to help identify resources in the console
        Tags.of(self).add(constants.SERVICE_NAME_TAG, constants.SERVICE_NAME)
        Tags.of(self).add(constants.OWNER_TAG, getpass.getuser())
