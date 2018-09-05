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

variable "docker_image" {
  type        = "string"
  description = "The Docker image to use"
  default     = "quintilesims/k8-auth0-authenticator:latest"
}

variable "auth0_domain" {
  type        = "string"
  description = "The domain of the Auth0 account to use"
  default     = "iqvia.auth0.com"
}

variable "auth0_client_id" {
  type        = "string"
  description = "The Auth0 Global Client ID"
}

variable "auth0_client_secret" {
  type        = "string"
  description = "The Auth0 Global Client Secret"
}

variable "scale" {
  type        = "string"
  description = "The scale of the service"
  default     = "1"
}

variable "acm_cert_domain" {
  type        = "string"
  description = "The domain of the ACM certificate to use"
  default     = "*.ims.io"
}

variable "domain" {
  type        = "string"
  description = "The top level domain of the route53 hosted zone to use"
  default     = "ims.io"
}

variable "subdomain" {
  type        = "string"
  description = "The desired subdomain for the service"
  default     = "k8-auth0-authenticator"
}
