const { App } = require('@slack/bolt');

// Initializes your app with your bot token and signing secret
const app = new App({
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  socketMode: true, // add this
  appToken: process.env.SLACK_APP_TOKEN, // add this
  port: process.env.PORT || 3000
});

// Listens to incoming messages that contain "hello"
app.message('https://github.com/ministryofjustice/cloud-platform-environments/pull/', async ({ message, say }) => {
  // say() sends a message to the channel where the event was triggered
  console.log('msg', message)
  await say(`Hey there <@${message.user}>!`);
});

(async () => {
  // Start your app
  await app.start();

  console.log('⚡️ Bolt app is running!');
})();
