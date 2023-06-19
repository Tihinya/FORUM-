read -p "Press enter to create a many posts :))))"
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "What an amazing new test post!",
  "content": "i predict that the id is going to be 501!",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Kickem"
  },
  "categories": ["beanz", "x", "heinz"]
}' -k https://localhost:8080/post/1
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "The Art of Cooking",
  "content": "Join me in a culinary adventure as we explore delicious recipes and culinary techniques.",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "master_ceil"
  },
  "categories": ["food", "tomcookery"]
}' -k https://localhost:8080/post/2
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "Finess meister tips and trickets",
  "content": "Learn effective finess strategies and get inspired to lead a healthy lifestyle.",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "finess_meisterman"
  },
  "categories": ["health", "finess"]
}' -k https://localhost:8080/post/3
echo
echo

read -p "Press enter to to create comments for post id 1"
curl -X POST -H "Content-Type: application/json" -d '{
  "content": "ONEEEEEEEEEEEEEEEE",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Kickem"
  }
}' -k https://localhost:8080/comment/1
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "content": "TWOOOOOOO",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Kickem2"
  }
}' -k https://localhost:8080/comment/1
echo

read -p "Press enter to to create comments for post id 2"
curl -X POST -H "Content-Type: application/json" -d '{
  "content": "SUUUUUUUUUUUUUUUUUUUUUUUUUU",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Nukem"
  }
}' -k https://localhost:8080/comment/2
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "content": "cristiano ronaldo",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "actualy cr"
  }
}' -k https://localhost:8080/comment/2
echo

read -p "Press enter to to update comment id 1"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "content": "testing for update IS THIS UPDATED AAAAAAAAAAAA"
}' -k https://localhost:8080/comment/1
echo

read -p "Press enter to delete comment id 1"
curl -X DELETE -k https://localhost:8080/comment/1
echo
echo "-------------------"

read -p "Press enter to FAIL comment creation for post id 240"
curl -X POST -H "Content-Type: application/json" -d '{
  "content": "SUUUUUUUUUUUUUUUUUUUUUUUUUU",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "Dick Nukem"
  }
}' -k https://localhost:8080/comment/240
echo

read -p "Press enter to FAIL updating comment id 69 (cuz it doesnt exist)"
curl -X PATCH -H "Content-Type: application/json" -d '{
  "content": "testing for update IS THIS UPDATED AAAAAAAAAAAA"
}' -k https://localhost:8080/comment/69
echo

read -p "Press enter to FAIL comment deletion id 79"
curl -X DELETE -k https://localhost:8080/comment/79
echo