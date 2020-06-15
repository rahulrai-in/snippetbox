docker run --name snippetbox-db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.0.20

docker container cp seed.sql snippetbox-db:/

docker exec -it snippetbox-db mysql -u root -ppassword -e "CREATE DATABASE IF NOT EXISTS snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

docker exec -it snippetbox-db mysql -u root -ppassword snippetbox 

REM source ./seed.sql