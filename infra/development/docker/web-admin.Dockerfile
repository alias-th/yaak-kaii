FROM node:20-alpine

WORKDIR /app

COPY web/yaak-kaii-admin/package*.json ./

RUN npm install

COPY web/yaak-kaii-admin ./

RUN npm run build

EXPOSE 3001

CMD ["npm", "start"]