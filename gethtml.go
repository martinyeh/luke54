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
	html "code.google.com/p/go.net/html"
) 

var (
	port      = flag.Int("p", 8080, "Luke54 Service Port")
	url	  = flag.String("url", "http://www.luke54.org/category", "Luke54 url")
	vurl      = flag.String("vurl", "http://www.luke54.org/view", "Luke54 url")
)


type ForumHandler struct {

}

type ForumData struct {
    Title string    `json:"title"` 
    Catg string	   `json:"catg"`
    Page string    `json:"page"`
    Fid string `json:"fid"`
}

type CatgHandler struct {
    	
}

type Message struct {
    Title string `json:"title"`
    Iurl string `json:"iurl"`
    Curl string `json:"curl"`
    Fid string `json:"fid"`
}

type ViewHandler struct {

}


type ViewData struct {
    Html string    `json:"html"`
    Vid string `json:"vid"`
}


func (hf *ForumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fd := ForumData{}
	fd.Title = "福音真理"
	fd.Catg = "9"
	fd.Page = "2"
	fd.Fid = "92"

	s := ""
	b, _ := json.Marshal(fd)
	s += fmt.Sprintf("%s,", b)	

	fd.Title = "福音小品"
	fd.Catg = "10"
	fd.Page = "3"
	fd.Fid = "103"	

	b, _ = json.Marshal(fd)
	s += fmt.Sprintf("%s", b)

	fmt.Fprintf(w, "[%s]", s)
}

func GetHtml(url string) *bytes.Buffer {
        res, err := http.Get(url)
        if err != nil {
                fmt.Println("http.Get", err)
                return nil
        }
        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
                fmt.Println("ioutil.ReadAll", err)
                return nil
        }

        reader := bytes.NewBufferString(string(body))
        
        return reader 
}


func (hf *CatgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

	r.ParseForm()                     // Parses the request body
    	catg := r.Form.Get("catg") 
     	page := r.Form.Get("page")

	url := fmt.Sprintf(*url+"/%s/%s.html", catg, page)
	fmt.Println(url)	

	doc, _ := html.Parse(GetHtml(url))

	m := Message{}
	m.Fid = catg+page
	//var msgmap map[string]Message
	msgmap := make(map[string]Message)

	var f func(*html.Node)
	f = func(n *html.Node) {
		
		if n.Type == html.ElementNode && n.Data == "meta" {

			for _, attr := range n.Attr{
				fmt.Println("attr.Val:" , attr.Val);
				if strings.Contains( attr.Val, "og:title"){
					m.Title = n.Attr[1].Val	
					m.Iurl = ""				
				}else if strings.Contains( attr.Val, "og:image"){
					m.Iurl = n.Attr[1].Val	
					msgmap[m.Title]	= m			
				}
			}
	
		} else if n.Type == html.ElementNode && n.Data == "a"{
			for _, attr := range n.Attr{
				fmt.Println("attr.Val:" , attr.Val);
				if _, ok := msgmap[attr.Val]; ok {
					 m := msgmap[attr.Val]
					 m.Curl = n.Attr[0].Val	
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
	for key, value := range msgmap {
	    fmt.Println("Key:", key, "Value:", value)
	    b, _ := json.Marshal(value)
	    s += fmt.Sprintf("%s,", b)		
	}
	fmt.Fprintf(w, "[%s]", s[:len(s)-1])
}

func (hf *ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

        r.ParseForm()                     // Parses the request body
        catg := r.Form.Get("catg")
        page := r.Form.Get("page")

	url := fmt.Sprintf(*vurl+"/%s/%s.html", catg, page)
        fmt.Println(url)
	doc, _ := html.Parse(GetHtml(url))

	vd := ViewData{}
        vd.Vid = catg+page

	content := "<html><head><title></title></head><body>"

	var f func(*html.Node)
        f = func(n *html.Node) {

                if n.Type == html.ElementNode && n.Data == "div" {

                        for _, attr := range n.Attr{
                                fmt.Println("attr.Val:" , attr.Val);
                                if strings.Contains( attr.Val, "itemFullText"){
					content+= n.Data
                                }
                        }

                } 

		for c := n.FirstChild; c != nil; c = c.NextSibling {
                        f(c)
                }
	}
	f(doc)

	content+= "</body></html>"
	vd.Html = content

	b, _ := json.Marshal(vd)
        s := fmt.Sprintf("%s", b)
        fmt.Fprintf(w, "[%s]", s)
}


func main() { 
	flag.Parse()


	http.Handle("/category", &CatgHandler{})
	http.Handle("/forum", &ForumHandler{})
	http.Handle("/view", &ViewHandler{})

	fmt.Println("server start...")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		

	if err != nil {
		log.Fatal("Listen and serve error: ", err)
	}
} 




