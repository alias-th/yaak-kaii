FROM node:20-alpine

WORKDIR /app

COPY web/yaak-kaii-public/package*.json ./

RUN npm install

COPY web/yaak-kaii-public ./

RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]