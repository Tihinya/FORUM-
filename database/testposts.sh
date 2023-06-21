read -p "Press enter to create a new post, don't forget to refresh ur browser :)"
curl -X POST -H "Content-Type: application/json" -d '{
  "title": "What an amazing new test post!",
  "content": "i predict that the id is going to be 6!",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Kickem"
  },
  "categories": ["beanz", "x", "heinz"]
}' -k https://localhost:8080/post
echo

read -p "Press enter to modify your post!"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "title": "NOT SO AMAZING NOW, EH",
  "content": "caps is annoying",
  "categories": ["updated", "x", "upgraded"]
}' -k https://localhost:8080/post/1
echo

read -p "Press enter to nuke your amazing post"
curl -X DELETE -k https://localhost:8080/post/1
echo