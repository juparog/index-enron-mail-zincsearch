# variables.tf

## Optionals ###
variable "settings_roles" {
  description = "value"
  default     = "test"
}

variable "ZINC_DATA_PATH" {
  description = "Path for ZincSearch data"
  type        = string
  default     = "/data"
}

variable "ZINC_FIRST_ADMIN_USER" {
  description = "Name for the admin user"
  type        = string
  default     = "admin"
}

variable "ZINC_FIRST_ADMIN_PASSWORD" {
  description = "Password for the admin user"
  type        = string
  default     = "Complexpass123"
}
################
