output "endpoint" {
  value = "https://${aws_route53_record.cname.fqdn}"
}

output "alb_dns_name" {
  value = "${aws_alb.this.dns_name}"
}
