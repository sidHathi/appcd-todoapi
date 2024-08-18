resource "helm_release" "my-postgres" {
  name       = "my-postgres"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "postgresql"
  version    = "15.5.21"

  timeout = 60

  values = [
    file("${path.module}/../values.yaml")
  ]
}
