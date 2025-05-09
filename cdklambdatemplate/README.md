# AWS CDK Go Demo: Lambda Template

This project serves as a template for creating and deploying AWS Lambda functions written in Go, using the AWS Cloud
Development Kit (CDK) for infrastructure management.

## Prerequisites

Before you begin, ensure you have the following installed:

* Go (version 1.23 or higher recommended, as specified in `services/demo/go.mod`)
* AWS CLI configured with appropriate credentials
* GNU Make

## Project Structure

* `infrastructure/`: Contains the AWS CDK code (in Go) for defining the Lambda function and related AWS resources.
* `services/demo/`: Contains the Go source code for the Lambda function.
    * `cmd/main.go`: The entry point for the Lambda function.
    * `handler/`: Business logic for the Lambda.
    * `pkg/`: Shared packages, like the logger.
* `Makefile`: Contains helper commands for building, deploying, testing, and managing the Lambda.
* `.build/`: Directory where the compiled Lambda binary is stored (created during the build process).

## Setup and Installation

**Clone the repository:**

```bash
git clone <your-repository-url>
cd <repository-name>
```

**Install CDK and project dependencies:**

```bash
go mod tidy
cd ..
```

**Install Go dependencies for the Lambda function:**

```bash
cd services/demo
go mod tidy
 cd ../..
```

## Development Workflow

### Building the Lambda

To compile the Go Lambda function:

```bash
make build
```

This command compiles the Lambda function defined in `services/demo/cmd/main.go` and places the binary in the
`.build/lambdas/demo/bootstrap` directory.

### Deploying the Infrastructure

To deploy the Lambda function and associated AWS resources using AWS CDK:

```bash
make deploy
```

This will synthesize the CloudFormation template from the CDK code in `infrastructure/app.go` and deploy it to your AWS
account. The Lambda function name will be outputted upon successful deployment.

### Testing

To run unit tests for the Lambda function:

```bash
make test
```

This command executes tests located within the `services/demo` directory.

### Linting

To run linters on both the service and infrastructure code:

```bash
make lint
```

### Viewing Lambda Logs

To view logs for the deployed Lambda function:

```bash
make logs <lambda-name>
```

Replace `<lambda-name>` with the actual name of your Lambda function (e.g., `LambdaStack-Demo`). This command fetches
logs from the last 24 hours.

### Destroying the Infrastructure

To remove all AWS resources created by this stack:

```bash
make destroy
```

This will tear down the CloudFormation stack and delete associated resources.

## Lambda Configuration

The Lambda function's behavior can be configured via environment variables, which are set in the `infrastructure/app.go`
file. Key environment variables include:

* `Environment`: Specifies the application environment (e.g., "staging", "production").
* `LOG_LEVEL`: Sets the logging level (e.g., "info", "debug", "error").
* `APP_VERSION`: Application version.

The logger (`services/demo/pkg/logger/logger.go`) is configured based on these environment variables. In "staging" or "
development" environments, it uses a colored, human-readable format. In other environments (like "production"), it
defaults to JSON formatted logs.

## Makefile Commands

The `Makefile` provides several useful commands:

* `help`: Displays a help screen with available targets.
* `lint`: Runs linters for Go code in `services/demo` and `infrastructure`.
* `test`: Runs unit tests for the `services/demo` code.
* `logs <lambda-name>`: Fetches logs for the specified Lambda function.
* `build`: Compiles the Lambda function.
* `deploy`: Deploys the infrastructure using AWS CDK.
* `destroy`: Destroys the deployed infrastructure.

Refer to the `Makefile` for more details on each command.
