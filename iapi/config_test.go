package iapi

import (
	"fmt"
	"testing"
)

func TestConfigListPackages(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	packages, err := Icinga2_Server.GetConfigPackages()

	if len(packages) < 1 {
		t.Error(fmt.Errorf("%d packages found", len(packages)))
	}

	if err != nil {
		t.Error(err)
	}
}

func TestConfigGetPackage(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	packages, err := Icinga2_Server.GetConfigPackages()

	if len(packages) < 1 {
		t.Error(fmt.Errorf("%d packages found", len(packages)))
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	pack, packErr := Icinga2_Server.GetConfigPackage(packages[0].Name)

	if packErr != nil {
		t.Error(packErr)
		return
	}

	if pack[0].Name != packages[0].Name {
		t.Error(fmt.Errorf("Package name does not match request"))
		return
	}
}

func TestConfigCreatePackage(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	var packagename = "test-package"

	packages, err := Icinga2_Server.CreateConfigPackage(packagename)

	if err != nil {
		t.Error(err)
		return
	}

	if packagename != packages[0].Name {
		t.Error(fmt.Errorf("Package name does not match request"))
		return
	}
}

func TestConfigCreateStage(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	var packagename = "test-package"

	var files = make(map[string]string)
	files["conf.d/test-host.conf"] = `object Host "test-host" {
		display_name = "Test Host"
		address = "127.0.0.1"

		check_command = "hostalive"
	}`

	stages, err := Icinga2_Server.CreateConfigStageRetriable(
		packagename,
		files,
		true,
		true,
		5,
		1000,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if !stages[0].Successful {
		t.Error("Stage unsuccessful")
		return
	}
}

func TestConfigListStageFiles(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	var packagename = "test-package"

	packages, err := Icinga2_Server.GetConfigPackage(packagename)

	if err != nil {
		t.Error(err)
		return
	}

	files, err := Icinga2_Server.ListConfigStageFiles(packagename, packages[0].ActiveStage)

	for _, file := range files {
		if file.Name == "conf.d/test-host.conf" && file.Type == "file" {
			return
		}
	}

	t.Error("Configuration file not found in active stage!")
}

func TestConfigCreateBrokenStage(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	var packagename = "test-package"

	var files = make(map[string]string)
	files["conf.d/test-host.conf"] = `object DoesNotExist "test-host" {
		display_name = "Test Host"
		address = "127.0.0.1"

		check_command = "hostalive"
	}`

	stages, err := Icinga2_Server.CreateConfigStageRetriable(
		packagename,
		files,
		true,
		true,
		5,
		1000,
	)

	if err != nil {
		t.Error(err)
		return
	}

	if stages[0].Successful {
		t.Error("Broken Stage Error was not detected!")
	}
}

func TestConfigDeletePackage(t *testing.T) {
	var Icinga2_Server = Server{ICINGA2_API_USER_ROOT, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}

	var packagename = "test-package"

	err := Icinga2_Server.DeleteConfigPackage(packagename)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = Icinga2_Server.GetConfigPackage(packagename)

	if err == nil {
		t.Error("Package could be found after delete")
	}
}
