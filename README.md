## REDDIT:

As part I of this project, you need to build an engine that (in part II) will be paired up with REST API/WebSockets to provide full functionality.  The specific API supported by Reddit can be found at:https://www.reddit.com/dev/api/. In part I you are only building the engine and a simulator, not the API + web clients. An overview of Reddit and its functionality can be found at: https://www.oberlo.com/blog/what-is-reddit

Specific things you have to do are: 

### Implement a Reddit-like engine with the following functionality:
Register account<br/>
Create & join sub-reddit; leave sub-reddit<br/>
Post in sub-reddit. Make the posts just simple text. No need to support images or markdown.<br/>
Comment in sub-reddit. Keep in mind that comments are hierarchical (i.e. you can comment on a comment)<br/>
Upvote+downvote + compute Karma<br/>
Get feed of posts<br/>
Get list of direct messages; Reply to direct messages<br/>

### Implement a tester/simulator to test the above
Simulate as many users as you can<br/>
Simulate periods of live connection and disconnection for users<br/>
Simulate a Zipf distribution on the number of sub-reddit members. For accounts with a lot of subscribers, increase the number of posts. Make some of these messages re-posts<br/>

### Other considerations:
The client part (posting, commenting, subscribing) and the engine (distribute posts, track comments, etc) have to be in separate processes. Preferably, you use multiple independent client processes that simulate thousands of clients and a single-engine process <br/>
You need to measure various aspects of your simulator and report performance <br/>
More detail in the lecture as the project progresses.<br/>

* Implement REST API interface for the engine you designed in project 4.1. Use structure similar to Redit's official API (does not need to be
identical). <br/>
* Implement a simple client that uses the REST API to perform each piece of functionality you supoprt <br/>
* Run your engine with multiple clients to show that functionality works (make video for demo) <br/>
