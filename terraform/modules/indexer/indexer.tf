resource "argocd_application" "indexer" {
  metadata {
    name      = "${var.name}-indexer"
    namespace = "argocd"
    labels = {
      app     = "${var.name}-indexer"
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
      path            = "k8s/indexer"
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
          name  = "env.rpc.http"
          value = var.rpc_endpoint
        }

        parameter {
          name  = "env.rpc.ws"
          value = var.rpc_ws_endpoint
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

