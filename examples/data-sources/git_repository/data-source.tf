data "git_repository" "example" {
  path = "./some-git-repository"
}

output "example" {
  value = {
    path    = data.git_repository.example.path
    branch  = data.git_repository.example.branch
    summary = data.git_repository.example.summary
    semver  = data.git_repository.example.semver
  }
}

terraform {
  required_providers {
    git = {
      source  = "ekristen/git"
      version = ">= 0.1.0"
    }
  }
}