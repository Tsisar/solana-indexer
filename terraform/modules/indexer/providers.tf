terraform {
  required_providers {
    onepassword = {
      source  = "1Password/onepassword"
      version = ">= 2.1.2"
    }

    argocd = {
      source  = "oboukili/argocd"
      version = ">= 6.1.1"
    }

    github = {
      source  = "integrations/github"
      version = ">= 6.3.0"
    }
  }
}

provider "github" {
  token = var.github_token
  owner = "Tsisar"
}