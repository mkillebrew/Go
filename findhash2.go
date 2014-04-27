/*
	unfuck loop, exclude 1918 and multicast - mostly unfucked
	don't build IP string in a stupid manner, use net lib as host and geoip will require it anyway - still pretty stupid
	don't convert entire bytearray, just what we need, to a fast compare instead of HasPrefix
	loop in GeoIP, and name resolution - done, don't need resolution
	filter country in code and not in grep - done, case insensitive
	input findhash2 <hash> <country>
	output ip - hash - city country
*/

package main
import (
	"fmt"
	"strconv"
	"strings"
	"encoding/hex"
	"crypto/md5"
    "github.com/oschwald/geoip2-golang"
	"net"
	"io"
	"os"
)

func main(){
	var findhash string = os.Args[1]
	var min, max int = 0, 255
    hasher := md5.New()
    db, err := geoip2.Open("GeoLite2-City.mmdb")

	if err != nil {
		panic(err)
 	}    

	if len(os.Args) == 3 {
		findcountry := os.Args[2]
		findhash=strings.ToLower(findhash)
		fmt.Printf("Country %v specified\n", findcountry)
	} else {
		fmt.Println("No country specified")
	}

	for a:=1; a<=223; a++{
		if (a==10 || a==127) {
			a++
		}
		for b:=min; b<=max; b++{
			for c:=min; c<=max; c++{
				for d:=min; d<=max; d++{
					ip := []string{strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d)}
					ipstring := strings.Join(ip, ".")
    				io.WriteString(hasher, ipstring)
    				hash:=(hasher.Sum(nil))
					if strings.HasPrefix(hex.EncodeToString(hash), findhash) {
						gip:=net.ParseIP(ipstring)
						record, err := db.City(gip)
						if err != nil {
							fmt.Println("poop")
						}
						if len(os.Args) == 3 {
							findcountry := os.Args[2]
							findhash=strings.ToLower(findhash)
							if strings.Contains(strings.ToLower(record.Country.Names["en"]), strings.ToLower(findcountry)) {
								fmt.Printf("%s - %x - %v %v\n", ipstring, hash, record.City.Names["en"], record.Country.Names["en"])
							}
						} else {
							fmt.Printf("%s - %x - %v %v\n", ipstring, hash, record.City.Names["en"], record.Country.Names["en"])
						}
					}
					hasher.Reset()
				}
			}
		}
	}

	db.Close()
}
