variable "host" {
  type = string
}

variable "client_certificate" {
  type = string
}

variable "client_key" {
  type = string
}

variable "cluster_ca_certificate" {
  type = string
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "postgres"
}

variable "db_pass" {
  description = "Database password"
  type        = string
}

variable "db_host" {
  description = "Database host name"
  type        = string
  default     = "my-postgres-postgresql"
}

variable "db_port" {
  description = "Database port number"
  type        = string
  default     = "5432"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "postgres"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}