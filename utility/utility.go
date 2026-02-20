package utility

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jwalton/go-supportscolor"
)

type TsHttp struct {
	HttpPort                          string `json:"http_port"`
	MaxAllowedSimultaneousConnections int    `json:"max_allowed_simultaneous_connections"`
}

type TsLog struct {
	LogFilePath   string `json:"logs_file_path"`
	ErrorFileName string `json:"error_filename"`
	DebugFileName string `json:"debug_filename"`
}

type TsQueue struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	QueueName string `json:"queue_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type TsJwt struct {
	Timeout    int `json:"timeout"`
	MaxRefresh int `json:"max_refresh"`
}

type TsTimezone struct {
	Location string `json:"location"`
	Format   string `json:"format"`
}

type Configuration struct {
	Http     TsHttp  `json:"http"`
	Queue    TsQueue `json:"queue"`
	Log      TsLog   `json:"log"`
	Accounts []struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"accounts"`
	Jwt      TsJwt      `json:"jwt"`
	AppPath  string     `json:"app_path"`
	Timezone TsTimezone `json:"timezone"`
	APIKeys  []string   `json:"api_keys"` // API Keys for Bento API authentication
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func GetApplicationPath() (string, error) {
	defer RecoverError()

	strPathSeparator := fmt.Sprintf("%c", os.PathSeparator)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	} else {
		currentPath := dir[len(dir)-1:]
		if strings.TrimSpace(currentPath) == strPathSeparator {
			return dir, nil
		} else {
			currentPath := dir + strPathSeparator
			return currentPath, nil
		}
	}
}

func LoadApplicationConfiguration(removeSuffixPath string) (config Configuration, err error) {
	defer RecoverError()

	config.AppPath, err = os.Getwd()
	if err != nil {
		return
	}

	if removeSuffixPath != "" {
		config.AppPath = strings.TrimSuffix(config.AppPath, removeSuffixPath)
	}

	configPath := filepath.Join(config.AppPath, "config.json")
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	if err != nil {
		return
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &config)
		return
	}
}

func ReadJsonBodyRequest(c *gin.Context) (string, error) {
	defer RecoverError()

	if c.Request.Body != nil {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return "", err
		} else {
			if json.Valid(jsonData) {
				return string(jsonData), nil
			} else {
				return "", errors.New("invalid json format")
			}
		}
	} else {
		return "", errors.New("cannot read or undefined request body")
	}
}

func PrettyStruct(data interface{}) (string, error) {
	defer RecoverError()

	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func PrintConsole(strPrint string, strStatus string) {
	defer RecoverError()

	if supportscolor.Stdout().SupportsColor {
		if strings.ToLower(strings.TrimSpace(strStatus)) == "info" {
			fmt.Println(Green + "[INFO]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "error" {
			fmt.Println(Red + "[ERROR]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "warning" {
			fmt.Println(Yellow + "[WARNING]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "logo" {
			fmt.Println(Red + strPrint + Reset)
		} else {
			fmt.Println(Reset + strPrint + Reset)
		}
	} else {
		if strings.ToLower(strings.TrimSpace(strStatus)) == "info" {
			fmt.Println("[INFO]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "error" {
			fmt.Println("[ERROR]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "warning" {
			fmt.Println("[WARNING]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "logo" {
			fmt.Println(strPrint)
		} else {
			fmt.Println(strPrint)
		}
	}
}

func FormatResponse(c *gin.Context, strResponseFormat string, intHttpResponseCode int, objData gin.H) {
	defer RecoverError()

	switch strResponseFormat {
	case "json":
		c.JSON(intHttpResponseCode, objData)
	case "jsonp":
		c.JSONP(intHttpResponseCode, objData)
	case "xml":
		c.XML(intHttpResponseCode, objData)
	case "yaml":
		c.YAML(intHttpResponseCode, objData)
	case "securejson":
		c.SecureJSON(intHttpResponseCode, objData)
	case "asciijson":
		c.AsciiJSON(intHttpResponseCode, objData)
	case "indentedjson":
		c.IndentedJSON(intHttpResponseCode, objData)
	case "purejson":
		c.PureJSON(intHttpResponseCode, objData)
	default:
		c.JSON(intHttpResponseCode, objData)
	}
}

func FormatResponseJson(c *gin.Context, strResponseFormat string, intHttpResponseCode int, objData interface{}) {
	defer RecoverError()

	switch strResponseFormat {
	case "json":
		c.JSON(intHttpResponseCode, objData)
	case "jsonp":
		c.JSONP(intHttpResponseCode, objData)
	case "xml":
		c.XML(intHttpResponseCode, objData)
	case "yaml":
		c.YAML(intHttpResponseCode, objData)
	case "securejson":
		c.SecureJSON(intHttpResponseCode, objData)
	case "asciijson":
		c.AsciiJSON(intHttpResponseCode, objData)
	case "indentedjson":
		c.IndentedJSON(intHttpResponseCode, objData)
	case "purejson":
		c.PureJSON(intHttpResponseCode, objData)
	default:
		c.JSON(intHttpResponseCode, objData)
	}
}

func InArray(val string, array []string) (exists bool, index int) {
	defer RecoverError()

	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}

func StringToMD5(text string) string {
	defer RecoverError()

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RecoverError() {
	if r := recover(); r != nil {
		PrintConsole(fmt.Sprintf("[ERROR][RECOVER]=> %v", r), "error")
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-Auth-Token, Content-Type, Content-Length, Authorization, Access-Control-Allow-Headers, Accept, Access-Control-Allow-Methods, Access-Control-Allow-Origin, Access-Control-Allow-Credentials")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, HEAD, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}

func SafeString(sql string) string {
	dest := make([]byte, 0, 2*len(sql))

	for i := 0; i < len(sql); i++ {
		c := sql[i]

		if c == '\'' {
			dest = append(dest, '\'', '\'')
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

// SafeJsonString escapes JSON string untuk SQL literal dengan proper handling backslash
// Menggunakan PostgreSQL E” string literal syntax untuk mendukung escape sequences
// Input: {"name":"John","email":"test@example.com"}
// Output: E'{"name":"John","email":"test@example.com"}'
func SafeJsonString(jsonStr string) string {
	// Escape backslash terlebih dahulu (harus yang pertama)
	result := strings.ReplaceAll(jsonStr, "\\", "\\\\")

	// Escape single quote
	result = strings.ReplaceAll(result, "'", "''")

	// Wrap dengan E'' PostgreSQL escape string literal syntax
	return "E'" + result + "'"
}
