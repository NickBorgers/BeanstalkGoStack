# Panda Health Reporting Tool

#### *It doesn't have to make sense, it just has to demonstrate a cloud-native architecture using AWS for CD*

[Product is running here](http://hellogo-env.mfrmtm2ahc.us-east-1.elasticbeanstalk.com/index.html)

## Technical summary
* Elastic Beanstalk for hosting
* Golang services
* Light, pure JavaScript UI

## Architecture
* Event-based
* Microservices architecture
* Stateless API, other than websockets
* Would not actually scale out in present form because analysis could be pushed to an instance without the requester attached; I should be using Pub/Sub for that part

![UML Sequence Diagram](/UML-Sequence.png)