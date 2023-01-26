# variables.tf

### Required ###
variable "ZINC_FIRST_ADMIN_USER" {
  description = "Name for the admin user"
  type        = string
}

variable "ZINC_FIRST_ADMIN_PASSWORD" {
  description = "Password for the admin user"
  type        = string
}
################

## Optionals ###
variable "AWS_REGION" {
  description = "Aws region name"
  type        = string
  default     = "us-east-1"
}

variable "ZINC_DATA_PATH" {
  description = "Path for ZincSearch data"
  type        = string
  default     = "/data"
}

variable "ZINC_MAX_DOCUMENT_SIZE" {
  description = "Buffer read size"
  type        = number
  default     = 1453063655
}

variable "IDX_ENRONTGZ_FIELDS" {
  description = "Enron email fields per file"
  type        = string
  default     = "Message-ID:Date:From:To:Subject:Mime-Version:Content-Type:Content-Transfer-Encoding:X-From:X-To:X-cc:X-bcc:X-Folder:X-Origin:X-FileName"
}

variable "IDX_ENRONTGZ_FORMAT" {
  description = "Format upload for ZinSearch"
  type        = string
  default     = "bulkv2"
}

variable "IDX_ENRONTGZ_IDXNAME" {
  description = "Index name for save in ZinSearch"
  type        = string
  default     = ""
}

variable "IDX_ENRONTGZ_SEPARATOR" {
  description = "Separator of fields and data in the files"
  type        = string
  default     = ":"
}

variable "IDX_ENRONTGZ_TERMINATOR" {
  description = "Terminator of records in the files"
  type        = string
  default     = "\\n\\r"
}
################
