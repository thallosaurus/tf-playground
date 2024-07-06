package tlsg108g

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	v8 "rogchap.com/v8go"
)

// Used to interact with the CGI Scripts
func DataRequest(endpoint DataRequestEndpoint, v url.Values) []byte {
	uri := fmt.Sprintf("http://%s/%s", host, endpoint)
	res, err := http.Post(uri, "application/x-www-urlencoded", strings.NewReader(v.Encode()))

	if nil != err {
		log.Fatal(err)
	}

	//log.Println("status", res.StatusCode)
	body, io_err := io.ReadAll(res.Body)

	if nil != io_err {
		log.Fatal(io_err)
	}

	//log.Println(string(body[:]))
	return body
}

func DataRequestParse(jsname string, endpoint DataRequestEndpoint, v url.Values) []byte {
	body := DataRequest(endpoint, v)

	return parse(jsname, bytes.NewReader(body))
}

type RequestEndpoint string
type DataRequestEndpoint string

const (
	VLAN_8021Q_RPM      RequestEndpoint = "Vlan8021QRpm.htm"
	VLAN_8021Q_PVID_RPM RequestEndpoint = "Vlan8021QPvidRpm.htm"
	INDEX               RequestEndpoint = ""
	LOGOUT              RequestEndpoint = "Logout.htm"

	LOGON    DataRequestEndpoint = "logon.cgi"
	QVlanSet RequestEndpoint     = "qvlanSet.cgi"
)

var host = ""

func SetHost(h string) {
	host = h
}

func RequestNoParse(endpoint RequestEndpoint) {
	uri := fmt.Sprintf("http://%s/%s", host, endpoint)

	// make http request
	resp, err := http.Get(uri)
	if nil != err {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	log.Println("reqnoparse", resp.StatusCode)
}

func RequestParam(params url.Values, jsname string, endpoint string) []byte {
	ep := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	return Request(jsname, ep)
}

// Used to request files from the frontend
func Request(jsname string, endpoint string) []byte {
	uri := fmt.Sprintf("http://%s/%s", host, endpoint)

	// make http request
	resp, err := http.Get(uri)
	if nil != err {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// read body response into buffer
	body, err := io.ReadAll(resp.Body)

	if nil != err {
		log.Fatal(err)
	}

	return parse(jsname, bytes.NewReader(body))
}

// Takes a bytes Reader that should read HTML and extracts the JavaScript Value of the given Variable Name as JSON bytes
func parse(jsname string, r *bytes.Reader) []byte {
	z, err := html.Parse(r)

	if nil != err {
		log.Fatal(err)
	}

	g := extract(jsname, z)
	g_b, err := g.MarshalJSON()
	if nil != err {
		log.Fatalln(err)
		return nil
	}

	return g_b
}

func extract(jsname string, n *html.Node) *v8.Value {
	if n.Type == html.ElementNode && n.Data == "script" {
		if n.FirstChild != nil && n.FirstChild.Data != "" {
			ctx := v8.NewContext()
			//log.Println("n.FirstChild.Data: ", n.FirstChild.Data)

			ctx.RunScript(n.FirstChild.Data, "parser.js")

			// try to return the value
			v, err := ctx.RunScript(jsname, "parser.js")

			if nil != err {
				log.Println(n.FirstChild.Data)
				log.Println(err)
				return nil
			}
			//ctx.Close()
			return v
		}
	} else {
		var v *v8.Value
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			v = extract(jsname, c)
			if nil != v {
				return v
			}
		}
		return v
	}
	return nil
}
