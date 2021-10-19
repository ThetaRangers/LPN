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

resource "aws_lambda_function" "replicaSet" {
   function_name = var.function_name
   filename      = var.lambda_zip

   # "lambda_sqs_consumer" is the filename within the zip file and "my_handler"
   # is the name of the handler function
   handler = var.handler
   runtime = "python3.8"

   role = aws_iam_role.iam_for_lambda.arn
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

# IAM role which dictates what other AWS services the Lambda function
 # may access.

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
          "dynamodb:BatchGetItem",
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchWriteItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem"
      ],
      "Resource": "arn:aws:dynamodb:*:*:*",
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
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "policyattach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}


resource "aws_dynamodb_table" "cloud_storage" {
  name = var.dynamo_table_name
  read_capacity = 20
  write_capacity = 20
  hash_key = "Key"

  attribute {
    name = "Key"
    type = "S"
  }
}

resource "aws_dynamodb_table" "tabellone" {
  name = var.dynamo_table_name1
  read_capacity = 20
  write_capacity = 20
  hash_key = "Key"

  attribute {
    name = "Key"
    type = "S"
  }
}