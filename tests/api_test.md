
## How to test the RESTful API

We can test the API by using curl by creating GET and POST messages. 

# Add data to the db
In separate command prompts start a mongo instance and open up a mongo shell.
Enter in some data:
Create a database called ballots, a collection called polls and enter some fake data:
```use ballots```
```db.polls.insert({"title":"Test poll","options":["happy","sad","fail","win"]})```



# Test Get
Get all the #hashtag words stored in the collection.

Start the server handling the api
```go build -o api```
``` ./api```
On a separate command prompt enter:
```curl -X GET http://localhost:8080/polls/?key=abc123```

The response:
```[{"id":"541728728ea48e5e5d5bb18a","title":"Test poll", "options": ["happy","sad","fail","win"]"apikey":"abc123"}]```
             
# Test Post
Create a new set of #hasttag words to track.

Enter:
```curl --data '{"title":"test","options":["john","paul","george","ringo"]}' -X POST http://localhost:8080/polls/?key=abc123```

Check if its there by doing a Get:
```curl -X GET http://localhost:8080/polls/?key=abc123```

The response:
```[{"id":"541728728ea48e5e5d5bb18a","title":"Test poll", "options": ["happy","sad","fail","win"]}, {"id":"5bfda6642c9b61459d9b68dd","title":"test","options":["john","paul","george","ringo"],"apikey":"abc123"}]```


**Note (Postman)[https://www.getpostman.com/] is a great alternative if you are regularly testing RESTful APIs, providing a nice interface to modify your request and to save requests for future use. Plus it works on windows**