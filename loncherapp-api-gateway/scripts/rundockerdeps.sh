   docker stop loncherapp-redis
    docker rm loncherapp-redis
    docker run -d --name loncherapp-redis -p 6380:6379  redis

    docker stop loncherapp-mariadb
    docker rm loncherapp-mariadb
    docker run -d -p 3307:3306 --name loncherapp-mariadb -e \
    MYSQL_ROOT_PASSWORD=password1234 -e MYSQL_DATABASE=loncherapp \
    -v "$(pwd)"/scripts/sql:/docker-entrypoint-initdb.d \
     mariadb:latest
