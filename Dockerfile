# Frontend stage
FROM node:18 as frontend-build
WORKDIR /app/client
COPY /client/package.json ./
RUN npm install
COPY /client ./
RUN npm run build

# Backend stage
FROM golang:1.20 as backend-build
WORKDIR /go/src/app
COPY /server/go.mod /server/go.sum ./
RUN go mod download
COPY /server ./
RUN GOOS=linux go build -o main .

# Finish stage
FROM node:18
WORKDIR /app
COPY --from=backend-build /go/src/app/main .
COPY /server/dev_config.json .
COPY /server/database/database.db ./database/database.db
COPY --from=frontend-build /app/client ./

EXPOSE 8080
EXPOSE 3000

# Running the container
COPY start.sh .
RUN chmod +x start.sh
CMD ["./start.sh"]