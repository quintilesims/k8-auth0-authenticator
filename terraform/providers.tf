provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

provider "auth0" {
  domain        = "${var.auth0_domain}"
  client_id     = "${var.auth0_client_id}"
  client_secret = "${var.auth0_client_secret}"
}
