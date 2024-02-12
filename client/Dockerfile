# Frontend stage
FROM node:18 as builder

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM nginx:alpine

WORKDIR /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY ssl/server.key ssl/server.crt /etc/nginx/

COPY --from=builder /app/dist .

EXPOSE 443

ENTRYPOINT ["nginx", "-g", "daemon off;"]
