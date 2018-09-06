resource "aws_ecs_cluster" "this" {
  name = "${var.name}"
}

resource "aws_alb" "this" {
  name            = "${var.name}"
  subnets         = ["${module.vpc.public_subnets}"]
  security_groups = ["${aws_default_security_group.default.id}", "${aws_security_group.https.id}"]

  tags = {
    Name = "${var.name}"
  }
}

resource "aws_alb_target_group" "this" {
  name        = "${aws_alb.this.name}"
  vpc_id      = "${module.vpc.vpc_id}"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"

  health_check = {
    path                = "/"
    protocol            = "HTTP"
    interval            = 5
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = 200
  }

  tags = {
    Name = "${var.name}"
  }
}

resource "aws_alb_listener" "this" {
  load_balancer_arn = "${aws_alb.this.id}"
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = "${data.aws_acm_certificate.this.arn}"

  default_action {
    target_group_arn = "${aws_alb_target_group.this.id}"
    type             = "forward"
  }
}

data "template_file" "container_definitions" {
  template = "${file("containers.json")}"

  vars {
    docker_image      = "${var.docker_image}"
    log_group_name    = "${aws_cloudwatch_log_group.this.name}"
    log_stream_prefix = "${var.name}"
    aws_region        = "${var.aws_region}"
    auth0_domain      = "${var.auth0_domain}"
    auth0_client_id   = "${auth0_client.this.client_id}"
  }
}

resource "aws_ecs_task_definition" "this" {
  family                   = "${var.name}"
  container_definitions    = "${data.template_file.container_definitions.rendered}"
  execution_role_arn       = "${aws_iam_role.logger.arn}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
}

resource "aws_ecs_service" "this" {
  name            = "${var.name}"
  cluster         = "${aws_ecs_cluster.this.id}"
  task_definition = "${aws_ecs_task_definition.this.arn}"
  launch_type     = "FARGATE"
  desired_count   = "${var.scale}"

  load_balancer {
    target_group_arn = "${aws_alb_target_group.this.id}"
    container_name   = "k8-auth0-authenticator"
    container_port   = "80"
  }

  network_configuration = {
    subnets         = ["${module.vpc.private_subnets}"]
    security_groups = ["${aws_default_security_group.default.id}", "${aws_security_group.https.id}"]
  }
}
