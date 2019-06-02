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

### Component descriptions
The three services are copuled by message queue and event syntax. No versioning of events was implemented.
The UI is hosted outside of the Go microservices by NGINX itself.

* Frontend service
    * [Main file](/frontendservice.go)
    * [Build command](/buildspec.yml#L20)
    * [Beanstalk mapping](/Procfile#L1)
    * Takes in requests for information about pandas by-name via HTTP
    * Reads [JSON object](/pandahealthstructs.go#L10-L17) analysis results off analysis queue
    * Delivers analysis via websocket to all connected clients
    * Should be consuming analysis reports from a Pub/Sub channel
* Data service
    * [Main file](/dataservice.go)
    * [Build command](/buildspec.yml#L22)
    * [Beanstalk mapping](/Procfile#L2)
    * Takes in requests for information about pandas by-name via request queue
    * Makes up repeatable data about requested panda and pushes it onto data queue
    * Takes in arbitrary string names and outputs [JSON objects](/pandahealthstructs.go#L3-L8)
    * Uses Panda health data generator (should be a package):
        * [/pandahealthdataretrieval.go](/pandahealthdataretrieval.go)
        * [/pandahealthindicators.go](/pandahealthindicators.go)
* Analysis service
    * [Main file](/analysisservice.go)
    * [Build command](/buildspec.yml#L24)
    * [Beanstalk mapping](/Procfile#L3)
    * Takes in requests for analysis about pandas via data queue
    * Performs analysis based on information included in received message
    * Pushes results onto analysis queue as [JSON objects](/pandahealthstructs.go#L10-L17)
* UI
    * [Main file](/html/index.html)
    * [Passthrough during build](/buildspec.yml#L35-L36)
    * [Beanstalk mapping](/.ebextensions/go-settings.config#L3-L4)
    * Pure JavaScript and CSS because I don't do UI frameworks
    * Does nothing but trigger HTTP GETs, make Websocket connection, and [format/display JSON analysis](/html/index.html#L91-L124) reports returned via JSON
* Shared Go code
    * These should be in packages, but I don't know how to do that and am ready to take the rest of my Sunday off
    * [/pandahealthstructs.go](/pandahealthstructs.go)
    * [/sqs.go](/sqs.go)
* Things that should be configuration files instead
    * [/queuenames.go](/queuenames.go)
    * [/constants.go](/constants.go)
    