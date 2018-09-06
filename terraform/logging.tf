resource "aws_iam_role" "logger" {
  name               = "${var.name}-ecs-logger"
  assume_role_policy = "${file("ecs_assume_role_policy.json")}"
}

resource "aws_iam_role_policy" "logger" {
  name   = "${var.name}-logger"
  role   = "${aws_iam_role.logger.id}"
  policy = "${file("logger_role_policy.json")}"
}

resource "aws_cloudwatch_log_group" "this" {
  name = "${var.name}"

  tags {
    Name = "${var.name}"
  }
}
