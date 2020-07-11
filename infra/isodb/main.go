package main

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Configuration struct {
	ListenAddress string
	DbType        string
	DbAddress     string
	DbHost        string
	DbUser        string
	DbName        string
	VerifyKey     *rsa.PublicKey
}

type InitRequest struct {
	Version string
	SQL     string
}

func randPassword() string {
	return randStringWithCharset(16, passwordCharset)
}

func randStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

var _configuration = Configuration{}

var dblock sync.Mutex

func initIsoDbRsp(w http.ResponseWriter) {
	w.Header().Add("Server", "isodb")
	w.Header().Add("Content-Type", "application/json")
}

func returnError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err.Error())
}

func initDbConn(name, main string) (*sql.DB, error) {
	dblock.Lock()
	defer dblock.Unlock()
	var version, cred string
	err := _db.QueryRow("SELECT version FROM metadb WHERE service=?", main).Scan(&version)
	if err != nil {
		return nil, err
	}
	name = "isodb_" + name + "_" + version
	err = _db.QueryRow("SELECT credentials FROM instancedb WHERE service=?", name).Scan(&cred)
	switch {
	case err == sql.ErrNoRows:
		var s string
		err = _db.QueryRow("SELECT `sql` FROM metadb WHERE service=?", main).Scan(&s)
		if err != nil {
			return nil, err
		}
		un := randPassword()
		pwd := randPassword()
		cred = un + ":" + pwd
		_, err = _db.Exec("CREATE DATABASE " + name)
		if err != nil {
			return nil, err
		}
		var cdb *sql.DB
		cdb, err = sql.Open(_configuration.DbType, _configuration.DbUser+"@"+_configuration.DbHost+"/"+name+"?multiStatements=true")
		if err != nil {
			return nil, err
		}
		defer cdb.Close()
		_, err = cdb.Exec(s)
		if err != nil {
			return nil, err
		}
		_, err = _db.Exec("CREATE USER '" + un + "'@'%' IDENTIFIED BY '" + pwd + "'")
		if err != nil {
			return nil, err
		}
		_, err = _db.Exec("GRANT SELECT, UPDATE, INSERT, DELETE ON " + name + ".* TO '" + un + "'@'%'")
		if err != nil {
			return nil, err
		}
		_, err = _db.Exec("FLUSH PRIVILEGES")
		if err != nil {
			return nil, err
		}
		_, err = _db.Exec("INSERT INTO instancedb(service, credentials) VALUES(?, ?)", name, cred)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}
	return sql.Open(_configuration.DbType, cred+"@"+_configuration.DbHost+"/"+name)
}

func listenAndServe(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/init", func(w http.ResponseWriter, r *http.Request) {
		initIsoDbRsp(w)
		dblock.Lock()
		defer dblock.Unlock()
		_, mainservice, err := getUsername(r.Header.Get("X-CTFProxy-JWT"))
		if err != nil {
			returnError(w, err)
			return
		}
		ver := r.URL.Query().Get("version")
		s, err := ioutil.ReadAll(r.Body)
		if err != nil {
			returnError(w, err)
			return
		}
		var version string
		err = _db.QueryRow("SELECT version FROM metadb WHERE service=?", mainservice).Scan(&version)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("%s %s init", mainservice, ver)
			_, err = _db.Exec("INSERT INTO metadb (version, `sql`, service) VALUES(?,?,?)", ver, s, mainservice)
			if err != nil {
				returnError(w, err)
				return
			}
			return
		case err != nil:
			log.Printf("%s %s init err", mainservice, ver)
			returnError(w, err)
			return
		}
		if version == ver {
			log.Printf("%s %s already init", mainservice, ver)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Printf("%s %s update", mainservice, ver)
		_, err = _db.Exec("UPDATE metadb SET version=?, `sql`=? WHERE service=?", ver, s, mainservice)
		if err != nil {
			returnError(w, err)
			return
		}
	})

	mux.HandleFunc("/api/sql", func(w http.ResponseWriter, r *http.Request) {
		initIsoDbRsp(w)
		service, mainservice, err := getUsername(r.Header.Get("X-CTFProxy-JWT"))
		db, err := initDbConn(service, mainservice)
		if err != nil {
			returnError(w, err)
			return
		}
		defer db.Close()

		var params []interface{}
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			returnError(w, err)
			return
		}
		rows, err := db.Query(params[0].(string), params[1:]...)

		if err != nil {
			returnError(w, err)
			return
		}
		fmt.Fprintf(w, "[%s]", strings.Join(jsonify.Jsonify(rows), ""))
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		initIsoDbRsp(w)
		fmt.Fprintln(w, "👍")
	})

	return http.ListenAndServe(addr, mux)
}

func readConfig() {
	var publicKeyPath string
	flag.StringVar(&_configuration.ListenAddress, "listen", "0.0.0.0:80", "http listen address")
	flag.StringVar(&_configuration.DbType, "dbtype", "", "database type")
	flag.StringVar(&_configuration.DbAddress, "dbaddr", "", "database address")
	flag.StringVar(&publicKeyPath, "jwt_public_key", "", "Path to JWT public key")
	flag.Parse()
	JwtPubKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	_configuration.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(JwtPubKey)
	if err != nil {
		panic(err)
	}
	_configuration.DbHost = strings.Split(strings.Split(_configuration.DbAddress, "@")[1], "/")[0]
	_configuration.DbUser = strings.Split(_configuration.DbAddress, "@")[0]
	_configuration.DbName = strings.Split(_configuration.DbAddress, "/")[1]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	readConfig()
	var err error
	_db, err = sql.Open(_configuration.DbType, _configuration.DbAddress)
	if err != nil {
		panic(err)
	}
	defer _db.Close()

	log.Panic(listenAndServe(_configuration.ListenAddress))
}
