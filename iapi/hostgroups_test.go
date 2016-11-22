package iapi

import "testing"

func TestGetValidHostgroup(t *testing.T) {

	name := "linux-servers"

	hostgroup, err := VagrantImage.GetHostgroup(name)

	if (err != nil) || (hostgroup == nil) {
		t.Errorf("Error : Failed to find hostgroup %s : ( %s <> %s ) ", name, err, hostgroup)
	}

	if hostgroup.Name != name {
		t.Errorf("Error : Did not get expected hostname. ( %s != %s )", hostgroup.Name, name)
	}

}

func TestGetInvalidHostgroup(t *testing.T) {

	name := "irix-servers"

	hostgroup, err := VagrantImage.GetHostgroup(name)
	if err != nil && hostgroup != nil {
		t.Errorf("Error : Did not get empty list of hostgroup. ( %v : %s )", err, hostgroup)
	}

}

func TestCreateHostGroup(t *testing.T) {

	name := "docker-servers"
	displayName := "Docker Host Servers"
	_, err := VagrantImage.CreateHostgroup(name, displayName)

	if err != nil {
		t.Errorf("Error : Failed to create hostgroup %s : %s", name, err)
	}

}

func TestDeleteHostgroup(t *testing.T) {

	name := "docker-servers"

	err := VagrantImage.DeleteHostgroup(name)
	if err != nil {
		t.Errorf("Error : Failed to delete hostgroup %s : %s", name, err)
	}

}

func TestDeleteNonExistentHostGroup(t *testing.T) {

	name := "docker-servers"
	err := VagrantImage.DeleteHostgroup(name)
	if err != nil {
		t.Errorf("Error : Failed to delete hostgroup %s : %s", name, err)
	}

}
