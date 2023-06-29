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
  "categories": ["health", "finess", "food"]
}' -k https://localhost:8080/post
echo
echo

read -p "Press enter to to create comments for post id 1"
curl -X POST -H "Content-Type: application/json" -d '{
  "content": "ONEEEEEEEEEEEEEEEE",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "Dick Kickem"
  }
}' -k https://localhost:8080/comment/1
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "content": "TWOOOOOOO",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "Dick Kickem2"
  }
}' -k https://localhost:8080/comment/1
echo

read -p "Press enter to to create comments for post id 2"
curl -X POST -H "Content-Type: application/json" -d '{
  "content": "SUUUUUUUUUUUUUUUUUUUUUUUUUU",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
    "username": "Dick Nukem"
  }
}' -k https://localhost:8080/comment/2
echo

curl -X POST -H "Content-Type: application/json" -d '{
  "content": "cristiano ronaldo",
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
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
  "user_info": {
    "avatar": "https://example.com/profile_picture.png",
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