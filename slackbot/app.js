import App from '@slack/bolt';
import fetch from 'node-fetch';

const NANO_SECOND = 1000000000

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

const postSuccess = async (data) => {
  if (data === null || data.length === 0) {
    await addEmoji("sparkles")
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

const postFail = async (data) => {
  const failed = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 1) : false)

  if (failed) {
    await addEmoji("x")
    return true
  }
  return false
}

const postReaction = async (data) => {
  const isSuccess = await postSuccess(data)

  if (isSuccess) {
    return true
  }

  const isFailed = await postFail(data)

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

const postPendingRecnt = async (data) => {
  const pendingRecent = data.some((pr) => pr.InvalidChecks.length ? pr.InvalidChecks.some((check) => check.Status === 2 && check.RetryInNanoSec > 0) : false)

  if (pendingRecent) {
    await addEmoji("repeat", message.ts)

    setTimeout(async () => {
      const data = await getStatus(pendingRecent.Id)

      const result = postReaction(data)

      await removeEmoji("repeat", message.ts)

      if (result) {
        return true
      }

      await postReply("It looks like checks on your pr are _still_ pending even after waiting a while. A Cloud Platform team member will come and take a look.", message.ts)
      await addEmoji("hourglass_flowing_sand", message.ts)
      await addEmoji("warning", message.ts)

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

  const result = await postReaction(data)

  if (result) {
    return
  }

  const pendingResult = await postPendingRecnt(data)

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
