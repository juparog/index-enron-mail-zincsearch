# main.tf

### Locals ###
locals {
  subnet_suffixes = toset(["a", "b", "c"])
}
##############

### Vpc and subnets ###
resource "aws_default_vpc" "default_vpc" {
}


resource "aws_default_subnet" "this" {
  for_each = local.subnet_suffixes

  availability_zone = "${var.region_name}${each.value}"
}
########################

### Security groups ###
resource "aws_security_group" "this" {
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
########################

### Lb app ###
resource "aws_alb" "this" {

  name               = var.lb_name
  load_balancer_type = "application"
  subnets            = [for subnet in aws_default_subnet.this : subnet.id]
  security_groups    = ["${aws_security_group.this.id}"]
}

resource "aws_lb_target_group" "this" {
  name        = var.target_group_name
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_default_vpc.default_vpc.id
  health_check {
    matcher = "200,301,302"
    path    = "/"
  }
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_alb.this.arn
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
###############
