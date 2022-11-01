package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"main/block"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/julienschmidt/httprouter"
)

func init() {

}
func main() {

	r := httprouter.New()

	r.GET("/", anasayfa)
	r.POST("/deneme", deneme)

	http.ListenAndServe(":8081", r)
}
func cmdOut() {

	cmd := exec.Command("ipfs", "daemon")
	cmd.CombinedOutput()
}

func anasayfa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	cmdOut()
	view, _ := template.ParseFiles("index.html")
	view.Execute(w, nil)
}

func deneme(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	r.ParseMultipartForm(10 << 20)
	a := r.Form["crypto"]
	fmt.Println(a)
	file, header, _ := r.FormFile("file")
	f, _ := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)

	sh := shell.NewShell("localhost:5001")

	cid, _ := sh.AddDir(f.Name())

	bytes := make([]byte, 32) //AES-256 için rastgele bir 32 bayt anahtar oluşturun.
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	fmt.Printf("\n")
	key := hex.EncodeToString(bytes) //anahtarı bayt cinsinden kodlayın ve gizli olarak saklayın, bir kasaya koyun
	fmt.Printf("key(anahtar) : %s\n", key)

	encrypted := block.Encrypt(cid, key)
	fmt.Printf("encrypted(şifreli) : %s\n", encrypted)

	decrypted := block.Decrypt(encrypted, key)
	fmt.Printf("decrypted(şifre çözüm) : %s\n", decrypted)

	fmt.Printf("https://ipfs.io/ipfs/%s", decrypted)
	Key := "furkan"
	fmt.Printf("https://ipfs.io/ipfs/%s", decrypted)
	isim := "https://ipfs.io/ipfs/" + decrypted + " KEY:   " + Key

	//burada şablon oluşturuyoruz
	şablon, _ := template.ParseFiles("index.html")

	//Burada da şablonu çalıştırmasını ve isim
	//değişkenini kullanmasını istiyoruz.

	şablon.Execute(w, isim)
}
