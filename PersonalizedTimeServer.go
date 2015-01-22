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
	cookie, _ := req.Cookie("UUID")
	name, ok := counter.m[cookie.Value]
	//html code
	fmt.Fprint(w, "<html><head><style> p{font-size: xx-large}")
	fmt.Fprint(w, "span.time {color:red}")
	fmt.Fprint(w, "</style></head><body><p> The time is now ")
	fmt.Fprint(w, "<span class=\"time\">")
	fmt.Fprint(w, time.Now().Format("3:04:04 PM"))
	if ok {
		fmt.Fprint(w, "</span>, ")
		fmt.Fprint(w, name)
		fmt.Fprint(w, "</p></body></html>")
	} else {
		fmt.Fprint(w, "</span>.</p></body></html>")
	}
}
// if user goes to / or index.html
func home(w http.ResponseWriter, req *http.Request){
	//if there is a cookie
	//grabbing name from broswer
	cookie, _ := req.Cookie("UUID")
	fmt.Println("cookie ")
	fmt.Println(cookie.Value)
	i, ok := counter.m[cookie.Value]
	fmt.Println("Name: " + i)
	fmt.Println(ok)
	if ok && i != "" {
		fmt.Fprint(w, "<html><body><p> Greetings, ")
		fmt.Fprint(w, i)
		fmt.Fprint(w, "</p></body></html>")
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
	cookie, _ := req.Cookie("UUID")
	_, ok := counter.m[cookie.Value]
	if ok {
		counter.countersLock.RLock()
		delete(counter.m,cookie.Value)
		counter.countersLock.RUnlock()
	}
	fmt.Fprint(w, "<html><head><META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprint(w, "<body><p>Good-bye.</p></body></html>")
}

// login
func loginHandler(w http.ResponseWriter, req *http.Request){
	name := req.FormValue("name")
	fmt.Println("path: " + req.URL.Path)
	redirectTarget :="/index.html"
	if req.URL.Path == "/login" {
		redirectTarget = "/"
	}
	if name != "" {
		value := getUniqueValue()
		tempValue := string(value[:])
		fmt.Println("login handler" + tempValue)
		
		cookie := &http.Cookie{
			Name: "UUID",
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
	http.HandleFunc("/index.html/login", loginHandler)
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