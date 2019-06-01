workflow "ci/cd" {
  resolves = ["cloud run deploy", "go test"]
  on = "push"
}

action "go test" {
  uses = "cedrickring/golang-action@1.3.0"
}

action "gcp authenticate" {
  needs = ["go test"]
  uses = "actions/gcloud/auth@master"
  secrets = ["GCLOUD_AUTH"]
}

action "build container image" {
  needs = ["gcp authenticate"]
  uses = "actions/gcloud/cli@master"
  args = "builds submit --tag gcr.io/rejstry/server"
}

action "cloud run deploy" {
  uses = "actions/gcloud/cli@master"
  needs = ["build container image"]
  args = "beta run deploy --image gcr.io/rejstry/server --allow-unauthenticated --region=us-central1 --timeout=8s"
}
