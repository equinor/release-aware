FROM node:12-alpine
RUN npm install -g serve
WORKDIR /code
COPY package.json /code/
RUN yarn --no-lockfile
COPY ./ /code/
RUN yarn build
EXPOSE 3000
CMD ["yarn", "start"]
