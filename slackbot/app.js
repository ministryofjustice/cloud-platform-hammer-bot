import App from '@slack/bolt';
import fetch from 'node-fetch';

// Initializes your app with your bot token and signing secret
const app = new App.App({
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  socketMode: true,
  appToken: process.env.SLACK_APP_TOKEN,
  port: process.env.PORT || 3000,
  logLevel: "debug"
});

const getStatus = async (ids) => {
  const response = await fetch(`https://hammer-bot.live.cloud-platform.service.justice.gov.uk/check-pr?id=${ids}`);

  return await response.json();
}

app.message('github.com/ministryofjustice/cloud-platform-environments/pull/', async ({ message }) => {
  console.log('msg', message)

  const pulls = message.text.match(/\d+/g);

  const ids = pulls.join(",")

  const data = await getStatus(ids)

  console.log(JSON.stringify(data))

  if (data === null || data.length === 0) {
    await app.client.reactions.add({
      name: "sparkle",
      channel: "C05EG79V8HW",
      timestamp: message.ts
    })
    return
  }

  const failed = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 1) : false)

  if (failed) {
    await app.client.reactions.add({
      name: "x",
      channel: "C05EG79V8HW",
      timestamp: message.ts
    })
    return
  }

  const pendingRecent = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec > 0) : false)

  if (pendingRecent) {
    await app.client.reactions.add({
      name: "repeat",
      channel: "C05EG79V8HW",
      timestamp: message.ts
    })
    return

    // TODO add retry here
    // setTimeout(async () => {
    //   await getStatus(pendingRecent.Id)
    // }, pendingRecent.RetryIn * 1000)
  }

  const pendingOlderThan10Mins = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec === 0) : false)

  if (pendingOlderThan10Mins) {
    await app.client.reactions.add({
      name: "hourglass_flowing_sand",
      channel: "C05EG79V8HW",
      timestamp: message.ts
    })
    return

    // TODO trigger empty commit and then check again in x mins
  }
});


(async () => {
  await app.start();

  console.log('⚡️ Bolt app is running!');
})();
