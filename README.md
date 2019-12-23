discord-vc-notifier
===================

[![License](https://img.shields.io/github/license/nownabe/discord-vc-notifier.svg?style=popout)](https://github.com/nownabe/discord-vc-notifier/blob/master/LICENSE.txt)

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

Deploy by GitHub Actions:

* GCP preparations:
  * Create a service account who has following roles:
    * App Engine Deployer
    * App Engine Service Admin
    * Cloud Build Editor
    * Storage Admin
  * Create JSON key of the service account
  * Enable App Engine Admin API
* Fork this repository
* Set required secrets
  * `GOOGLE_CREDENTIALS`: JSON key
  * `GOOGLE_PROJECT_ID`: GCP project ID
  * `DISCORD_BOT_TOKEN`
* Set optional secrets
  * `DISCORD_CHANNEL_ID`: if notify discord channel
