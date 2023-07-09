read -p "Press enter to create a many posts :))))"
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "What an amazing new test post!",
  "content": "i predict that the id is going to be 501!",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "Dick Kickem"
  },
  "categories": ["beanz", "x", "heinz"]
}' -k https://localhost:8080/post
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "The Art of Cooking",
  "content": "Join me in a culinary adventure as we explore delicious recipes and culinary techniques.",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "master_ceil"
  },
  "categories": ["food", "tomcookery", "beanz"]
}' -k https://localhost:8080/post
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "Finess meister tips and trickets",
  "content": "Learn effective finess strategies and get inspired to lead a healthy lifestyle.",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "finess_meisterman"
  },
  "categories": ["beanz", "test", "hello"]
}' -k https://localhost:8080/post
echo
echo

read -p "Press enter to modify your post!"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "title": "NOT SO AMAZING NOW, EH",
  "content": "caps is annoying",
  "categories": ["food", "dude", "beanz", "tomcookery", "LONELY SURVIVOR"]
}' -k https://localhost:8080/post/3
echo

read -p "Press enter to nuke post id 3"
curl -X DELETE -k https://localhost:8080/post/3
echo