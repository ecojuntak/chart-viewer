FROM node:latest as builder

WORKDIR /app

COPY package*.json ./
COPY ./ .
RUN npm install
RUN npm install -g @vue/cli

ARG VUE_APP_API_SERVER_HOST

RUN touch .env
RUN echo VUE_APP_API_SERVER_HOST=$VUE_APP_API_SERVER_HOST > .env

RUN npm run build

FROM nginx as dist

COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]