# Terraform configuration for the API

terraform {
    required_providers {
        render = {
            source = "render-oss/render"
            version = "1.5.0"
        }
    }
}

provider "render" {
    api_key = var.render_access_token
}

resource "render_web_service" "web" {
    name               = "aeternum"
    plan               = "free"
    region             = "singapore"
    num_instances      = 3
    health_check_path  = "/healthz"

    runtime_source = {
        docker = {
            auto_deploy     = false
            branch          = "main"
            dockerfile_path = "./Dockerfile"
            repo_url        = "https://github.com/jgfranco17/aeternum-api"
        }
    }
}
