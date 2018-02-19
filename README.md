# crawlers

i found the assignment in my inbox on thursday and chose 'go' for the implementation. why 'go'? for one we had discussed it in our meeting last monday and also it’s the only one i had never heard of! fun times!

note: i’m sure there are improvements that can be made in the code (i have 3 days expertise with go!), but overall i found go to be very cool and useful, intuitive and mostly pain free. i will add go to my resume and recommend it! ... the html templating is very handy and easy to use. other things too.

is the go http server suitable for production environments? not sure, but possibly, ---need to look into this. i also have questions about error handling in go, as in, i don’t know enough about it to say for sure if im doing things properly (will research design patterns/best practices in the next few days)

the design is pretty simple and matches the req (i am big proponent of simple, effective solutions, sans extra bells and whistles).

arch (in a nutshell): html page --->go http server/backend logic --->MySQL database.

of note, design-wise: i’m using the 'count' column in the 'crawler' table as a double duty value-add (hence it being a varchar). when things go wrong during a search, this column gets updated with a description of the error. it’s an easy/useful way to keep things simple and show the user some good data. try typing in https://www.goggle.com in the app and you'll see how this works.

installation instructions:
1.	get the go package and install it as per instructions found here: https://golang.org/doc/install
2.	set up a mySQL database with root/admin as username/password.
3.	create a schema called netapp.
4.	execute the crawler.sql file against the netapp schema --- as you've guessed, this will create the single table required.
5.	the server.exe file should work, as-is. launch it!ive only tested this from the command line – recommend you do that same.
6.	if you want to mess with the code and build it yourself, you will need the ‘go-sql-driver’ - i had to install this using a go command

run the following command from a command prompt in the project root directory: 
    go get -u github.com/go-sql-driver/mysql 
(i think that’s what i did, not 100% sure on this, but pretty sure).

the infrastructure should be all setup at this point.

to test and see the code, pull the whole thing from GitHub and put it wherever you like. in my setup, the actual project root is found at c:/go-work/source/crawler.

the files:
- crawler.sql --- ddl for the table creation
- server.exe --- the windows executable
- server.go --- the go source code
- templates/crawler.html --- the html view

as mentioned, server.exe should work on its own, as-is.

if you want to build this yourself, open a command window in the server directory and type: go build ... this will regenerate the server.exe.

to test: start the application, open a browser and type this url: http://localhost:8080/home. voila!

ps - i didnt get the part where i was supposed to use git to track changes, i can tell you all the steps i took to build this though. it started out as a simple 'hello world' http server, --- the change log would look pretty silly anyhow!!!

here is the basic genesis:
- http server that prints 'hello world'
- add code that gets a hard-coded website. get the body, turn it into a string, get the length of the string.
- add the database stuff.
- create an html page, add the go template stuff.
- do some cleanup and run some tests. more cleanup.

i may make some additional changes --- i'll be sure to check them in this time :). overall, this was a lot of fun! thanks!

oh just fyi --- i started this on thursday night with zero knowledge of go and all the 'initial upload' code was complete by last night (saturday). i did read about go on thursday before selecting it. cheers!
