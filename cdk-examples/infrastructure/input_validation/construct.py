"""
Input Validation Lambda Construct
"""

from aws_cdk import aws_lambda as _lambda, Duration
from aws_cdk.aws_logs import RetentionDays
from constructs import Construct

import infrastructure.constants as constants


class InputValidationLambdaConstruct(Construct):
    def __init__(self, scope: Construct, construct_id: str, lambda_layers=None) -> None:
        super().__init__(scope, construct_id)
        self.construct_id = construct_id
        self.lambda_layers = [] if lambda_layers is None else lambda_layers
        self.lambda_function = self._build_labmda_function()

    def _build_labmda_function(
        self,
    ) -> _lambda.Function:
        """
        Build lambda function
        """
        return _lambda.Function(
            self,
            "InputValidationLambdaFunction",
            function_name=self.construct_id,
            runtime=_lambda.Runtime.PYTHON_3_12,
            code=_lambda.Code.from_asset(constants.BUILD_FOLDER),
            handler="service.handlers.input_validation.lambda_handler",
            environment={
                "POWERTOOLS_SERVICE_NAME": constants.SERVICE_NAME,  # for logger, tracer and metrics
                "POWERTOOLS_TRACE_DISABLED": "true",  # for tracer
                "LOG_LEVEL": constants.LOG_LEVEL,  # for logger
            },
            tracing=_lambda.Tracing.DISABLED,
            retry_attempts=0,
            timeout=Duration.seconds(constants.HANDLER_LAMBDA_TIMEOUT),
            memory_size=constants.HANDLER_LAMBDA_MEMORY_SIZE,
            layers=self.lambda_layers,
            # role=role,
            log_retention=RetentionDays.ONE_DAY,
            log_format=_lambda.LogFormat.JSON.value,
            system_log_level=_lambda.SystemLogLevel.INFO.value,
        )
