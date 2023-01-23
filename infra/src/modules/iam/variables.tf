# variables.tf

##### Required variables #####
variable "role_name" {
  description = "Role name for create iam role"
  type        = string
  validation {
    condition     = var.role_name != null && var.role_name != ""
    error_message = "Role name not null or empty"
  }
}

variable "arn_policies" {
  description = "Arn list to attach the policies"
  type        = list(string)
  validation {
    condition     = length(var.arn_policies) > 0
    error_message = "Arn list can not empty"
  }
}

variable "identifiers" {
  description = "Identifier service list for assume rol"
  type        = list(string)
  validation {
    condition     = length(var.identifiers) > 0
    error_message = "Identifiers list can not empty"
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
