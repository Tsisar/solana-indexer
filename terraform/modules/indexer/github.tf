variable "repository" {
  description = "Repository of application source code"
  type        = string
}

variable "environment" {
  description = "GitHub Environment"
  type        = string
}

variable "argo_cd_host" {
  description = "ArgoCD Host"
  type        = string
}

variable "repository_secrets" {
  description = "GitHub Repository secrets"
  type = map(string)
  default = {}
}

variable "repository_variables" {
  description = "GitHub Repository variables"
  type = map(string)
  default = {}
}

variable "github_token" {
  description = "GitHub Token"
  type        = string
}

# Docker Hub Read & Write token
data "onepassword_item" "docker_credentials" {
  vault = var.op_vault_uuid
  title = "DOCKER_RW_CREDENTIALS"
}

# Secret for github-actions
locals {
  repository_secrets = merge(
    var.repository_secrets,
    {
      "ARGOCD_API_URL"    = var.argo_cd_host
      "ARGOCD_AUTH_TOKEN" = argocd_account_token.github_action.jwt
      "DOCKER_USERNAME"   = data.onepassword_item.docker_credentials.username
      "DOCKER_PASSWORD"   = data.onepassword_item.docker_credentials.password
    }
  )

  repository_variables = merge(
    var.repository_variables,
    {
      "ARGOCD_APP_NAME" = argocd_application.indexer.metadata[0].name
    }
  )
}

# Token for account github-action
resource "argocd_account_token" "github_action" {
  account = "github-actions"
}

resource "github_repository_environment" "dynamic" {
  repository  = var.repository
  environment = var.environment

  deployment_branch_policy {
    protected_branches     = false
    custom_branch_policies = true
  }
}

resource "github_actions_environment_secret" "dynamic" {
  for_each        = local.repository_secrets
  repository      = var.repository
  environment     = github_repository_environment.dynamic.environment
  secret_name     = each.key
  plaintext_value = each.value
}

resource "github_actions_environment_variable" "dynamic" {
  for_each      = local.repository_variables
  repository    = var.repository
  environment   = github_repository_environment.dynamic.environment
  variable_name = each.key
  value         = each.value
}