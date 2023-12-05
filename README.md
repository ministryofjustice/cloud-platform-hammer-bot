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

