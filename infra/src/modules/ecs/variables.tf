# variables.tf

##### Required variables #####
variable "cluster_name" {
  description = "ECS Cluster name"
  type        = string
  validation {
    condition     = var.cluster_name != null && var.cluster_name != ""
    error_message = "ESC Cluster name not null or empty"
  }
}

variable "task_name" {
  description = "Ecs task family name"
  type        = string
  validation {
    condition     = var.task_name != null && var.task_name != ""
    error_message = "ECS task name not null or empty"
  }
}

variable "task_img_name" {
  description = "ECS task imgage name"
  type        = string
  validation {
    condition     = var.task_img_name != null && var.task_img_name != ""
    error_message = "ECS task image name not null or empty"
  }
}

variable "target_group_arn" {
  description = "Target group arn"
  type        = string
  validation {
    condition     = var.target_group_arn != null && var.target_group_arn != ""
    error_message = "Target group arn not null or empty"
  }
}

variable "role_task_execution_arn" {
  description = "Arn role for ecs task execution"
  type        = string
  validation {
    condition     = var.role_task_execution_arn != null && var.role_task_execution_arn != ""
    error_message = "Rol arn for task execution not null or empty"
  }
}

variable "service_name" {
  description = "ECS service name in cluster"
  type        = string
  validation {
    condition     = var.service_name != null && var.service_name != ""
    error_message = "ECS service name not null or empty"
  }
}

variable "subnet_ids" {
  description = "Subnets ids for ecs service name"
  type        = list(any)
  validation {
    condition     = length(var.subnet_ids) > 0
    error_message = "Subnets list ids not empty"
  }
}
##############################

##### Optional variables #####

variable "env_task" {
  description = "Environment variables for task definition"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

##############################
