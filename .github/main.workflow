workflow "deploy" {
  resolves = ["cloud run deploy"]
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
  args = "builds submit --tag gcr.io/npmfs-242515/server"
  env = {
    CLOUDSDK_CORE_PROJECT = "npmfs-242515"
  }
}

action "cloud run deploy" {
  uses = "actions/gcloud/cli@master"
  needs = ["build container image"]
  args = "--quiet beta run deploy --image gcr.io/npmfs-242515/server --allow-unauthenticated --region=us-central1 --timeout=8s server"
  env = {
    CLOUDSDK_CORE_PROJECT = "npmfs-242515"
  }
}
