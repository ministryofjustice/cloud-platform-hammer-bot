FROM node:20.10-alpine

WORKDIR /usr/src/app

COPY package*.json  .

COPY app.js .

RUN npm i

EXPOSE 3000

ENV HOST=0.0.0.0

RUN chown -R node:node /usr/src/app

USER node

CMD ["node", "/usr/src/app/app.js" ]
