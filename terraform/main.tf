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
  # config_vars = {
  #   DATABASE_URL = var.db_url
  # }
}

resource "heroku_addon" "database" {
  app_id = var.app-id
  plan   = "heroku-postgresql:hobby-dev"
}
