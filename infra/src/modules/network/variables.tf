# variables.tf

##### Required variables #####
variable "lb_name" {
  description = "Load balance name"
  type        = string
  validation {
    condition     = var.lb_name != null && var.lb_name != ""
    error_message = "Load balance name not null or empty"
  }
}

variable "target_group_name" {
  description = "Target group name"
  type        = string
  validation {
    condition     = var.target_group_name != null && var.target_group_name != ""
    error_message = "Target group name not null or empty"
  }
}
##############################

##### Optional variables #####
variable "region_name" {
  description = "Aws region name"
  type        = string
  default     = "us-east-1"
}
##############################
