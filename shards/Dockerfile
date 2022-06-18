FROM gomicro/goose:latest

WORKDIR /usr/src
COPY . ./

ADD *.sql ./shard_migrations/
ADD entrypoint.sh ./shards/

RUN chmod +x ./shards/entrypoint.sh

ENTRYPOINT ["./shards/entrypoint.sh"]