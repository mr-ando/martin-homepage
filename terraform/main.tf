
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  required_version = ">= 1.2"

  backend "s3" {
    bucket  = "martin-terraform-58994"
    key     = "terraform/state/terraform.tfstate"
    region  = "eu-north-1"
    encrypt = true
  }
}

provider "aws" {
  region = "eu-north-1"
}

# Use S3 Backend.
resource "aws_amplify_app" "homepage" {
  name         = "martin-homepage"
  repository   = "https://github.com/mr-ando/martin-homepage"
  access_token = var.github_access_token

  enable_auto_branch_creation = true
  enable_branch_auto_build    = true

  build_spec = <<-EOT
version: 1
frontend:
  phases:
    preBuild:
      commands:
        - npm ci --cache .npm --prefer-offline
    build:
      commands:
        - npm run build
  artifacts:
    baseDirectory: dist
    files:
      - '**/*'
  cache:
    paths:
      - .npm/**/*
      - node_modules/**/*
  EOT

  #custom_rule {
  #  source = "/<*>"
  #  status = "404"
  #  target = "/404.html" # Make sure you have a 404.html page in Astro
  #}


  # Static assets - let them pass through normally
  custom_rule {
    source = "/_astro/<*>"
    status = "200"
    target = "/_astro/<*>"
  }

  # Handle other static assets
  custom_rule {
    source = "/assets/<*>"
    status = "200"
    target = "/assets/<*>"
  }

  # Handle favicon and other root files
  custom_rule {
    source = "/*.ico"
    status = "200"
    target = "/<*>"
  }

  custom_rule {
    source = "/*.png"
    status = "200"
    target = "/<*>"
  }

  custom_rule {
    source = "/*.jpg"
    status = "200"
    target = "/<*>"
  }

  custom_rule {
    source = "/*.svg"
    status = "200"
    target = "/<*>"
  }


  # Environment variables
  environment_variables = {
    ENV = "prod"
  }

  # Platform can be WEB or WEB_COMPUTE
  platform = "WEB"

  tags = {
    Environment = "production"
  }
}

# Create a branch
resource "aws_amplify_branch" "main" {
  app_id      = aws_amplify_app.homepage.id
  branch_name = "main"

  enable_auto_build = true


  stage = "PRODUCTION"

  environment_variables = {
    PUBLIC_BASE_API_URL = "https://api.example.com"
    PUBLIC_BASE_WS_URL  = "example.com"
  }
}

# Auto-deploy from repository
resource "aws_amplify_webhook" "main" {
  app_id      = aws_amplify_app.homepage.id
  branch_name = aws_amplify_branch.main.branch_name
  description = "triggerMain"
}
# Create Amplify resource or S3 for frontend

# Create ec2 instance for golang backend.
