terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
   region = var.zone
}
 
 
module "vpc"  {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.77.0"

  name                 = "terraformVpc"
  cidr               = "10.0.0.0/16"
  azs                  = ["us-east-1a", "us-east-1b", "us-east-1c"]
  public_subnets       = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_security_group" "secGroup" {
  name        = "security"
  description = "Security group for AWS lambda and AWS RDS connection"
  vpc_id      = module.vpc.vpc_id
  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["127.0.0.1/32"]
    self = true
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}

resource "aws_db_subnet_group" "subnetTerraform" {
  name       = "main"
  subnet_ids = module.vpc.public_subnets

  tags = {
    Name = "DB subnet group"
  }
}

#RDS instance
resource "aws_db_instance" "registryDB" {
  identifier = "registry"
  allocated_storage    = 5
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t3.micro"
  name                 = "registry"
  username             = var.db_username
  password             = var.db_password
  db_subnet_group_name = aws_db_subnet_group.subnetTerraform.id
  vpc_security_group_ids = [aws_security_group.secGroup.id]
  final_snapshot_identifier = "dummyid"
  skip_final_snapshot  = true
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_policy" {
  name        = "lambda_policy"
  path        = "/"
  description = "IAM policy for a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
          "rds:*"
      ],
      "Resource": "arn:aws:rds:*:*:*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    },
        {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeNetworkInterfaces",
        "ec2:CreateNetworkInterface",
        "ec2:DeleteNetworkInterface",
        "ec2:DescribeInstances",
        "ec2:AttachNetworkInterface"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "policyattach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

#Setup DB lambda
resource "aws_lambda_function" "setupDB" {
   function_name = "LambdaDB"
   filename      = "rds_setup.zip"

   handler = "app.handler_setup"
   runtime = "python3.8"

   role = aws_iam_role.iam_for_lambda.arn

   vpc_config {
      subnet_ids = module.vpc.public_subnets
      security_group_ids = [aws_security_group.secGroup.id]
  }
  environment {
    variables = {
      rds_endpoint = aws_db_instance.registryDB.endpoint
      db_username = var.db_username
      db_password = var.db_password
      db_name = aws_db_instance.registryDB.name
    }
  }
}

#Reg service lambda
resource "aws_lambda_function" "regService" {
  function_name = "RegService"
  filename = "registry_service.zip"

  handler = "app.handler_registry"
  runtime = "python3.8"

  role = aws_iam_role.iam_for_lambda.arn

  vpc_config {
      subnet_ids = module.vpc.public_subnets
      security_group_ids = [aws_security_group.secGroup.id]
  }
  
  environment {
    variables = {
      rds_endpoint = aws_db_instance.registryDB.endpoint
      db_username = var.db_username
      db_password = var.db_password
      db_name = aws_db_instance.registryDB.name
    }
  } 
}

#DynamoDB for cloud storage
resource "aws_dynamodb_table" "cloud_storage" {
  name = "cloud_storage"
  read_capacity = 20
  write_capacity = 20
  hash_key = "Key"

  attribute {
    name = "Key"
    type = "S"
  }
}

#Client lambda
resource "aws_lambda_function" "allIP" {
  function_name = "ClientFunc"
  filename      = "client.zip"

  handler = "app.client_handler"
  runtime = "python3.8"

  role = aws_iam_role.iam_for_lambda.arn

  vpc_config {
    subnet_ids = module.vpc.public_subnets
    security_group_ids = [aws_security_group.secGroup.id]
  }
  environment {
    variables = {
      rds_endpoint = aws_db_instance.registryDB.endpoint
      db_username = var.db_username
      db_password = var.db_password
      db_name = aws_db_instance.registryDB.name
    }
  }
}

#Gateway
resource "aws_api_gateway_rest_api" "API" {
  name        = "API"
  description = "API gateway"
}

resource "aws_api_gateway_resource" "Resource" {
  rest_api_id = aws_api_gateway_rest_api.API.id
  parent_id   = aws_api_gateway_rest_api.API.root_resource_id
  path_part   = "nodesIP"
}

resource "aws_api_gateway_method" "Method" {
  rest_api_id   = aws_api_gateway_rest_api.API.id
  resource_id   = aws_api_gateway_resource.Resource.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = aws_api_gateway_rest_api.API.id
  resource_id             = aws_api_gateway_resource.Resource.id
  http_method             = aws_api_gateway_method.Method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.allIP.invoke_arn
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.allIP.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.API.execution_arn}/*/*"
}
resource "aws_api_gateway_deployment" "dev" {
  depends_on = [aws_api_gateway_integration.integration]
  rest_api_id = aws_api_gateway_rest_api.API.id
  stage_name = "dev"
}
