variable "aws_access_key" {
  type        = "string"
  description = "The access key portion of an IAM key pair"
}

variable "aws_secret_key" {
  type        = "string"
  description = "The secret key portion of an IAM key pair"
}

variable "aws_region" {
  type        = "string"
  description = "The aws region to create resources"
  default     = "us-west-2"
}

variable "name" {
  type        = "string"
  description = "An identifier used to name, tag, and/or group resources"
  default     = "k8-auth0-authenticator"
}
