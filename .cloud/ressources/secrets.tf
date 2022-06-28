resource "aws_secretsmanager_secret" "service" {
  name = "service/webhooks"
}

resource "aws_secretsmanager_secret_version" "service" {
  secret_id = aws_secretsmanager_secret.service.id
  secret_string = jsonencode({
    MONGODB_CONN_STRING = "mongodb+srv://${var.app_env_name}:${random_password.password.result}@${local.mongodb_atlas_url}"
    MONGODB_PASSWORD = random_password.password.result
    MONGODB_USERNAME = var.app_env_name
    MONGODB_HOST = local.mongodb_atlas_url })
}