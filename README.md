# cloud-platform-hammer-bot

TODO: 

AFTER:
- we should also think about in the future pushing an empty commit to retrigger the check (separate api route)
- remove api ingress and call the api from the slack bot via the api service


1. if the checks are still pending then retry based on the retryIn field from the api

![api diagram](./images/api_diagram.png)
![Go diagram](./images/go_diagram.png)

