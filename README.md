# forum

Features of web forum project:

+ communication between users.
+ associating posts to categories.
+ like and dislike posts and comments.
+ filtering posts.

> [Task objectives and audit](https://github.com/01-edu/public/tree/master/subjects/forum)

## to run:
1. Clone the repo:
>`https://01.kood.tech/git/Denis/forum.git`

2. If you having any troubles with for example creating post, check console log. If it saying <Session token expired>, just relogin 

### from terminal

> `docker build -t forum:latest .`
    
    First command will build project

> `docker run -p 8080:8080 -p 3000:3000 forum:latest`

    This command will launch forum

### or use Docker

To use Docker for launching the forum, Docker should be installed from [here](https://docs.docker.com/get-docker/)

 Follow the instructions. Wait for image to be built. 

2. Link to objectives and audit: [here](https://github.com/01-edu/public/tree/master/subjects/forum).
 
3. To run SQL tests:

    > `$ sqlite3 forumData.db 'SELECT * FROM users;'`

    > `$ sqlite3 forumData.db 'SELECT * FROM post;' `

    > `$ sqlite3 forumData.db 'SELECT * FROM comments;'`

### to note

Sqlite database may change permissions when moving the file (like downloading it). In that case, adding to the database will not work. After cloning the repo, please check the database permissions and set all to read and write.

## Implementation
- Backend: `Golang`
- Frontend: `HTML`, `CSS`, `JS` `TS`
- Database: `Sqlite3`
- Container: `Docker`

## Authors

[Stepan Tihinya](Discord: StepanTI), 
[Martin Vahe](Discord: mvahe), 
[Deniss Orlov](Discord: Denis), 
[Denis Petrov](Discord: Dolphin), 
