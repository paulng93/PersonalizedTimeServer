/**
 * @author Paul Nguyen
 * @Date: 1/13/15
 * @Name: timeserver.go
 * @Descrption: a simple time server that tells you time
 */
package main
//imports for all packages used
import "net/http"
import "log"
import "fmt"
import "flag"
import "os"
import "strconv"
import "time"
import "sync"
import "os/exec"
import "strings"
//import "net/url"
//import "bytes"

//variables used 
var default_port = flag.String("port", "8080", "Default port number is 8080")
var counter = struct{
	countersLock sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}
	/**
	 * Time server for when going to localhost.PORTNUMBER/time
	 */
func TimeServer(w http.ResponseWriter, req *http.Request) {
	// if user goes to another website after time/...
	if req.URL.Path != "/time/" {
		errorHandler(w, req, http.StatusNotFound)
		return
	}
	//html code
	fmt.Fprint(w, "<html><head><style> p{font-size: xx-large}")
	fmt.Fprint(w, "span.time {color:red}")
	fmt.Fprint(w, "</style></head><body><p> The time is now ")
	fmt.Fprint(w, "<span class=\"time\">")
	fmt.Fprint(w, time.Now().Format("3:04:04 PM"))
	fmt.Fprint(w, "</span>.</p></body></html>")
}
// if user goes to / or index.html
func home(w http.ResponseWriter, req *http.Request){
	//if there is a cookie
	/*
	_, ok := counter.m[GlobalName]
	if (ok){
		fmt.Fprint(w, "<html><body><p> Greetings,")
		fmt.Fprint(w, GlobalName)
		fmt.Fprint(w, "</p></body></html>")
		*/
	if(false){
	}else {
		//if there isn't a cookie yet
		fmt.Fprint(w, "<html><body><p><form action=\"login\"> What is your name, Earthling?")
		fmt.Fprint(w, "<input type=\"text\" name=\"name\" size=\"50\">")
		fmt.Fprint(w, "<input type=\"submit\"></form></p></body></html>")
	}
}

func logout(w http.ResponseWriter, req *http.Request){
	//if there isn't a cookie yet
	//TODO: REMOVE FROM COOKIE MAP
	if req.URL.Path != "/logout/" {
		errorHandler(w, req, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "<html><head><META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprint(w, "<body><p>Good-bye.</p></body></html>")
}

// login
func loginHandler(w http.ResponseWriter, req *http.Request){
	name := req.FormValue("name")
	redirectTarget :="/"
	if name != "" {
		value := getUniqueValue()
		tempValue := string(value[:])
		fmt.Println(tempValue)
		
		cookie := &http.Cookie{
			Name: name,
			Value: strings.Trim(tempValue, "\n"),
		}
		http.SetCookie(w,cookie)
		
		counter.countersLock.RLock()
		counter.m[strings.Trim(tempValue, "\n")] = name
		counter.countersLock.RUnlock()
	}else {
		fmt.Fprint(w, "<html><body><p> C'mon, I need a name.")
		fmt.Fprint(w, "</p></body></html>")
	}
	http.Redirect(w,req,redirectTarget, 302)
	// check if name is in cookie map
	//
}
func getUniqueValue() []byte{
	out, error := exec.Command("uuidgen").Output()
	if error != nil {
		log.Fatal(error)
	}
	return out
}

// if user goes to different site localhost.PORTNUMBER/..
func homeHandler(w http.ResponseWriter, req *http.Request){
	errorHandler(w, req, http.StatusNotFound)
}

// This function handles errors and goes to page
func errorHandler(w http.ResponseWriter, req *http.Request, status int){
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "<html><body><p> These are not the URLs you're looking for.")
		fmt.Fprint(w, "</p></body></html>")
	}
}
// simple function that appends a colon to the string passed in
func appendColon(temp string)string {
	temp = ":"+ temp
	return temp
}
// checks the port to make sure that it's not a reserved port
// returns: a boolean value 
func checkPort()bool {
	i, err := strconv.Atoi(*default_port)
	if err != nil {
		fmt.Println(err)
	}
	if i < 1024 {
		return false
	} else {
		return true
	}
}


func main() {
	//create server
	flag.Parse()
	http.HandleFunc("/time/", TimeServer)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", home)
	http.HandleFunc("/index.html/", home)
	http.HandleFunc("/logout/", logout)
	fmt.Println("SERVER ONLINE")
	if !checkPort() {
		fmt.Printf("Error trying to connect to privledged port\n")
		os.Exit(404)
	}
	err := http.ListenAndServe(appendColon(*default_port), nil)
	if err != nil {
		fmt.Printf("Error\n")
		log.Fatal("ListenAndServe: ", err)
		os.Exit(404)		
	}	
}