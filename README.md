# cloud-platform-hammer-bot

TODO: 

- CheckPrStatus function handle queued checks (these will pick up concourse plans which are stuck) we should also think about in the future pushing an empty commit to retrigger the check (separate api route)
- Add route to get checks
- test with actual prs by curling their pr number at the api
- get api deployed
- final step call the api from the slackbot and then post relevant emojis and if a check is queued for a long time this is probably concourse stuck so we should fire off another call to the api to push an empty commit
