provider "aws" {
  region = "eu-central-1"
}

variable "db_password" {
  description = "Password for RDS instance"
  type        = string
  sensitive   = true
}

resource "aws_db_instance" "weatherapp_db" {
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "15.7"
  instance_class       = "db.t3.micro"  # Free tier eligible
  db_name                 = "weatherapp"
  username             = "adminesraa"
  password             = var.db_password
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true
  publicly_accessible  = true  # So your local K8s can access it
}

output "rds_endpoint" {
  value = aws_db_instance.weatherapp_db.endpoint
}
