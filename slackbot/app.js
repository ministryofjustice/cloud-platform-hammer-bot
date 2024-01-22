import App from '@slack/bolt';
import fetch from 'node-fetch';

const NANO_SECOND = 1000000000

const LOG_LEVEL = process.env.ENVIRONMENT === "production" ? "info" : "debug"

const CHANNEL_ID = process.env.ENVIRONMENT === "production" ? "C57UPMZLY" : "C06EH6AR7DH"

const API_URL = process.env.ENVIRONMENT === "production" ? "http://api.cloud-platform-hammer-bot.svc.cluster.local:3001" : "http://localhost:3000"

const app = new App.App({
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  socketMode: true,
  appToken: process.env.SLACK_APP_TOKEN,
  port: process.env.PORT || 3000,
  logLevel: LOG_LEVEL
});

const getStatus = async (ids) => {
  const response = await fetch(`${API_URL}/check-pr?id=${ids}`);

  return await response.json();
}

const pushEmptyCommit = async (branch) => {
  const response = await fetch(`${API_URL}/retrigger-checks?branch=${branch}`);

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
    channel: CHANNEL_ID,
    timestamp: ts
  })
}

const addEmoji = async (emoji, ts) => {
  return await app.client.reactions.add({
    name: emoji,
    channel: CHANNEL_ID,
    timestamp: ts
  })
}

const postFail = async (data, ts) => {
  console.log("debug fail data some error", data)
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
    channel: CHANNEL_ID,
    text: message,
    icon_emoji: "robot_face",
    thread_ts: ts
  })
}

const retryLater = (ids, ts, retryInMs) => {
  setTimeout(async () => {
    const data = await getStatus(ids)

    const result = await postReaction(data, ts)

    await removeEmoji("spinning-circle-of-death", ts)

    if (result) {
      return true
    }

    await postReply("It looks like checks on your pr are _still_ pending even after waiting a while. A Cloud Platform team member will come and take a look.", ts)
    await addEmoji("hourglass_flowing_sand", ts)
    await addEmoji("warning", ts)

  }, retryInMs)
}

const postPendingRecent = async (data, ts, ids) => {
  const pendingRecent = data.filter((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec > 0) : false)

  if (pendingRecent.length && pendingRecent.length > 0) {
    const retryIn = pendingRecent.map((pr) => pr.InvalidChecks.map((check) => check.RetryInNanoSec)).flat().sort((a, b) => a - b)[0]

    await addEmoji("spinning-circle-of-death", ts)

    retryLater(ids, ts, (retryIn / NANO_SECOND) * 1000 + 10)

    return true
  }

  return false
}

app.message('github.com/ministryofjustice/cloud-platform-environments/pull/', async ({ message }) => {
  console.log('msg', message)

  const pulls = message.text.match(/\/pull\/\d+/g);

  const pullIds = pulls.map((match) => match.split("/pull/")[1]);

  const ids = pullIds.join(",")

  const data = await getStatus(ids)

  console.log(JSON.stringify(data))

  const result = await postReaction(data, message.ts)

  if (result) {
    return
  }

  const pendingResult = await postPendingRecent(data, message.ts, ids)

  if (pendingResult) {
    return
  }

  const pendingOlderThan10Mins = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec === 0) : false)

  if (pendingOlderThan10Mins) {
    await addEmoji("hourglass_flowing_sand", message.ts)
    await postReply("Looks like concourse needs a kick, Hammer-Bot has pushed an empty commit to your pull request", message.ts)

    const branch = data.map((pr) => pr.Branch).join(",")
    console.log("branch", branch)

    await pushEmptyCommit(branch)
    await removeEmoji("hourglass_flowing_sand", message.ts)
    await addEmoji("spinning-circle-of-death", message.ts)

    retryLater(ids, message.ts, 480000) // 8 mins in ms
  }
});


(async () => {
  await app.start();

  console.log('⚡️ Bolt app is running!');
})();
