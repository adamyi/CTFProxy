package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
)

type Flag struct {
	Id          string
	DisplayName string
	Category    string
	Points      int
	Type        string
	Flag        string
	Prefix      string
	Owner       string
}

type VerificationResponse struct {
	Success int
	Message string
	Flag    Flag
}

type Configuration struct {
	ListenAddress string
	FlagKey       string
	VerifyKey     *ed25519.PublicKey
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var _configuration = Configuration{}

var _flags map[string]Flag
var d2f map[string]Flag

func FlagSplit(r rune) bool {
	return r == '{' || r == '.' || r == '}'
}

func UserSplit(r rune) bool {
	return r == '+' || r == '@'
}

func readConfig() {
	var publicKeyPath, flagConfigPath string
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&_configuration.FlagKey, "flag_key", "", "flag signature key")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.StringVar(&flagConfigPath, "flag_config", "flags.json", "path to flags configuration json")
	flag.Parse()
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = eddsa.ParseEdPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
	var flags []Flag
	file, _ := os.Open(flagConfigPath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&flags)
	if err != nil {
		panic(err)
	}
	_flags = make(map[string]Flag)
	d2f = make(map[string]Flag)
	for _, flag := range flags {
		_, exists := _flags[flag.Id]
		if exists {
			panic("duplicate flag id: " + flag.Id)
		}
		if flag.Type != "fixed" && flag.Type != "dynamic" {
			panic("wrong flag type for flag id: " + flag.Id)
		}
		fk := flag.Prefix + flag.Flag
		_, exists = d2f[fk]
		if exists {
			panic("duplicate flag for flag id: " + flag.Id)
		}
		_flags[flag.Id] = flag
		d2f[fk] = flag
		log.Println(fk)
	}
}

func initFGRsp(rsp http.ResponseWriter) {
	rsp.Header().Add("Server", "flaganizer")
	rsp.Header().Add("Content-Type", "text/plain")
}

func genFlag(prefix, username, flag, id string) string {
	key := []byte(_configuration.FlagKey)
	h := hmac.New(md5.New, key)
	h.Write([]byte(username + "_" + id))
	return prefix + "{" + flag + "." + base64.StdEncoding.EncodeToString([]byte(username)) + "." + base64.StdEncoding.EncodeToString(h.Sum(nil)) + "}"
}

func (f Flag) GetFlag(username string) string {
	userparts := strings.FieldsFunc(username, UserSplit)
	if userparts[0] != f.Owner {
		return "NOT_A_FLAG{NO_PERMISSION_TO_GET_FLAG_ASK_COURSE_STAFF}"
	}
	if f.Type == "fixed" {
		return f.Prefix + "{" + f.Flag + "}"
	}
	if f.Type == "dynamic" {
		if len(userparts) != 3 {
			return "NOT_A_FLAG{NO_SUBACC_FOR_DYNAMIC_FLAG_ASK_COURSE_STAFF}"
		}
		return genFlag(f.Prefix, userparts[1], f.Flag, f.Id)
	}
	return "NOT_A_FLAG{WRONG_FLAG_TYPE_ASK_COURSE_STAFF}"
}

func GenerateFlag(rsp http.ResponseWriter, req *http.Request) {
	initFGRsp(rsp)

	tknStr := req.Header.Get("X-CTFProxy-JWT")

	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{eddsa.SigningMethodEdDSA.Alg()}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			rsp.Write([]byte("NOT_A_FLAG{JWT_SIGNATURE_INVALID_ASK_COURSE_STAFF}"))
			return
		}
		rsp.Write([]byte("NOT_A_FLAG{JWT_ERROR_ASK_COURSE_STAFF}"))
		return
	}

	if !tkn.Valid {
		rsp.Write([]byte("NOT_A_FLAG{JWT_INVALID_ASK_COURSE_STAFF}"))
		return
	}

	username := claims.Username

	flagObj, ok := _flags[req.URL.Query().Get("id")]
	if !ok {
		rsp.Write([]byte("NOT_A_FLAG{WRONG_FLAG_ID_ASK_COURSE_STAFF}"))
		return
	}

	rsp.Write([]byte(flagObj.GetFlag(username)))
}

func VerifyFlag(rsp http.ResponseWriter, req *http.Request) {
	initFGRsp(rsp)
	rsp.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rsp)

	var r = VerificationResponse{}
	r.Success = 0

	tknStr := req.Header.Get("X-CTFProxy-JWT")

	claims := &Claims{}

	p := jwt.Parser{ValidMethods: []string{eddsa.SigningMethodEdDSA.Alg()}}
	tkn, err := p.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return _configuration.VerifyKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			r.Message = "jwt signature invalid"
			enc.Encode(r)
			return
		}
		r.Message = "jwt error"
		enc.Encode(r)
		return
	}

	if !tkn.Valid {
		r.Message = "jwt invalid"
		enc.Encode(r)
		return
	}

	username := claims.Username
	userparts := strings.FieldsFunc(username, UserSplit)

	if userparts[0] != "ctfd" {
		r.Message = "Calling service not allowlisted"
		enc.Encode(r)
		return
	}

	flagStr := strings.TrimSpace(req.FormValue("flag"))
	log.Println(flagStr)
	flagParts := strings.FieldsFunc(flagStr, FlagSplit)
	log.Println(flagParts)

	if len(flagParts) < 2 {
		log.Println("len(flagParts) < 2")
		r.Message = "Incorrect flag"
		enc.Encode(r)
		return
	}
	fk := flagParts[0] + flagParts[1]
	flagObj, ok := d2f[fk]
	if !ok {
		log.Println("Didn't find flagObj " + fk)
		r.Message = "Incorrect flag"
		enc.Encode(r)
		return
	}
	if flagObj.Type == "fixed" && len(flagParts) == 2 {
		log.Println("static success")
		r.Success = 1
		r.Flag = flagObj
		enc.Encode(r)
		return
	}
	if flagObj.Type == "dynamic" && len(flagParts) == 4 {
		fusername, err := base64.StdEncoding.DecodeString(flagParts[2])
		log.Printf("Dynamic flag - %v (%d) - %v (%d)", string(fusername), len(string(fusername)), userparts[1], len(userparts[1]))
		if err == nil && string(fusername) == userparts[1] {
			rf := genFlag(flagObj.Prefix, string(fusername), flagObj.Flag, flagObj.Id)
			if rf == flagStr {
				log.Println("Dynamic success")
				r.Success = 1
				r.Flag = flagObj
				enc.Encode(r)
				return
			}
		}
	}
	log.Printf("fallback incorrect")
	r.Message = "Incorrect flag"
	enc.Encode(r)
	return

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	initFGRsp(w)
	w.Write([]byte("ok"))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	readConfig()
	http.HandleFunc("/generate", GenerateFlag)
	http.HandleFunc("/verify", VerifyFlag)
	http.HandleFunc("/healthz", HealthCheckHandler)
	err := http.ListenAndServe(_configuration.ListenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
