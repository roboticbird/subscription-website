FROM node:13.12.0-alpine as build
WORKDIR $GOPATH/src/subscription-website/frontend-react
ENV PATH /frontend-react/node_modules/.bin:$PATH

COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent
RUN npm install axios
RUN npm install react-bootstrap

COPY . ./

CMD ["npm", "start"]


