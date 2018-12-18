workflow "New workflow" {
  on = "push"
  resolves = ["GitHub Action for Docker"]
}

action "GitHub Action for Docker" {
  uses = "actions/docker/cli@0c53e4a"
  args = "build -t csweichel/ruruku:latest ."
}
