package jsondb

import "github.com/schollz/jsonstore"

type JsonStore struct {
	store *jsonstore.JSONStore
}

func (jsonStore *JsonStore) setupDB() {
	jsonStore.store = new(jsonstore.JSONStore)
	if err := jsonstore.Save(jsonStore.store, "holepunch.json.gz"); err != nil {
		panic(err)
	}
}

func (jsonStore *JsonStore) getSubdomainID(subdomainName string) (int, error) {
	var id int
	err := jsonStore.store.Get("subdomain:"+subdomainName, &id)
	if err != nil {
		panic(err)
	}
	return id, nil
}
func (jsonStore *JsonStore) refreshSubdomainIDs(subdomains []string, ids []int) error {

	return nil
}
func (jsonStore *JsonStore) addSubdomainID(subdomainName string, id int) error {
	err := jsonStore.store.Set("subdomain:"+subdomainName, &id)
	if err != nil {
		panic(err)
	}
	return nil
}
func (jsonStore *JsonStore) getTunnelID(subdomainName string) (int, error) {
	var id int
	err := jsonStore.store.Get("tunnel:"+subdomainName, &id)
	if err != nil {
		panic(err)
	}
	return id, nil
}
func (jsonStore *JsonStore) refreshTunnelIDs(subdomains []string, ids []int) error {

	return nil
}
func (jsonStore *JsonStore) addTunnelID(subdomainName string, id int) error {
	err := jsonStore.store.Set("tunnel:"+subdomainName, &id)
	if err != nil {
		panic(err)
	}
	return nil
}
