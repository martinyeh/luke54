package main 

import ( 
	"flag"
        "fmt" 
	"log"
        "net/http" 
        "io/ioutil" 
	"bytes"
	"encoding/json"
	"strings"
	//transform "code.google.com/p/go-html-transform/html/transform"
	//h5 "code.google.com/p/go-html-transform/h5"
	html "code.google.com/p/go.net/html"
) 

var (
	port      = flag.Int("p", 8080, "Luke54 Service Port")
	url	  = flag.String("url", "http://www.luke54.org/category", "Luke54 url")
)


type ForumHandler struct {

}

type ForumData struct {
    title string
    catg string
    page string
}

type catgHandler struct {
    	
}

type Message struct {
    title string
    iurl string
    curl string
}

func (hf *ForumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fd := ForumData{}
	fd.title = "福音真理"
	fd.catg = "9"
	fd.page = "2"

	s := ""
	b, _ := json.Marshal(fd)
	s += fmt.Sprintf("%s,", b)	

	fd.title = "福音小品"
	fd.catg = "10"
	fd.page = "3"	

	b, _ = json.Marshal(fd)
	s += fmt.Sprintf("%s", b)

	fmt.Fprintf(w, "[%s]", s)
}

func (hf *catgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

	r.ParseForm()                     // Parses the request body
    	catg := r.Form.Get("catg") 
     	page := r.Form.Get("page")

	url := fmt.Sprintf(*url+"/%s/%s.html", catg, page)
	fmt.Println(url)	

	res, err := http.Get(url) 
        if err != nil { 
                fmt.Println("http.Get", err) 
                return 
        } 
        defer res.Body.Close() 
         
        body, err := ioutil.ReadAll(res.Body) 
        if err != nil { 
                fmt.Println("ioutil.ReadAll", err) 
                return 
        } 

	reader := bytes.NewBufferString(string(body))

	doc, _ := html.Parse(reader)

	m := Message{}
	//var msgmap map[string]Message
	msgmap := make(map[string]Message)

	var f func(*html.Node)
	f = func(n *html.Node) {
		
		if n.Type == html.ElementNode && n.Data == "meta" {

			for _, attr := range n.Attr{
				fmt.Println("attr.Val:" , attr.Val);
				if strings.Contains( attr.Val, "og:title"){
					m.title = n.Attr[1].Val	
					m.iurl = ""				
				}else if strings.Contains( attr.Val, "og:iurl"){
					m.iurl = n.Attr[1].Val	
					msgmap[m.title]	= m			
				}
			}
	
			/*b, err := json.Marshal(m)
			if err==nil{
				fmt.Println("b:" , string(b));
			}*/
		} else if n.Type == html.ElementNode && n.Data == "a"{
			for _, attr := range n.Attr{
				fmt.Println("attr.Val:" , attr.Val);
				if _, ok := msgmap[attr.Val]; ok {
					 m := msgmap[attr.Val]
					 m.curl = n.Attr[0].Val	
					 msgmap[attr.Val] = m										
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	fmt.Println("msgmap:" , msgmap);
	var s string
	for _, value := range msgmap {
	    //fmt.Println("Key:", key, "Value:", value)
	    b, _ := json.Marshal(value)
	    s += fmt.Sprintf("%s,", b)		
	}
	fmt.Fprintf(w, "[%s]", s[:len(s)-1])
}


func main() { 
	flag.Parse()

        
        //lenp := len(body) 
        //if maxp := 60; lenp > maxp { 
          //      lenp = maxp 
        //} 
        //fmt.Println(len(body), string(body[:lenp])) 
	//fmt.Println(string(body))


	http.Handle("/category", &catgHandler{})
	http.Handle("/forum", &ForumHandler{})

	fmt.Println("server start...")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		

	if err != nil {
		log.Fatal("Listen and serve error: ", err)
	}
} 




