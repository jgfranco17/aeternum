variable "render_access_token" {
    type      = string
    sensitive = true
}

variable "deployment_environment" {
    type    = string
    default = "dev"
}
