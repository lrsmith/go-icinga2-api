
docker_start:
	docker pull jordan/icinga2
	docker run -d --name icinga2 -p 8080:80 -p 8443:443 -p 5665:5665 -it jordan/icinga2:latest

docker_add_user:
	docker exec icinga2 bash -c 'echo -e "object ApiUser \"icinga-test\" {\n  password = \"icinga\"\n  permissions = [ \"*\" ]\n}" >> /etc/icinga2/conf.d/api-users.conf'

docker_restart:
	docker exec icinga2 supervisorctl restart icinga2

docker_clean:
	docker stop icinga2
	docker rm icinga2

docker_setup: docker_start docker_add_user docker_restart

test:
	go test ./...
