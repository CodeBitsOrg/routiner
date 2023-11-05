## Run database (MySQL)
docker run --rm --name gogogo -v gogogo:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=gogogo --network net -p 3306:3306 -d mysql:8

## For Session start Redis
docker run --name redis -p 6379:6379 -d redis