package icinga2

import "testing"

func TestGetValidHostgroup(t *testing.T) {

	name := "linux-servers"

	hostgroup, err := VagrantImage.GetHostgroup(name)

	if (err != nil) || (hostgroup == nil) {
		t.Errorf("Error : Failed to find hostgroup %s : ( %s <> %s ) ", name, err, hostgroup)
	}

	if hostgroup[0].Name != name {
		t.Errorf("Error : Did not get expected hostname. ( %s != %s )", hostgroup[0].Name, name)
	}

}

func TestGetInvalidHostgroup(t *testing.T) {

	name := "irix-servers"

	hostgroup, err := VagrantImage.GetHostgroup(name)
	if (err != nil) || (len(hostgroup) != 0) {
		t.Errorf("Error : Did not get empty list of hostgroup. ( %v : %s )", err, hostgroup)
	}

}

func TestCreateHostGroup(t *testing.T) {

	name := "docker-servers"

	err := VagrantImage.CreateHostgroup(name)

	if err != nil {
		t.Errorf("Error : Failed to create hostgroup %s : %s", name, err)
	}

}

/*
func TestDeleteHost(t *testing.T) {

	hostname := "go-icinga2-api-1"

	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

}

func TestDeleteNonExistentHost(t *testing.T) {
	hostname := "go-icinga2-api-1"
	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

} */
