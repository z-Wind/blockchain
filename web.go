package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"./core"
	"github.com/pkg/errors"
)

// Instantiate the Blockchain
var blockchain = core.NewBlockchain()

func main() {
	host := flag.String("host", "127.0.0.1", "port to listen on")
	port := flag.String("port", "6060", "port number")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/chain/", fullChain)
	mux.HandleFunc("/transactions/new", newTransactions)
	mux.HandleFunc("/mine/", mine)
	mux.HandleFunc("/nodes/register", registerNodes)
	mux.HandleFunc("/nodes/resolve", consensus)

	addr := *host + ":" + *port
	log.Println("Listening on ", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func fullChain(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain)
}

func newTransactions(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	form := make(map[string]string)
	for _, k := range []string{"sender", "recipient", "amount"} {
		v, ok := req.Form[k]
		if !ok {
			fmt.Fprintf(w, "Missing %s\n", k)
			continue
		}
		form[k] = v[0]
	}

	amount, err := strconv.ParseFloat(form["amount"], 64)
	if err != nil {
		fmt.Fprintf(w, "amount:%s is not a number", form["amount"])
		return
	}

	blockchain.NewTransaction(form["sender"], form["recipient"], amount)
	fmt.Fprintf(w, "CurrentTransactions: %s", blockchain.CurrentTransactions)
}

func mine(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	blockchain.NewBlock()
	fmt.Fprintf(w, "Chain: %s", blockchain.Chain)
}

func registerNodes(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(w, "error:%s", err)
		return
	}

	addrs, ok := req.Form["addr"]
	if !ok {
		fmt.Fprintln(w, "Missing addr")
		return
	}
	for _, addr := range addrs {
		blockchain.RegisterNode(addr)
	}

	for node := range blockchain.Nodes {
		fmt.Fprintf(w, "%s,\n", node)
	}
}

func consensus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newBlockchain := &core.Blockchain{}
	for k := range blockchain.Nodes {
		s := k + "/chain/"
		log.Printf("update by %s\n", s)
		err := getJSON(s, newBlockchain)
		if err != nil {
			log.Printf("getJSON Error: %s\n", err)
			continue
		}

		changed, err := blockchain.ResolveConflicts(newBlockchain)
		if err != nil {
			log.Printf("blockchain.ResolveConflicts Error: %s\n", err)
			continue
		}

		if changed {
			fmt.Fprintln(w, "updated")
		}

	}

	fmt.Fprintf(w, "chain: %s", blockchain.Chain)
}

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "http.Get(%s)", url)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
