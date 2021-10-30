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
