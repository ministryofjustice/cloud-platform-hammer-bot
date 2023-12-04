import App from '@slack/bolt';
import fetch from 'node-fetch';

const NANO_SECOND = 1000000000

const API_URL = process.env.ENVIRONMENT === "production" ? "http://api.cloud-platform-hammer-bot.svc.cluster.local:3001" : "https://hammer-bot.live.cloud-platform.service.justice.gov.uk"

// Initializes your app with your bot token and signing secret
const app = new App.App({
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  socketMode: true,
  appToken: process.env.SLACK_APP_TOKEN,
  port: process.env.PORT || 3000,
  logLevel: "info"
});

const getStatus = async (ids) => {
  const response = await fetch(`${API_URL}/check-pr?id=${ids}`);

  return await response.json();
}

const postSuccess = async (data, ts) => {
  if (data === null || data.length === 0) {
    await addEmoji("sparkles", ts)
    return true
  }
  return false
}

const removeEmoji = async (emoji, ts) => {
  return await app.client.reactions.remove({
    name: emoji,
    timestamp: ts
  })
}

const addEmoji = async (emoji, ts) => {
  return await app.client.reactions.add({
    name: emoji,
    channel: "C05EG79V8HW",
    timestamp: ts
  })
}

const postFail = async (data, ts) => {
  const failed = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 1) : false)

  if (failed) {
    await addEmoji("x", ts)
    return true
  }
  return false
}

const postReaction = async (data, ts) => {
  const isSuccess = await postSuccess(data, ts)

  if (isSuccess) {
    return true
  }

  const isFailed = await postFail(data, ts)

  if (isFailed) {
    return true
  }

  return false
}

const postReply = async (message, ts) => {
  return await app.client.chat.postMessage({
    channel: "C05EG79V8HW",
    text: message,
    icon_emoji: "robot_face",
    thread_ts: ts
  })
}

const postPendingRecnt = async (data, ts) => {
  const pendingRecent = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec > 0) : false)

  if (pendingRecent) {
    await addEmoji("repeat", ts)

    setTimeout(async () => {
      const data = await getStatus(pendingRecent.Id)

      const result = postReaction(data, ts)

      await removeEmoji("repeat", ts)

      if (result) {
        return true
      }

      await postReply("It looks like checks on your pr are _still_ pending even after waiting a while. A Cloud Platform team member will come and take a look.", ts)
      await addEmoji("hourglass_flowing_sand", ts)
      await addEmoji("warning", ts)

    }, pendingRecent.RetryIn / NANO_SECOND + 10)

    return true
  }

  return false
}


app.message('github.com/ministryofjustice/cloud-platform-environments/pull/', async ({ message }) => {
  console.log('msg', message)

  const pulls = message.text.match(/\d+/g);

  const ids = pulls.join(",")

  const data = await getStatus(ids)

  console.log(JSON.stringify(data))

  const result = await postReaction(data, message.ts)

  if (result) {
    return
  }

  const pendingResult = await postPendingRecnt(data, message.ts)

  if (pendingResult) {
    return
  }

  const pendingOlderThan10Mins = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec === 0) : false)

  if (pendingOlderThan10Mins) {
    await addEmoji("hourglass_flowing_sand", message.ts)
    await postReply("Looks like concourse needs a kick, plesase push an empty commit to retrigger the checks `git commit --allow-empty -m 'Empty - Commit'`", message.ts)
    // TODO trigger empty commit and then check again in x mins
  }
});


(async () => {
  await app.start();

  console.log('⚡️ Bolt app is running!');
})();
