#!/bin/bash

# Successful registration
read -p "Press Enter to send a request for successful registration..."
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

curl -X POST -H "Content-Type: application/json" -d '{
  "username": "erik",    
  "email": "erik@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

curl -X POST -H "Content-Type: application/json" -d '{
  "username": "andrei",
  "email": "andrei@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

# View all users
read -p "Press Enter to View all users"
curl -X GET -k https://localhost:8080/users/get

echo "----------------------------"

# Missing email and password
read -p "Press Enter to send a request with missing email and password..."
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

# Invalid email format
read -p "Press Enter to send a request with an invalid email format..."
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "johnexample.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

# Email or username already taken
read -p "Press Enter to create user"
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "jane_doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

read -p "Press Enter to send a request with a username already taken..."
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "jane@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/user/create

echo "----------------------------"

# Get user by ID
read -p "Press Enter to get user by ID 1"
curl -X GET -k https://localhost:8080/user/1/get

echo "----------------------------"

# Update user by ID
read -p "Press Enter to update user by ID 1"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "username": "john_smith",
  "email": "john@example.com"
}' -k https://localhost:8080/user/1/update

echo "----------------------------"

# Get user by ID
read -p "Press Enter to get user by ID 1"
curl -X GET -k https://localhost:8080/user/1/get

echo "----------------------------"

# Delete user by ID
read -p "Press Enter to delete user by ID 1"
curl -X DELETE -k https://localhost:8080/user/1/delete

echo "----------------------------"

# View all users
read -p "Press Enter to View all users"
curl -X GET -k https://localhost:8080/users/get

# Delete user by ID
read -p "Press Enter to delete user by ID 2"
curl -X DELETE -k https://localhost:8080/user/2/delete

echo "----------------------------"

# View all users
read -p "Press Enter to View all users"
curl -X GET -k https://localhost:8080/users/get

# Delete user by ID
read -p "Press Enter to delete user by ID 3"
curl -X DELETE -k https://localhost:8080/user/3/delete

echo "----------------------------"

# View all users
read -p "Press Enter to View all users"
curl -X GET -k https://localhost:8080/users/get
