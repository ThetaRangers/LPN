variable "function_name"{
	type = string
}

variable "lambda_zip" {
	type = string
}

variable "handler"{
	type = string
}

variable "dynamo_table_name"{
	type = string
}

variable "dynamo_table_name1"{
	type = string
}

variable "zone" {
	type = string
	default = "us-east-1"
}
