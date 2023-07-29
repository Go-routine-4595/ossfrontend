package filedb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Filedb struct {
	tenantMap map[string]interface{}
}

func NewFileDB(keyfile string) *Filedb {
	var (
		err error
		f   Filedb
	)

	dbk, _ := readDB(keyfile)
	err = json.Unmarshal(dbk, &f.tenantMap)

	if err != nil {
		fmt.Println("No Public Key file found or bad Public Key file")
	}
	return &f
}

func (f *Filedb) UserSelector(t string) string {
	return f.tenantMap[t].(map[string]interface{})["user"].(string)
}

func (f *Filedb) RoleSelector(t string) string {
	return f.tenantMap[t].(map[string]interface{})["role"].(string)
}

func (f *Filedb) KeySelector(t string) []byte {
	tmp := f.tenantMap[t].(map[string]interface{})
	b, err := base64.StdEncoding.DecodeString(tmp["pubkey"].(string))
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func readDB(s string) ([]byte, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)

	return b, err
}
