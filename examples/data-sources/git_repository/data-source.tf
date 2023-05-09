data "git_repository" "example" {
  path = "./some-git-repository"
}

output "example" {
  value = data.git_repository.example
}

terraform {
  required_providers {
    git = {
      source  = "ekristen/git"
      version = ">= 0.1.0"
    }
  }
}