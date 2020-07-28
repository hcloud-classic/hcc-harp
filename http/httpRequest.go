package http

import (
	"encoding/json"
	"errors"
	harpData "hcc/harp/data"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getModuleHTTPInfo(moduleName string) (time.Duration, string, error) {
	var timeout time.Duration
	var url = "http://"
	switch moduleName {
	case "flute":
		timeout = time.Duration(config.Flute.RequestTimeoutMs)
		url += config.Flute.ServerAddress + ":" + strconv.Itoa(int(config.Flute.ServerPort))
		break
	case "violin":
		timeout = time.Duration(config.Violin.RequestTimeoutMs)
		url += config.Violin.ServerAddress + ":" + strconv.Itoa(int(config.Violin.ServerPort))
		break
	default:
		return 0, "", errors.New("unknown module name")
	}

	return timeout, url, nil
}

// DoHTTPRequest : Send http request to other modules with GraphQL query string.
func DoHTTPRequest(moduleName string, needData bool, data interface{}, query string) (interface{}, error) {
	timeout, url, err := getModuleHTTPInfo(moduleName)
	if err != nil {
		return nil, err
	}

	url += "/graphql?query=" + queryURLEncoder(query)
	client := &http.Client{Timeout: timeout * time.Millisecond}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(config.HTTP.RequestRetryCount); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("DoHTTPRequest(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + moduleName + " module")
				logger.Logger.Println("DoHTTPRequest(): Failed info - ModuleName=" + moduleName + ", Query=" + query)
			}
			logger.Logger.Println("DoHTTPRequest(): Will retry for " + moduleName + " module after " +
				strconv.Itoa(int(config.HTTP.RequestRetryDelaySec)) + " seconds " +
				"(" + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.HTTP.RequestRetryCount)) + ")")
			time.Sleep(time.Duration(config.HTTP.RequestRetryDelaySec) * time.Second)
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				result := string(respBody)

				if strings.Contains(result, "errors") {
					return nil, errors.New(result)
				}
				if needData {
					if data == nil {
						return nil, errors.New("needData marked as true but data is nil")
					}

					switch moduleName {
					case "flute":
						listNodeData := data.(harpData.ListNodeData)
						err = json.Unmarshal([]byte(result), &listNodeData)
						if err != nil {
							return nil, err
						}
						return listNodeData, nil
					case "violin":
						allServerData := data.(harpData.AllServerData)
						err = json.Unmarshal([]byte(result), &allServerData)
						if err != nil {
							return nil, err
						}

						return allServerData, nil
					default:
						return nil, errors.New("DoHTTPRequest(): data is not supported for " + moduleName + " module")
					}
				}

				return result, nil
			}

			return nil, err
		}
	}

	return nil, errors.New("DoHTTPRequest(): retry count exceeded for " + moduleName + " module")
}
