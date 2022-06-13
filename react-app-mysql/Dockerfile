FROM node:18-alpine3.14
WORKDIR /react/app
COPY ./app/package.json .
RUN npm install -g npm@latest
RUN npm install
ENV PORT=3003
EXPOSE 3003
CMD npm start
