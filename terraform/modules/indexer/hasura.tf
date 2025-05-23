resource "argocd_application" "hasura" {
  metadata {
    name      = "${var.name}-hasura"
    namespace = "argocd"
    labels = {
      app     = "${var.name}-hasura"
      type    = var.type
      project = var.project
      network = var.network
    }
  }

  cascade = true

  spec {
    project = "default"

    destination {
      name      = "in-cluster"
      namespace = kubernetes_namespace.indexer.metadata[0].name
    }

    source {
      repo_url        = "git@github.com:desync-labs/splyce-infrastructure.git"
      path            = "k8s/hasura"
      target_revision = var.branch

      helm {
        value_files = [local.values_yaml]

        parameter {
          name  = "env.postgres.db"
          value = var.postgres_db
        }

        parameter {
          name  = "env.postgres.host"
          value = "${argocd_application.postgres.metadata[0].name}-service"
        }

        parameter {
          name  = "env.ingress.host"
          value = var.host
        }
      }
    }

    sync_policy {
      automated {
        prune     = true
        self_heal = true
      }
    }
  }

  depends_on = [
    kubernetes_namespace.indexer
  ]
}

resource "random_password" "hasura_password" {
  length  = 12
  special = false
}

resource "kubernetes_secret" "hasura_credentials" {
  metadata {
    name      = "hasura-credentials"
    namespace = kubernetes_namespace.indexer.metadata[0].name

    labels = {
      owner                        = "desynclabs"
      "app.kubernetes.io/instance" = kubernetes_namespace.indexer.metadata[0].name
      "app.kubernetes.io/name"     = "${var.name}-hasura"
    }
  }

  type = "Opaque"

  data = {
    password = random_password.hasura_password.result
  }
}

resource "onepassword_item" "hasura_credentials" {
  vault = var.op_vault_uuid

  title    = "HASURA_CREDENTIALS"
  category = "login"
  username = "admin"
  password = random_password.hasura_password.result
  url      = "${var.name}-hasura-service"

  tags = var.tags

  note_value = "Created by Terraform"
}

output "hasura_user" {
  value = var.hasura_user
}

output "hasura_password" {
  value = random_password.hasura_password.result
}
