package main

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/adamyi/CTFProxy/third_party/eddsa"
	"golang.org/x/crypto/bcrypt"
)

var _db *sql.DB

type Configuration struct {
	BackdoorPwd   string
	ListenAddress string
	DbType        string
	DbAddress     string
	VerifyKey     *ed25519.PublicKey
}

var _configuration = Configuration{}
var ctf_domain string

func initGaiaRsp(w http.ResponseWriter) {
	w.Header().Add("Server", "gaia")
}

func verifyPassword(username string, password string) bool {
	if username == "" {
		return false
	}
	if password == _configuration.BackdoorPwd && _configuration.BackdoorPwd != "" {
		return true
	}
	var storedPassword string
	err := _db.QueryRow("SELECT password FROM users WHERE ldap=?", username).Scan(&storedPassword)

	if err != nil {
		log.Println(err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func listenAndServe(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		initGaiaRsp(w)

		data := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Malformed Data", http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		if verifyPassword(data.Username, data.Password) {
			fmt.Fprintln(w, "👍")
			return
		}

		http.Error(w, "uhh", http.StatusForbidden)
	})

	mux.HandleFunc("/api/getgroups", func(w http.ResponseWriter, r *http.Request) {
		initGaiaRsp(w)

		rows, err := _db.Query("SELECT `group` FROM groups WHERE ldap=?", r.URL.Query().Get("ldap"))

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.Write([]byte("[]"))
			fmt.Println(err)
			return
		}

		result := make([]string, 0)
		for rows.Next() {
			grp := ""
			rows.Scan(&grp)
			result = append(result, grp)
		}

		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/api/addtogroup", func(w http.ResponseWriter, r *http.Request) {
		initGaiaRsp(w)
		if r.URL.Query().Get("group") == "" || r.URL.Query().Get("ldap") == "" {
			http.Error(w, "invalid request", http.StatusInternalServerError)
			return
		}
		name, err := getUsername(r.Header.Get("X-CTFProxy-JWT"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "I don't know what happened", http.StatusInternalServerError)
			return
		}
		name += "." + r.URL.Query().Get("group") + "@groups." + ctf_domain
		_db.Exec("INSERT INTO groups(ldap, `group`) VALUES(?, ?)", r.URL.Query().Get("ldap"), name)
		fmt.Fprintln(w, "👍")
	})

	mux.HandleFunc("/api/getusers", func(w http.ResponseWriter, r *http.Request) {
		initGaiaRsp(w)

		data := map[string]string{}

		rows, err := _db.Query("select ldap, affiliation from users WHERE hidden != 1")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "I don't know what happened", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			name := ""
			affiliation := ""
			rows.Scan(&name, &affiliation)
			data[name] = affiliation
		}

		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		initGaiaRsp(w)
		fmt.Fprintln(w, "👍")
	})

	return http.ListenAndServe(addr, mux)
}

func readConfig() {
	var publicKeyPath string
	flag.StringVar(&_configuration.BackdoorPwd, "backdoor_password", "", "use this password to login as arbitrary user (empty to disable)")
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&_configuration.DbType, "dbtype", "", "database type")
	flag.StringVar(&_configuration.DbAddress, "dbaddr", "", "database address")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.Parse()
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = eddsa.ParseEdPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
}

func main() {
	readConfig()
	var err error
	_db, err = sql.Open(_configuration.DbType, _configuration.DbAddress)
	if err != nil {
		panic(err)
	}
	defer _db.Close()

	log.Panic(listenAndServe(_configuration.ListenAddress))
}
