package iapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (server *Server) GetConfigPackages() ([]ConfigPackageStruct, error) {
	var packages []ConfigPackageStruct

	results, err := server.NewAPIRequest("GET", "/config/packages", nil)
	if err != nil {
		return nil, err
	}

	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	if unmarshalErr := json.Unmarshal(jsonStr, &packages); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return packages, err
}

func (server *Server) GetConfigPackage(packagename string) ([]ConfigPackageStruct, error) {
	var packages []ConfigPackageStruct

	packages, err := server.GetConfigPackages()

	if err != nil {
		return nil, err
	}
	if packages == nil {
		return nil, fmt.Errorf("Empty struct")
	}
	for _, item := range packages {
		if item.Name == packagename {
			var packResult []ConfigPackageStruct
			return append(packResult, item), nil
		}
	}

	return nil, fmt.Errorf("Package %s not found", packagename)
}

func (server *Server) CreateConfigPackage(packagename string) ([]ConfigPackageStruct, error) {
	results, err := server.NewAPIRequest("POST", "/config/packages/"+packagename, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		return server.GetConfigPackage(packagename)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

func (server *Server) DeleteConfigPackage(packagename string) error {
	results, err := server.NewAPIRequest("DELETE", "/config/packages/"+packagename, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}

func (server *Server) CreateConfigStage(packagename string, files map[string]string, reload bool, activate bool) ([]ConfigStageStruct, error) {
	if !activate && reload {
		return nil, fmt.Errorf("if activate is set to false, reload must also be set to false")
	}

	var newAttrs ConfigStageAttrs
	newAttrs.Files = files
	newAttrs.Activate = activate
	newAttrs.Reload = reload

	payloadJSON, marshalErr := json.Marshal(newAttrs)
	if marshalErr != nil {
		return nil, marshalErr
	}

	results, err := server.NewAPIRequest("POST", "/config/stages/"+packagename, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		jsonStr, marshalErr := json.Marshal(results.Results)
		if marshalErr != nil {
			return nil, marshalErr
		}

		var stages []ConfigStageStruct

		if unmarshalErr := json.Unmarshal(jsonStr, &stages); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return stages, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

func (server *Server) DetermineLogStatus(log string) bool {
	lines := strings.Split(log, "\n")
	var lastLine string
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			lastLine = line
		}
	}
	return strings.Contains(lastLine, "Finished")
}

func (server *Server) CreateConfigStageRetriable(packagename string, files map[string]string, reload bool, activate bool, retryCount int, delayMS time.Duration) ([]ConfigStageStruct, error) {
	var globalerr error
	for i := 1; i < retryCount; i++ {
		stages, err := server.CreateConfigStage(packagename, files, reload, activate)
		if err == nil {
			log, err := server.GetConfigStageFile(packagename, stages[0].Name, "startup.log")
			if err != nil {
				return nil, err
			}
			stages[0].Log = log
			stages[0].Successful = server.DetermineLogStatus(log)

			return stages, nil
		}
		time.Sleep(delayMS * time.Millisecond)
		globalerr = err
	}
	return nil, globalerr
}

func (server *Server) ListConfigStageFiles(packagename string, stage string) ([]ConfigStageFileStruct, error) {
	results, err := server.NewAPIRequest("GET", "/config/stages/"+packagename+"/"+stage, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		jsonStr, marshalErr := json.Marshal(results.Results)
		if marshalErr != nil {
			return nil, marshalErr
		}

		var files []ConfigStageFileStruct

		if unmarshalErr := json.Unmarshal(jsonStr, &files); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return files, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

func (server *Server) GetConfigStageFile(packagename string, stage string, file string) (string, error) {
	result, err := server.NewPlainTextRequest("GET", "/config/files/"+packagename+"/"+stage+"/"+file, nil)
	if err != nil {
		return "", err
	}

	return result, nil
}
