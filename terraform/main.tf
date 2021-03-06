terraform {
  required_providers {
    heroku = {
      source  = "heroku/heroku"
      version = "5.0.2"
    }
  }
}

provider "heroku" {
  email   = var.email
  api_key = var.api-key
}

resource "heroku_app" "app" {
  name   = var.app-name
  region = "us"
  # (Optional) not necessary if you are only using assets precompile.used only assets:precompile
  buildpacks = ["heroku/go"]
  # (Optional) enviroment variables.
  config_vars = {
    ACCESS_SECRET  = var.access_secret
    REFRESH_SECRET = var.refresh_secret
  }
}

resource "heroku_addon" "database" {
  app_id = var.app-id
  plan   = "heroku-postgresql:hobby-dev"
}

resource "heroku_addon" "redis" {
  app_id = var.app-id
  plan   = "heroku-redis:hobby-dev"
}
