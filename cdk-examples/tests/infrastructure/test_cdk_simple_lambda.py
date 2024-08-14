"""
Test Infrastructure
"""

import aws_cdk as core
import aws_cdk.assertions as assertions

from infrastructure.component import CdkExamplesStack


# example tests. To run these tests, uncomment this file along with the example
# resource in cdk_examples/cdk_examples_stack.py
def test_simple_lambda_created():
    app = core.App()
    stack = CdkExamplesStack(app, "cdk-examples")
    template = assertions.Template.from_stack(stack)

    # Extra function for log rotation
    template.resource_count_is("AWS::Lambda::Function", 2)
