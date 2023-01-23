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
##############################

##### Optional variables #####
variable "region_name" {
  description = "Aws region name"
  type        = string
  default     = "us-east-1"
}
##############################
