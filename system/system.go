package system

import (
	"path/filepath"
	"strings"

	"flag"
	"fmt"

	"io/ioutil"
	"os"

	"log"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	//	"option.bzza.com/models/redis"
)

var (
	yamlConfigFile = flag.String("file", "./config.yaml", "(Simple) YAML file to read") //配置文件

	Conf Config //config.yaml配置文件
)

type Config struct { //config.yaml配置文件
	Web struct {
		IsMasterServer bool   `yaml:"isMasterServer"`
		Port           string `yaml:"port"`
		JsConfigPath   string `yaml:"jsConfigPath"`
	}

	Redis struct {
		Server   string `yaml:"server"`
		Password string `yaml:"password"`
	}
	Kdb struct {
		Server1 string   `yaml:"server1"`
		Server2 []string `yaml:"server2"`
	}
	Db struct {
		DriverName     string `yaml:"driverName"`
		DataSourceName string `yaml:"dataSourceName"`
	}
	Certificate struct {
		Pem string `yaml:"pem"`
		Key string `yaml:"key"`
	}
	Datasource struct {
		Aliyun struct {
			HttpORhttps string `yaml:"httpORhttps"`
			AppKey      string `yaml:"appKey"`
			AppSecret   string `yaml:"appSecret"`
		}
		A7a2 struct {
			WsORwss string `yaml:"wsORwss"`
		}
	}
}

var ChanOptionSettingsArrayInit = make(chan bool, 10)

func FileChangeWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					switch event.Name {
					case "config.yaml":
						YamlInit()
					case "optionSettingsArray.json":
						ChanOptionSettingsArrayInit <- true
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func YamlInit() {
	flag.Parse()
	b := ReadFile(*yamlConfigFile)
	err := yaml.Unmarshal([]byte(*b), &Conf)
	if err != nil {
		log.Fatalf("readfile(%q): %s", *yamlConfigFile, err)
	}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func OpenFile(path string) (fi *os.File, err error) {
	if CheckFileIsExist(path) { //如果文件存在
		fi, err = os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fi, err = os.Create(path) //创建文件
		if err != nil {
			fmt.Println(err)
		}
	}
	return fi, err
}

func WriteFile(path string, bytes []byte) error {
	err := ioutil.WriteFile(path, bytes, 0666) //写入文件(字节数组)
	return err
}

func ReadFile(path string) *[]byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return &fd
}
