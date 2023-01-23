# outputs.tf

output "lb_dns_name" {
  description = "DNS of load balancer applications"
  value       = module.network.lb_dns_name
}
