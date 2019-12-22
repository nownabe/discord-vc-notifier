discord-vc-notifier
===================

# Prepare

* [Create new application](https://discordapp.com/developers/applications/)
* Add a bot
* Go to https://discordapp.com/oauth2/authorize?client_id=YOUR_CLIENT_ID&scope=bot&permission=0
  * You can find Client ID at OAuth2 page
* `cp .envrc.credentials.example .envrc.credentials`
* Add the bot token into .envrc.credentials

# Deploy

## Google App Engine

Preparation:

* Create GCP project
* Create an application of App Engine in the project

Deploy manually:

* Configure gcloud
  * `gcloud auth login`
  * `gcloud config set project YOUR_PROJECT_ID`
* Create `app.yaml`
  * envsubst < app.yaml.tmpl > app.yaml
* `gcloud app deploy`
