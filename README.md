# cloud-platform-hammer-bot

TODO: 

NEXT:
- add error case for checkPrStatus()
- look into what the value is being returned by retry
- add diagrams

AFTER:
- we should also think about in the future pushing an empty commit to retrigger the check (separate api route)
- prep deploy by copying over the deploy files

SLACKBOT:
- final step call the api from the slackbot and then post relevant emojis and if a check is queued for a long time this is probably concourse stuck so we should fire off another call to the api to push an empty commit

FINAL:
- get api deployed
