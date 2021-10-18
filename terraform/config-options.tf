variable "function_name"{
	type = string
}

variable "s3_bucket" {
	type = string
}

variable "s3_key" {
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
