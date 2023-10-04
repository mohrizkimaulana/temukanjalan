FROM node:14

WORKDIR /app_node

COPY ./app_node/package*.json ./

COPY ./app_node/node_modules ./node_modules

RUN npm install

COPY . .

EXPOSE 3000

CMD ["node", "./app_node/app.js"]
