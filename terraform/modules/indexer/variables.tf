variable "op_vault_uuid" {
  description = "1Password Vault UUID"
  type        = string
}

variable "name" {
  description = "The name of the project"
  type        = string
}

variable "namespace" {
  description = "Namespace"
  type        = string
}

variable "type" {
  description = "Type"
  type        = string
  default     = "indexer"
}

variable "project" {
  description = "Project"
  type        = string
}

variable "network" {
  description = "Network"
  type        = string
}

variable "branch" {
  description = "Branch to deploy"
  type        = string
}

variable "tags" {
  description = "labels"
  type = list(string)
}

variable "host" {
  description = "The external host"
  type        = string
}

variable "hasura_user" {
  description = "Hasura user"
  type        = string
  default     = "admin"
}

variable "rpc_ws_endpoint" {
  description = "RPC WebSocket endpoint"
  type        = string
}

variable "rpc_endpoint" {
  description = "RPC endpoint"
  type        = string
}
