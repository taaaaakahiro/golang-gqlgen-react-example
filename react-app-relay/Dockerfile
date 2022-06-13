FROM node:18-alpine3.14
WORKDIR /react/app
COPY ./app/package.json .
RUN npm install -g npm@8.9.0
RUN npm install
ENV PORT=3002
EXPOSE 3002
CMD npm start
