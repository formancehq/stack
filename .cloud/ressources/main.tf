data "aws_caller_identity" "current" {}

locals {
  mongodb_atlas_url = trimprefix(mongodbatlas_advanced_cluster.test.connection_strings.0.private_srv, "mongodb+srv://")
}

resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = var.atlas_project_id
  name                   = var.app_env_name
  cluster_type           = "REPLICASET"
  backup_enabled         = true
  mongo_db_major_version = "5.0"
  pit_enabled = true
  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      auto_scaling {
        disk_gb_enabled = true
        compute_enabled = true
        compute_scale_down_enabled = true
        compute_min_instance_size = "M10"
        compute_max_instance_size = "M40"
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}

resource "random_password" "password" {
  length           = 16
  special          = false
}

resource "mongodbatlas_database_user" "user" {
  username           = var.app_env_name
  password           = random_password.password.result
  project_id         = var.atlas_project_id
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = "webhooks"
  }

  scopes {
    type = "CLUSTER"
    name = mongodbatlas_advanced_cluster.test.name
  }
}