resource "helm_release" "postgres" {
  name       = "my-postgres"
  chart      = "bitnami/postgresql"
  version    = "10.x.x"
  namespace  = "default"

  values = [
    file("${path.module}/values.yaml")
  ]
}