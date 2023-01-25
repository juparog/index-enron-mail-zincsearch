# main.tf

### Cluster ###
resource "aws_ecs_cluster" "this" {
  name = var.cluster_name
}
###############

### ECS tasks ###
resource "aws_ecs_task_definition" "this" {
  family                   = var.task_name
  container_definitions    = <<DEFINITION
  [
    {
      "name": "${var.task_name}",
      "image": "${var.task_img_name}",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 4080,
          "hostPort": 4080
        }
      ],
      "environment": ${jsonencode(var.env_task)},
      "memory": 512,
      "cpu": 256
    }
  ]
  DEFINITION
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  memory                   = 512
  cpu                      = 256
  execution_role_arn       = var.role_task_execution_arn
}
#################

### ECS service ###
resource "aws_ecs_service" "this" {
  name            = var.service_name
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.this.arn
  launch_type     = "FARGATE"
  desired_count   = 1

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = aws_ecs_task_definition.this.family
    container_port   = 4080
  }

  network_configuration {
    subnets          = var.subnet_ids
    assign_public_ip = true
  }
}
###################
