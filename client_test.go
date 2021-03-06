package ezk

import (
	"fmt"
	cv "github.com/glycerine/goconvey/convey"
	zook "github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

func Test001ClientRetryGetsDefault(t *testing.T) {
	cv.Convey("Given that we don't configure a Retry function, we should get the DefaultRetry function in our ClientConfig", t, func() {

		cli := NewClient(ClientConfig{})
		cv.So(cli.Cfg.Retry, cv.ShouldEqual, DefaultRetry)

	})
}

// In the example code, the goal is to set
// the value newURL into the node /chroot/service-name/config/server-url-list
func ExampleClient() {
	newURL := "http://my-new-url.org:343/hello/enhanced-zookeeper-client"

	base := "/chroot/"
	path := "service-name/config/server-url-list"
	zkCfg := ClientConfig{
		Servers:        []string{"127.0.0.1:2181"},
		Acl:            zook.WorldACL(zook.PermAll),
		Chroot:         base,
		SessionTimeout: 10 * time.Second,
	}
	zk := NewClient(zkCfg)
	err := zk.Connect()
	if err != nil {
		panic(err)
	}

	defer zk.Close()

	err = zk.CreateDir(path, nil)
	if err != nil {
		panic(err)
	}

	err = zk.DeleteNode(path) // delete any old value
	if err != nil {
		panic(err)
	}

	err = zk.CreateNode(path)
	if err != nil {
		panic(err)
	}

	_, err = zk.Set(path, []byte(newURL), -1)
	if err != nil {
		panic(err)
	}
}

func Test002RemoveChroot(t *testing.T) {
	cv.Convey("Given an absolute Chrooted path, the RemoveChoot() function should return the relative path without the Chroot prefix", t, func() {

		// map from input to expected output
		m := map[string]string{
			"/mybase/myservice/config":     "myservice/config",
			"/myroot/alist":                "alist",
			"/hello/":                      "",
			"/poorlyFormedChrootPrefix":    "",
			"/properlyFormedChrootPrefix/": "",
			"relative/path/unchanged":      "relative/path/unchanged",
			"/":          "",
			"//":         "",
			"abc":        "abc",
			"a/b/c/d/e/": "a/b/c/d/e/",
		}

		for k, v := range m {
			fmt.Printf("\n checking '%s' -> '%s'\n", k, v)
			cv.So(RemoveChroot(k), cv.ShouldEqual, v)
		}
	})
}
