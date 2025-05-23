locals {
  values_yaml = var.branch == "master" ? "values.yaml" : "values-dev.yaml"
}

# Create namespace
resource "kubernetes_namespace" "indexer" {
  metadata {
    name = var.namespace
  }
}

# Setup docker credentials to pull images
module "docker_credentials" {
  source        = "../docker-credentials"
  namespace     = kubernetes_namespace.indexer.metadata[0].name
  op_vault_uuid = var.op_vault_uuid
}