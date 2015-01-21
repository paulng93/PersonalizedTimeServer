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

//variables used 
var default_port = flag.String("port", "8080", "Default port number is 8080")

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
	http.HandleFunc("/", TimeServer)
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