# twittervotes

## Overview

This project can be found in Chapter 5 and 6 of the Go Programming Blueprints text book written by Mat Ryer. 

All I have done is typed the code from the textbook, managed to get it working, refactored it in a way I thought was clearer and have learnt a ton along the way.

This project covers the concepts of building a distributed system and exposing data and user functionalities over a RESTful API.

Massive credit goes to Mat Ryer for the awesome textbook!

## The project

The project uses Twitter's 'tracking' API to track hashtagged words that Twitter user's include in their tweets. 
The idea is, that the user will choose random words, and then this system will track them on twitter and will expose the total count of these words in a pie chart on your browser. 
The user can then see what hashtagged words are popular throughout twitter!

## How it works


## Dependencies
Obviously Go must be installed on your computer. 
Mongodb will always need to be installed.
After that, a ```go get github.com/corno93/twittervotes``` will install this code for you locally. 
You will need to have Twitter developer credentials. These will need to be stored as environment variables on your machine.

## How to run
Open 7 command prompts and navigate your way to the go project.
Type in separate prompts:
- ```mongod```
- ```nsqlookupd```
- ```nsqd --lookupd-tcp-address=127.0.0.1:4160```
- ```go build -o twittervotes``` then ```./twittervotes```
- ```go build -o counter``` then ```./counter```
- ```go build -o api``` then ```./api```
- ```go build -o web``` then ```./web```

After go to ```http://localhost:8081/```



### Coming soon - cloud implementation...
