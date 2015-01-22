/**
 * @author Paul Nguyen
 * @Date: 1/21/15
 * @Name: PersonalizedTimeServer.go
 * @Descrption: a more complex server that keeps track of cookies, tells time, and 
 * 				logs you out
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

//variables used 
var default_port = flag.String("port", "8080", "Default port number is 8080")
var counter = struct{
	countersLock sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}
//--------------------------------------------------------------------------------------
/**
 * TimeServer function is the recipent when user goes to /time
 * Parameter Responsewriter and http.Request
 * 
 */
func TimeServer(w http.ResponseWriter, req *http.Request) {
	// if user goes to another website after time/...
	if req.URL.Path != "/time/" {
		errorHandler(w, req, http.StatusNotFound)
		return
	}
	cookie, _ := req.Cookie("UUID")
	existCheck := false
	temp2 := ""
	if cookie != nil { // if cookie exist set flags
		name, check := counter.m[cookie.Value]
		fmt.Println(name)
		existCheck = check
		temp2 = name
		fmt.Println(check)
	}
	//html code
	fmt.Fprint(w, "<html><head><style> p{font-size: xx-large}")
	fmt.Fprint(w, "span.time {color:red}")
	fmt.Fprint(w, "</style></head><body><p> The time is now ")
	fmt.Fprint(w, "<span class=\"time\">")
	fmt.Fprint(w, time.Now().Format("3:04:04 PM"))
	if existCheck { //prints the name if cookie exist
		fmt.Fprint(w, "</span>, ")
		fmt.Fprint(w, temp2)
		fmt.Fprint(w, "</p></body></html>")
	} else {
		fmt.Fprint(w, "</span>.</p></body></html>")
	}
}
//--------------------------------------------------------------------------------------
/**
 * home function is the recipent when user goes to / or /index.html
 * Parameter Responsewriter and http.Request
 * if user goes to home
 */
func home(w http.ResponseWriter, req *http.Request){
	//if there is a cookie
	//grabbing name from broswer
	cookie,_ := req.Cookie("UUID")
	temp := false
	//checking of cookie exist
	if cookie != nil {
		_, ok := counter.m[cookie.Value]
		temp = ok
	}
	//only true if cookie exist AND is the correct cookie (in map)
	if(temp){
		i, ok := counter.m[cookie.Value]
		if ok && i != "" {//last check to see if name exist
			fmt.Fprint(w, "<html><body><p> Greetings, ")
			fmt.Fprint(w, i)
			fmt.Fprint(w, "</p></body></html>")
		}
	}else {
		//if there isn't a cookie yet ask user for name and redirect to login
		fmt.Fprint(w, "<html><body><p><form action=\"login\"> What is your name, Earthling?")
		fmt.Fprint(w, "<input type=\"text\" name=\"name\" size=\"50\">")
		fmt.Fprint(w, "<input type=\"submit\"></form></p></body></html>")
	}
}
//--------------------------------------------------------------------------------------
/**
 * logout function is the recipent when user goes to /logout
 * Parameter Responsewriter and http.Request
 * this will log the user out aka destroy the internal record of the cookie
 */
func logout(w http.ResponseWriter, req *http.Request){
	//if user goes to wrong url
	if req.URL.Path != "/logout/" {
		errorHandler(w, req, http.StatusNotFound)
		return
	}//getting cookie
	cookie, _ := req.Cookie("UUID")
	if cookie != nil {// cookie check
		_, ok := counter.m[cookie.Value]
		if ok {
			counter.countersLock.RLock()
			delete(counter.m,cookie.Value)
			counter.countersLock.RUnlock()
		}
	}
	fmt.Fprint(w, "<html><head><META http-equiv=\"refresh\" content=\"10;URL=/\">")
	fmt.Fprint(w, "<body><p>Good-bye.</p></body></html>")
}
//--------------------------------------------------------------------------------------
/**
 * login function is the recipent when user goes to /time
 * Parameter Responsewriter and http.Request
 * This will log the user in
 */
func loginHandler(w http.ResponseWriter, req *http.Request){
	name := req.FormValue("name")
	redirectTarget :="/"
	if name != "" { //make sure user has input a name
		value := getUniqueValue() // generate UUID
		tempValue := string(value[:]) // turn into string
		//creating cookie
		cookie := &http.Cookie{
			Name: "UUID",
			Value: strings.Trim(tempValue, "\n"),
			Path: "/",
		}
		//set cookie
		http.SetCookie(w,cookie)
		
		counter.countersLock.RLock()
		counter.m[strings.Trim(tempValue, "\n")] = name
		counter.countersLock.RUnlock()
	}else {// if user has not input a name
		fmt.Fprint(w, "<html><body><p> C'mon, I need a name.")
		fmt.Fprint(w, "</p></body></html>")
	}
	http.Redirect(w,req,redirectTarget, 302)
}
//--------------------------------------------------------------------------------------
/**
 * GetUniqueValue : generate a UUID
 * Parameter none
 * Renturn byte array
 */
func getUniqueValue() []byte{
	out, error := exec.Command("uuidgen").Output()
	if error != nil {
		log.Fatal(error)
	}
	return out
}
//--------------------------------------------------------------------------------------
/**
 * ErrorHandler function is the recipent when user goes to a wrong url
 * Parameter Responsewriter and http.Request
 * 
 */
func errorHandler(w http.ResponseWriter, req *http.Request, status int){
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "<html><body><p> These are not the URLs you're looking for.")
		fmt.Fprint(w, "</p></body></html>")
	}
}
//--------------------------------------------------------------------------------------
// simple function that appends a colon to the string passed in
func appendColon(temp string)string {
	temp = ":"+ temp
	return temp
}
//--------------------------------------------------------------------------------------
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
//--------------------------------------------------------------------------------------

func main() {
	//create server
	flag.Parse()
	http.HandleFunc("/time/", TimeServer)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/index.html/login", loginHandler)
	http.HandleFunc("/", home)
	http.HandleFunc("/index.html", home)
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