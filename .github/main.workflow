workflow "deploy" {
  resolves = ["cloud run deploy"]
  on = "push"
}

action "go test" {
  uses = "cedrickring/golang-action@1.3.0"
}

action "gcp authenticate" {
  needs = ["go test"]
  uses = "g-harel/gcloud/auth@253"
  secrets = ["GCLOUD_AUTH"]
}

action "build container image" {
  needs = ["gcp authenticate"]
  uses = "g-harel/gcloud/cli@253"
  args = "builds submit --tag gcr.io/npmfs-242515/website"
  env = {
    CLOUDSDK_CORE_PROJECT = "npmfs-242515"
  }
}

action "cloud run deploy" {
  uses = "g-harel/gcloud/cli@253"
  needs = ["build container image"]
  args = "--quiet beta run deploy --image gcr.io/npmfs-242515/website --allow-unauthenticated --region=us-central1 --timeout=8s website"
  env = {
    CLOUDSDK_CORE_PROJECT = "npmfs-242515"
  }
}
