# outputs.tf

output "iam_execution_role_arn" {
  description = "Arn of the role for execution"
  value       = aws_iam_role.this.arn
}
