resource "kubernetes_deployment" "todo-api-deployment" {
  metadata {
    name      = "todo-api"
    labels = {
      App = "todo-api"
    }
  }

  timeouts {
    create="5m"
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        App = "todo-api"
      }
    }

    template {
      metadata {
        labels = {
          App = "todo-api"
        }
      }

      spec {
        container {
          name  = "todo-api"
          image = "sidhathi/appcd-todo:latest"
          image_pull_policy = "IfNotPresent"

          port {
            container_port = 8000
            host_port = 8000
          }

          env {
            name  = "DB_HOST"
            value = var.db_host
          }

          env {
            name  = "DB_PORT"
            value = var.db_port
          }

          env {
            name  = "DB_USER"
            value = var.db_username
          }

          env {
            name  = "DB_NAME"
            value = var.db_name
          }

          env {
            name  = "DB_PASSWORD"
            value = var.db_pass
          }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
  depends_on = [helm_release.my-postgres]
}

resource "kubernetes_service" "todo-api-service" {
  metadata {
    name = "todo-api"
  }

  timeouts {
    create = "1m"
  }

  spec {
    selector = {
      App = "todo-api"
    }

    port {
      port        = 8000
      target_port = 8000
    }

    type = "ClusterIP"
  }
}

output "lb_ip" {
  value = kubernetes_service.todo-api-service.status.0.load_balancer.0.ingress.0.hostname
}
