# outputs.tf

output "lb_target_group_arn" {
  description = "Arn of the target group"
  value       = aws_lb_target_group.this.arn
}

output "default_subnet_ids" {
  description = "Ids of the default subnets"
  value       = [for subnet in aws_default_subnet.this : subnet.id]
}

output "lb_dns_name" {
  description = "DNS of load balancer applications"
  value       = aws_alb.this.dns_name
}
