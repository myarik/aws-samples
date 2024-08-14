"""
Python constructs for Lambda Layers
"""

from aws_cdk import aws_lambda as _lambda, RemovalPolicy
from aws_cdk.aws_lambda_python_alpha import PythonLayerVersion
from constructs import Construct

import infrastructure.constants as constants


class LambdaLayerConstruct(Construct):
    def __init__(self, scope: Construct, construct_id: str) -> None:
        super().__init__(scope, construct_id)
        self.construct_id = construct_id
        self.common_layer = self._build_common_layer()

    def _build_common_layer(self) -> PythonLayerVersion:
        """
        Build common layer
        """
        return PythonLayerVersion(
            self,
            f"{self.construct_id}_common",
            entry=constants.COMMON_LAYER_BUILD_FOLDER,
            compatible_runtimes=[_lambda.Runtime.PYTHON_3_12],
            removal_policy=RemovalPolicy.DESTROY,
        )
