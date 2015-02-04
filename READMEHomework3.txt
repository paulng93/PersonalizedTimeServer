Paul Nguyen
2/3/15
ReadMe instructions on how to run a refactored PersonalizedTimeServer.

IMPORTANT READ ON HOW TO SET UP ENVIRONMENT BEFORE RUN!
What I will include will be a folder called template containing all the html templates that will be rendered.
A CSS file, and PersonalizedTimeServer which will all be inside the src folder of Go directory. 

For my program I run it using go install personalizedTimeServer command which creates a executable file stored in the bin 
folder in go directory. Because of that, and not running it using the command go run personalizedTimeServer.go, i have to specify 
the location of the CSS file exsplicity in my program and not use the relative path. 

The two main paths that you will have to set up for the run environment are the absolute path of the CSS file in line
24 and 25 of my program. for me the path is 'cssPath = "C:/Users/Paul/Documents/Home/go/src/css"' and the path of where the template 
are 'temp_Location = "Home/go/src/templates"'. The path of the tempates is that becuase i run it from my command line and my command
line is in the diretory MyDocuments. The location from mydocuments aka where your running it all is C:/Users/Paul/Documents/Home/go/src/PersonalizedTimeServer
You can specify the location of the templates using --template flag below. The location must be relative from where your 
console is (where your running it) and where the templates are stored.

MAKEFILE:
i have included a make file. it will run the program if you are in the folder before gopath. EX: i run it from C:/Users/Paul/Documents
echo $GOPATH = c:/Users/Paul/Documents/Home/go
to set that up you can use export $GOPATH= "INSERT DESIRED GOPATH" 
then use 
export PATH=$PATH:$GOPATH/bin

HELP: to set up environment go to: https://golang.org/doc/code.html

Once the location constants are set up we can look at the flags that are set up. 

FLAGS:
--port: can change port number used for server must be > 1024. 			Default is 8080
--V: change the version number of server					Default: Version 2
--templates: change the location of where the program looking for the templates	Default: "Home/go/src/templates"
--log: change the location of where the the logfile is saved. 			Default: "./ect/timeserver.log"

Notes: Using log.seelog to log error but i also kept in the http.Errors as well. 

DIRECTORY LOCATION:
ect/  //This is where i'm saving the log files but you can change this by using flag
Home/
	go/
		src/
			templates/		location of all templates
			PersonalizedTimeServer/	main folder
			CookieJar/		Taking cookies out
			css/			css file
			gothub.com/ 		location of seelog
		bin/
			location of executable file after run make file

FILES INCLUDED:
	I will zip my entire go library in a tar file. Please ignore following directory
	/hello
	/timeserver
	
