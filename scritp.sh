docker build -t stepanti/real-time-forum-backend:latest ./server
docker build -t stepanti/real-time-forum-frontend:latest ./client
docker push stepanti/real-time-forum-backend
docker push stepanti/real-time-forum-frontend