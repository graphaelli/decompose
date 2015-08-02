package decompose

import "testing"

var fileServiceMap = map[string][]ComposeService{
	"testdata/order.yml": []ComposeService{
		{Name: "one",
			Image: "image1"},
		{Name: "two",
			Image: "image2"},
		{Name: "three",
			Image: "image3"},
		{Name: "four",
			Image: "image4"},
		{Name: "five",
			Image: "image5"},
		{Name: "six",
			Image: "image6"},
		{Name: "seven",
			Image: "image7"},
		{Name: "eight",
			Image: "image8"},
		{Name: "nine",
			Image: "image9"},
		{Name: "ten",
			Image: "image10"},
	},
	"testdata/all.yml": []ComposeService{
		{Name: "all",
			Build:         "/path/to/build/dir",
			CapAdd:        []string{"ALL"},
			CapDrop:       []string{"NET_ADMIN", "SYS_ADMIN"},
			Command:       "bundle exec thin -p 3000",
			CPUSet:        "0,1",
			CPUShares:     "73",
			Devices:       []string{"/dev/ttyUSB0:/dev/ttyUSB0"},
			DNS:           []string{"8.8.8.8", "9.9.9.9"},
			DNSSearch:     []string{"dc1.example.com", "dc2.example.com"},
			Entrypoint:    "/code/entrypoint.sh",
			EnvFile:       []string{"./common.env", "./apps/web.env", "/opt/secrets.env"},
			Environment:   []string{"RACK_ENV=development", "SESSION_SECRET"},
			Expose:        []string{"3000", "8000"},
			ExternalLinks: []string{"redis_1", "project_db_1:mysql", "project_db_1:postgresql"},
			ExtraHosts:    []string{"somehost:162.242.195.82", "otherhost:50.31.209.229"},
			Hostname:      "foo",
			Image:         "ubuntu",
			Labels:        []string{"com.example.description=Accounting webapp", "com.example.department=Finance", "com.example.label-with-empty-value"},
			Links:         []string{"db", "db:database", "redis"},
			LogDriver:     "json-file",
			MemLimit:      "1000000000",
			Net:           "bridge",
			PID:           "host",
			Ports:         []string{"3000", "8000:8000", "49100:22", "127.0.0.1:8001:8001"},
			Privileged:    true,
			ReadOnly:      true,
			Restart:       "always",
			SecurityOpt:   []string{"label:user:USER", "label:role:ROLE"},
			StdinOpen:     true,
			TTY:           true,
			User:          "postgresql",
			Volumes:       []string{"/var/lib/mysql", "cache/:/tmp/cache", "~/configs:/etc/configs/:ro"},
			VolumesFrom:   []string{"service_name", "container_name"},
			WorkingDir:    "/code",
		},
	},
}

func TestComposeServiceString(t *testing.T) {
	expCommand := []string{
		"docker run --name=one image1",
		"docker run --name=two image2",
		"docker run --name=three image3",
		"docker run --name=four image4",
		"docker run --name=five image5",
		"docker run --name=six image6",
		"docker run --name=seven image7",
		"docker run --name=eight image8",
		"docker run --name=nine image9",
		"docker run --name=ten image10",
	}

	// more testing required!

	expServs := fileServiceMap["testdata/order.yml"]
	for i := range expServs {
		if expServs[i].String() != expCommand[i] {
			t.Error("unexpected docker command")
			t.Errorf("expected: %s; found: %s\n", expServs[i], expCommand[i])
			t.FailNow()
		}
	}
}

func TestParseComposeFile(t *testing.T) {
	for fp, expServs := range fileServiceMap {
		servs, err := ParseComposeFile(fp)
		if err != nil {
			t.Fatal(err)
		}

		if len(expServs) != len(servs) {
			t.Error("unexpected services count")
			t.Errorf("expected: %d; found: %d\n", len(expServs), len(servs))
			t.FailNow()
		}
	}

}
