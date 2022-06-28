data "aws_caller_identity" "this" {}

data "aws_eks_cluster" "this" {
  name = "cluster-${var.env}"
}

resource "aws_iam_role" "eks_role" {
  name = "app-webhooks"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sts:AssumeRoleWithWebIdentity"
        ]
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.this.account_id}:oidc-provider/${trimprefix(data.aws_eks_cluster.this.identity[0].oidc[0].issuer, "https://")}"
        }
      },
    ]
  })
  managed_policy_arns = [aws_iam_policy.secret_manager_policy.arn, aws_iam_policy.s3_policy.arn]
}

resource "aws_iam_role_policy_attachment" "eks_secret_manager_policy" {
  role       = aws_iam_role.eks_role.name
  policy_arn = aws_iam_policy.secret_manager_policy.arn
}

resource "aws_iam_policy" "secret_manager_policy" {
  name        = "SecretManagerPolicyForwebhooks"
  path        = "/"
  description = "Policy for EKS access to Secret Manager"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_policy" "s3_policy" {
  name        = "S3AccessForwebhooks"
  path        = "/"
  description = "Policy for webhooks access to S3"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:*Object",
          "s3:*Object*",
          "s3:List*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}