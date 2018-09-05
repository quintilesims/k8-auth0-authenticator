data "aws_availability_zones" "available" {}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "1.14.0"

  name = "${var.name}"
  cidr = "10.0.0.0/16"

  azs = [
    "${data.aws_availability_zones.available.names[0]}",
    "${data.aws_availability_zones.available.names[1]}",
    "${data.aws_availability_zones.available.names[2]}",
  ]

  private_subnets    = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets     = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
  enable_nat_gateway = true
  single_nat_gateway = true

  tags = {
    Name = "${var.name}"
  }
}

resource "aws_default_security_group" "default" {
  vpc_id = "${module.vpc.vpc_id}"

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.name}-default"
  }
}

resource "aws_security_group" "https" {
  name        = "${var.name}-https"
  vpc_id      = "${module.vpc.vpc_id}"
  description = "Allow communication on port 443"

  ingress {
    protocol    = "tcp"
    from_port   = "443"
    to_port     = "443"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.name}-https"
  }
}

data "aws_acm_certificate" "this" {
  domain   = "${var.acm_cert_domain}"
  statuses = ["ISSUED"]
}

data "aws_route53_zone" "top" {
  name = "${var.domain}"
}

resource "aws_route53_record" "cname" {
  zone_id = "${data.aws_route53_zone.top.zone_id}"
  name    = "${var.subdomain}.${var.domain}"
  type    = "CNAME"
  ttl     = "300"
  records = ["${aws_alb.this.dns_name}"]
}
