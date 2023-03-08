data "git_repository" "example" {
  path = "./some-git-repository"
}

output "example" {
  value = {
    path      = data.git_repository.example.path
    branch    = data.git_repository.example.branch
    tag       = data.git_repository.example.tag
    is_dirty  = data.git_repository.example.is_dirty
    is_tag    = data.git_repository.example.is_tag
    is_branch = data.git_repository.example.is_branch
    summary   = data.git_repository.example.summary
    semver    = data.git_repository.example.semver
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