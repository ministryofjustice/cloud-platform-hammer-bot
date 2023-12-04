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

app.message('github.com/ministryofjustice/cloud-platform-environments/pull/', async ({ message, say }) => {
  // say() sends a message to the channel where the event was triggered
  console.log('msg', message)

  const pulls = message.text.match(/\d+/g);

  const response = await fetch(`https://hammer-bot.live.cloud-platform.service.justice.gov.uk/check-pr?id=${pulls.join(",")}`);

  const data = await response.json();

  console.log(data);
  // await say(`Hey there <@${message.user}>, ${message.text.split("pull/")[1]}`);

});


(async () => {
  await app.start();

  console.log('⚡️ Bolt app is running!');
})();
