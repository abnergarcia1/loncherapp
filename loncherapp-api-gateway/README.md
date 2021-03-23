# How to run locally


###Prerequisites:

- Docker
- Go v1.14

###Steps
Run `make run-local` inside folder in a new terminal and wait for the message `[GIN-debug] Listening and serving HTTP on :8081`

The "make" command runs a Redis image with access in `localhost:6380` 

if you receive the error `ERRO[0000] Error when trying to get redis db: dial tcp [::1]:6380: connect: connection refused  RedisClient="Redis<localhost:6380 db:0>"`
try to delete the current redis image in docker and rerun the make command