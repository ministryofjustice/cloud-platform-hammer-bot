FROM node:20.10-alpine

RUN addgroup -g 1000 -S appgroup && \
  adduser -u 1000 -S appuser -G appgroup

WORKDIR /usr/src/app

COPY slackbot/package*.json  .

COPY slackbot/app.js .

RUN npm i

EXPOSE 3000

ENV HOST=0.0.0.0

RUN chown -R node:node /usr/src/app

USER 1000

CMD ["node", "/usr/src/app/app.js" ]
