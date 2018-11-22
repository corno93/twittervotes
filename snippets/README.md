
Instructions on how to run:

Two NSQ terminals:
nsqlookupd
nsqd --lookupd-tcp-address=127.0.0.1:4160
nsq_tail --topic="votes" --lookupd-http-address=localhost:4161 (check the messages on the topic)


Mongo:
mongod

Twitter:
go build -o twittervotes
./twittervotes

