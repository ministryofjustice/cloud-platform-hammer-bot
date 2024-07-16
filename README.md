# cloud-platform-hammer-bot

## Purpose

> Increase visibility of user pull request checks

## Why?

> Users were posting prs before they had finished their checks meaning cloud platform team would waste time checking a pr only to see it had failed or they would have to wait to see the result. The bot also automates some common PR check scenarios.

## How?

The flow of the hammer bot is to check slack messages containing a pull request pr and then check the statuses and apply an emoji to indicate to the reviewer the pr is ready for review.

![architecture diagram](./images/api_diagram.png)

There are 2 components:

1. the [slackbot](slackbot/) which watches the slack channel for relevant messages. When a result is received from the api the slackbot will then add an emoji to the message in the slack channel.
2. the hammer-bot api which processes the pr and looks up the status of the pr checks and returns the result back to the slackbot

![Go diagram](./images/go_diagram.png)

TODO:

AFTER:
- we should also think about in the future pushing an empty commit to retrigger the check (separate api route)
- remove api ingress and call the api from the slack bot via the api service (not sure if we definetely need to do this)

## Local Testing
app.js is the main entry point for the slackbot. It can be run locally with the following commands:
the tokens can be found in the slack app settings page under 'OAuth & Permissions' and 'Basic Information'

```bash
export SLACK_BOT_TOKEN=<slack bot token>
export SLACK_SIGNING_SECRET=<slack signing secret>
export SLACK_APP_TOKEN=<slack app token>

node app.js
``` 

main.go is the main entry point for the api. It can be run locally with the following commands:

```bash
docker build . -t hammerbot:latest

docker run -p 3000:3000 --hostname=cf432824e79f --user=1000 --env=GITHUB_TOKEN=<create a token in github to use here> --env=GITHUB_URL=https://github.com/ministryofjustice/cloud-platform-environments --env=GITHUB_USER=<set your user for testing> --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --workdir=/ --runtime=runc -d hammerbot:latest
```