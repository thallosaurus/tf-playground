package tlsg108g

import (
	"encoding/json"
	"log"
	"net/url"
)

func CheckLogin() bool {
	login_status_res := Request("logonInfo", string(QVlanSet))

	var res []int64
	json_err := json.Unmarshal(login_status_res, &res)

	if nil != json_err {
		log.Fatal(json_err)
	}

	log.Println("check login ", res)

	// 1 means no session
	return res[0] == 1
}

func Login(username string, password string) bool {
	v := url.Values{}
	v.Add("username", username)
	v.Add("password", password)
	v.Add("cpassword", "")
	v.Add("logon", "Login")

	linfo := DataRequestParse("logonInfo", LOGON, v)

	var buf []int64
	err := json.Unmarshal(linfo, &buf)

	if nil != err {
		log.Fatal(err)
	}
	return buf[0] == 0
	//}
}

func Logout() {
	RequestNoParse(LOGOUT)

	//log.Println(res)
}
