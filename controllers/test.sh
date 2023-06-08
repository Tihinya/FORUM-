read -p "Press enter to use POST method for post creation id 1:"
curl -X POST -H "Content-Type: application/json" -d '{
  "title": "Test post",
  "text": "This is the content of my post",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "john_doe"
  },
  "likes": 0,
  "dislikes": 0,
  "categories": ["technology", "programming"]
}' -k https://localhost:8080/post/1
echo

read -p "Press enter to use POST method for post creation id 2:"
curl -X POST -H "Content-Type: application/json" -d '{
  "title": "2222222222",
  "text": "2 This is the content of my post 2",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "john_doe"
  },
  "categories": ["technology", "programming"]
}' -k https://localhost:8080/post/2
echo

read -p "Press enter to use GET method for reading post id 1"
curl -X GET -k https://localhost:8080/post/1
echo

read -p "Press enter to use GET method for reading post id 2"
curl -X GET -k https://localhost:8080/post/2
echo

read -p "Press enter to use PATCH method for updating post id 1"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "title": "UPDATED UPDATED UPDATED",
  "content": "Updated Updated Updated?",
  "categories": ["updated", "the whats?"]
}' -k https://localhost:8080/post/1
echo

read -p "Press enter to use GET method for reading post id 1"
curl -X GET -k https://localhost:8080/post/1
echo

read -p "Press enter to use GET method for reading post id 2"
curl -X GET -k https://localhost:8080/post/2
echo

read -p "Press enter to use DELETE method for deleting post 1"
curl -X DELETE -k https://localhost:8080/post/1
echo

read -p "Press enter to use GET method for reading post 1"
curl -X GET -k https://localhost:8080/post/1
