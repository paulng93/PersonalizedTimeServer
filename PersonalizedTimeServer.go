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
import "fmt"
import "flag"
import "os"
import "strconv"
import "time"
import "CookieJar"
import log "github.com/cihub/seelog"
import "html/template"
import "path"
import "strings"

const (
	logFilePath = "./ect/"
	cssPath = "C:/Users/Paul/Documents/Home/go/src/css"
	temp_Location = "Home/go/src/templates"
	versionNo = "Version No 3"
	)
//flag variables 
var default_port = flag.String("port", "8080", "Default port number is 8080")
var version = flag.String("V", "No", "Change to 'Yes' if you want to print")
//DEFAULT location set for my working environment 
var template_Location = flag.String("templates",temp_Location, "This allows to find location of templates" )
var logfile_Name = flag.String("log",logFilePath + "timeserver.log", " to specify the name of the log configuration file" )
//Cookie jar for taking cookie out 
var cookieJar = CookieJar.NewCookieJar()
//profile structure is needed to be passed in for templates 
// has 2 value name and current time
type Profile struct {
	Name string
	CurrentTime string
	LogPath string
}//&Profile{LogPath : name}
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
	cookie := cookieJar.GetCookie(req, "UUID")
	profile := Profile{"",time.Now().Format("3:04:04 PM"),""}
	if cookie != nil { // if cookie exist set flags
		name, check := cookieJar.GetValue(cookie.Value)
		profile = Profile{name,time.Now().Format("3:04:04 PM"),""}
		//logging info
		log.Info("Persons name is " + name)
		value := "no"
		if check {
			value = "yes"
		}
		log.Info("Name Exist? " + value)
	}
	//for templates
	lp := path.Join(*template_Location, "layout.html")
	fp := path.Join(*template_Location, "time.html")
	headerPath := path.Join(*template_Location, "menu.html")
	tmpl, err := template.ParseFiles(lp, fp, headerPath)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, profile); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	cookie := cookieJar.GetCookie(req, "UUID")
	checkCookie := false
	lp := path.Join(*template_Location, "layout.html")
	//checking of cookie exist
	if cookie != nil {
		_, check := cookieJar.GetValue(cookie.Value)
		checkCookie = check
	}
	//only true if cookie exist AND is the correct cookie (in map)
	if(checkCookie){
		name, ok := cookieJar.GetValue(cookie.Value)
		if ok && name != "" {//last check to see if name exist
			//profile = Profile{name, ""}
			fp := path.Join(*template_Location, "greeting.html")
			headerPath := path.Join(*template_Location, "menu.html")
			tmpl, err := template.ParseFiles(lp, fp, headerPath)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := tmpl.Execute(w, &Profile{Name : name}); err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}else {
		//if there isn't a cookie yet ask user for name and redirect to login
		fp := path.Join(*template_Location, "loginform.html")
		headerPath := path.Join(*template_Location, "menu.html")
		tmpl, err := template.ParseFiles(lp, fp, headerPath)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	}//deleting cookie
	cookie := cookieJar.GetCookie(req, "UUID")
	if cookie != nil {
		cookieJar.DeleteCookie(w, cookie.Value)
	}
	lp := path.Join(*template_Location, "layout.html")
	fp := path.Join(*template_Location, "logout.html")
	headerPath := path.Join(*template_Location, "menu.html")
	tmpl, err := template.ParseFiles(lp, fp, headerPath)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		cookieJar.CreateCookie(w, name)
	}else {// if user has not input a name
		lp := path.Join(*template_Location, "layout.html")
		fp := path.Join(*template_Location, "noName.html")
		headerPath := path.Join(*template_Location, "menu.html")
		tmpl, err := template.ParseFiles(lp, fp, headerPath)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w,req,redirectTarget, 302)
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
		lp := path.Join(*template_Location, "layout.html")
		fp := path.Join(*template_Location, "notFound.html")
		headerPath := path.Join(*template_Location, "menu.html")
		tmpl, err := template.ParseFiles(lp, fp, headerPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
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
		log.Error(err)
		//fmt.Println(err)
	}
	if i < 1024 {
		return false
	} else {
		return true
	}
}
//--------------------------------------------------------------------------------------
func loadAppConfig() {
	testConfig := 
		`<seelog>
		  <outputs formatid="common">
		    <filter levels="trace,debug,info,warn,error,critical">
		      <file path="`+ *logfile_Name +`"/>
		    </filter>
		  </outputs>
		  <formats>
		    <format id="common" format="%Date%t[%LEVEL]%t%FullPath:%Func:%Line%t%Msg%n" />
		  </formats>
		</seelog>`

	logger, err := log.LoggerFromConfigAsBytes([]byte(testConfig))
	//fmt.Println(testConfig)
	if err != nil {
		fmt.Println("seelog error!")
		fmt.Println(err)
	}
	
	loggerErr := log.ReplaceLogger(logger)
	
	if loggerErr != nil {
		fmt.Println(loggerErr)
	}
	log.Trace("Test message!")
	doLog()
	
}

func doLog() {
	for i:=0; i < 5; i++ {
		log.Tracef("%d", i)
	}
}

func doVersion(){
	if(strings.EqualFold("Yes", *version)){
		fmt.Print("This is currently Version ")
		fmt.Println(versionNo)
	}
}

func main() {
	defer log.Flush()
	loadAppConfig()	
	doVersion()
    log.Info("SERVER ONLINE!")
	//create server
	flag.Parse()
	//to handle different url
	http.HandleFunc("/time/", TimeServer)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/index.html/login", loginHandler)
	http.HandleFunc("/", home)
	http.HandleFunc("/index.html", home)
	http.HandleFunc("/logout/", logout)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssPath))))
	if !checkPort() {
		log.Critical("Error trying to connect to privledged port\n")
		os.Exit(404)
	}
	err := http.ListenAndServe(appendColon(*default_port), nil)
	if err != nil {
		log.Critical("Server's ListenAndServer critical error")
		os.Exit(404)		
	}	
}