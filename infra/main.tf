terraform {
  backend "s3" {
    bucket       = "henrique-todo-tfstate-bucket"
    key          = "todo-api/terraform.tfstate"
    region       = "eu-north-1"
    use_lockfile = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

variable "region" {
  default = "eu-north-1"
}

provider "aws" {
  region = var.region
}

locals {
  binary_path = "${path.module}/../bin/todo-api"
  binary_key  = "todo-api"
  init_sql    = "${path.module}/../db/init.sql"
}

# ---------- AMI ----------
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

# ---------- Default VPC ----------
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# ---------- S3: holds the binary and init.sql ----------
resource "aws_s3_bucket" "artifacts" {
  bucket_prefix = "todo-api-artifacts-"
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "artifacts" {
  bucket                  = aws_s3_bucket.artifacts.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_object" "binary" {
  bucket = aws_s3_bucket.artifacts.id
  key    = local.binary_key
  source = local.binary_path
  etag   = filemd5(local.binary_path)
}

resource "aws_s3_object" "init_sql" {
  bucket = aws_s3_bucket.artifacts.id
  key    = "init.sql"
  source = local.init_sql
  etag   = filemd5(local.init_sql)
}

# ---------- Database password: generated, stored in SSM ----------
resource "random_password" "db" {
  length  = 24
  special = false
}

resource "aws_ssm_parameter" "db_password" {
  name  = "/todo-api/db-password"
  type  = "SecureString"
  value = random_password.db.result
}

# ---------- IAM: lets the instance read S3 and SSM ----------
resource "aws_iam_role" "instance" {
  name = "todo-api-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Action    = "sts:AssumeRole"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy" "read_artifacts" {
  name = "read-artifacts"
  role = aws_iam_role.instance.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = "${aws_s3_bucket.artifacts.arn}/*"
      },
      {
        Effect   = "Allow"
        Action   = ["ssm:GetParameter"]
        Resource = aws_ssm_parameter.db_password.arn
      }
    ]
  })
}

# Lets you open a shell via Session Manager (no SSH, no key pair, no port 22)
resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.instance.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "instance" {
  name = "todo-api-instance-profile"
  role = aws_iam_role.instance.name
}

# ---------- Security groups ----------
resource "aws_security_group" "api" {
  name = "todo-api-sg"

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Only the API's security group may talk to the database
resource "aws_security_group" "db" {
  name   = "todo-db-sg"
  vpc_id = data.aws_vpc.default.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.api.id]
  }
}

# ---------- RDS (Postgres) ----------
resource "aws_db_subnet_group" "db" {
  name       = "todo-db-subnets"
  subnet_ids = data.aws_subnets.default.ids
}

resource "aws_db_instance" "db" {
  identifier             = "todo-db"
  engine                 = "postgres"
  engine_version         = "16"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  db_name                = "todo"
  username               = "todoadmin"
  password               = random_password.db.result
  db_subnet_group_name   = aws_db_subnet_group.db.name
  vpc_security_group_ids = [aws_security_group.db.id]
  publicly_accessible    = false
  skip_final_snapshot    = true # fine for learning, not for production
}

# ---------- EC2 ----------
resource "aws_instance" "example" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = "t3.micro"
  iam_instance_profile   = aws_iam_instance_profile.instance.name
  vpc_security_group_ids = [aws_security_group.api.id]

  # user_data only runs on first boot, so replace the instance when it changes
  user_data_replace_on_change = true

  user_data = <<-EOF
    #!/bin/bash
    set -e
    # change detection: ${filemd5(local.binary_path)} ${filemd5(local.init_sql)}

    dnf install -y postgresql16

    aws s3 cp s3://${aws_s3_bucket.artifacts.id}/${local.binary_key} /usr/local/bin/todo-api --region ${var.region}
    aws s3 cp s3://${aws_s3_bucket.artifacts.id}/init.sql /tmp/init.sql --region ${var.region}
    chmod +x /usr/local/bin/todo-api

    DB_PASSWORD=$(aws ssm get-parameter --name ${aws_ssm_parameter.db_password.name} --with-decryption --query Parameter.Value --output text --region ${var.region})

    # config.yaml read by the Go app (root-only), matches the Config struct
    mkdir -p /etc/todo-api
    cat > /etc/todo-api/config.yaml <<YAML
    database:
      host: "${aws_db_instance.db.address}"
      port: 5432
      user: "${aws_db_instance.db.username}"
      password: "$DB_PASSWORD"
      name: "${aws_db_instance.db.db_name}"
      sslmode: "require"
    YAML
    chmod 600 /etc/todo-api/config.yaml

    # Load the schema (init.sql must be idempotent: CREATE TABLE IF NOT EXISTS ...)
    PGPASSWORD=$DB_PASSWORD psql -h ${aws_db_instance.db.address} -U ${aws_db_instance.db.username} -d ${aws_db_instance.db.db_name} -f /tmp/init.sql

    cat > /etc/systemd/system/todo-api.service <<'UNIT'
    [Unit]
    Description=ToDo API
    After=network-online.target
    StartLimitIntervalSec=0

    [Service]
    RestartSec=5
    WorkingDirectory=/etc/todo-api
    ExecStart=/usr/local/bin/todo-api --config /etc/todo-api/config.yaml serve --no-docker --addr :8080
    Restart=always

    [Install]
    WantedBy=multi-user.target
    UNIT

    systemctl daemon-reload
    systemctl enable todo-api
    systemctl start todo-api || true

    # Print diagnostics to the system log (visible without a shell)
    sleep 5
    echo "=== todo-api diagnostics ==="
    uname -m
    ls -l /usr/local/bin/todo-api /etc/todo-api/ || true
    file /usr/local/bin/todo-api || true
    systemctl status todo-api --no-pager || true
    journalctl -u todo-api -n 50 --no-pager || true
    echo "=== end diagnostics ==="
  EOF

  depends_on = [aws_s3_object.binary, aws_s3_object.init_sql]

  tags = {
    Name = "todo-api"
  }
}

# ---------- Outputs ----------
output "public_ip" {
  value = aws_instance.example.public_ip
}

output "db_endpoint" {
  value = aws_db_instance.db.address
}

