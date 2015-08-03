package decompose

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	// Prefix is the string that will be inserted before service names.
	Prefix string
	// Generate docker run commands with detach flag.
	Detach bool
	// Generate docker run commands with remove flag.
	Remove bool
)

// ComposeService defines one service declared in a Compose YAML file.
// Currently, dockerfile, domainname, and extends are not supported.
type ComposeService struct {
	Build     string   `yaml:"build"`      // dockerfile
	CapAdd    []string `yaml:"cap_add"`    // --cap-add=[]
	CapDrop   []string `yaml:"cap_drop"`   // --cap-drop=[]
	Command   string   `yaml:"command"`    // [command]
	CPUSet    string   `yaml:"cpuset"`     // --cpuset-cpus=
	CPUShares string   `yaml:"cpu_shares"` // --cpu-shares=
	Devices   []string `yaml:"devices"`    // --device=[]
	DNS       []string `yaml:"dns"`        // --dns=[]
	DNSSearch []string `yaml:"dns_search"` // --dns-search=[]
	// Dockerfile    string   `yaml:"dockerfile"`  // dockerfile (another)
	// DomainName    string   `yaml:"domainname"`  // ??
	Entrypoint  string   `yaml:"entrypoint"`  // --entrypoint=
	EnvFile     []string `yaml:"env_file"`    // --env-file=[]
	Environment []string `yaml:"environment"` // --env=[]
	Expose      []string `yaml:"expose"`      // --expose=[]
	// Extends       string   `yaml:"extends"`
	ExternalLinks []string `yaml:"external_links"` // --link=[]
	ExtraHosts    []string `yaml:"extra_hosts"`    // --add-host=[]
	Hostname      string   `yaml:"hostname"`       // --hostname=
	Image         string   `yaml:"image"`          // image
	Labels        []string `yaml:"labels"`         // --label=[]
	Links         []string `yaml:"links"`          // --link=[]
	LogDriver     string   `yaml:"log_driver"`     // --log-driver=
	MemLimit      string   `yaml:"mem_limit"`      // --memory=
	Name          string   // --name=
	Net           string   `yaml:"net"`          // --net=bridge
	PID           string   `yaml:"pid"`          // --pid=
	Ports         []string `yaml:"ports"`        // --publish=[]
	Privileged    bool     `yaml:"privileged"`   // --privileged=false
	ReadOnly      bool     `yaml:"read_only"`    // --read-only=false
	Restart       string   `yaml:"restart"`      // --restart=no
	SecurityOpt   []string `yaml:"security_opt"` // --security-opt=[]
	StdinOpen     bool     `yaml:"stdin_open"`   // --interactive=false
	TTY           bool     `yaml:"tty"`          // --tty=false
	User          string   `yaml:"user"`         // --user=
	Volumes       []string `yaml:"volumes"`      // --volume=[]
	VolumesFrom   []string `yaml:"volumes_from"` // --volumes-from=[]
	WorkingDir    string   `yaml:"working_dir"`  // --workdir=
}

// String returns the docker run commands to recreate a service. If a service
// specifies build, then a docker build command will also be created.
func (c ComposeService) String() string {
	var buf bytes.Buffer

	// docker build [options] dockerfile
	if c.Build != "" {
		buf.WriteString(fmt.Sprintf("docker build --tag %s %s\n", c.Name, c.Build))
		c.Image = c.Name
	}

	// docker run [options] image [command] [arg...]
	buf.WriteString("docker run ")

	// [options]
	if v := c.Name; v != "" {
		buf.WriteString("--name=" + Prefix + v + " ")
	}
	for _, v := range c.CapAdd {
		buf.WriteString("--cap-add=" + v + " ")
	}
	for _, v := range c.CapDrop {
		buf.WriteString("--cap-drop=" + v + " ")
	}
	if v := c.CPUSet; v != "" {
		buf.WriteString("--cpuset-cpus=" + v + " ")
	}
	if v := c.CPUShares; v != "" {
		buf.WriteString("--cpu-shares=" + v + " ")
	}
	if Detach {
		buf.WriteString("--detach ")
	}
	for _, v := range c.Devices {
		buf.WriteString("--device=" + v + " ")
	}
	for _, v := range c.DNS {
		buf.WriteString("--dns=" + v + " ")
	}
	for _, v := range c.DNSSearch {
		buf.WriteString("--dns-search=" + v + " ")
	}
	if v := c.Entrypoint; v != "" {
		buf.WriteString("--entrypoint=" + v + " ")
	}
	for _, v := range c.EnvFile {
		buf.WriteString("--env-file=" + v + " ")
	}
	for _, v := range c.Environment {
		buf.WriteString("--env=" + v + " ")
	}
	for _, v := range c.Expose {
		buf.WriteString("--env=" + v + " ")
	}
	for _, v := range c.ExternalLinks {
		buf.WriteString("--link=" + v + " ")
	}
	for _, v := range c.ExtraHosts {
		buf.WriteString("--add-host=" + v + " ")
	}
	if v := c.Hostname; v != "" {
		buf.WriteString("--hostname=" + v + " ")
	}
	for _, v := range c.Labels {
		buf.WriteString("--label=" + v + " ")
	}
	for _, v := range c.Links {
		buf.WriteString("--link=" + Prefix + v + " ")
	}
	if v := c.LogDriver; v != "" {
		buf.WriteString("--log-driver=" + v + " ")
	}
	if v := c.MemLimit; v != "" {
		buf.WriteString("--memory=" + v + " ")
	}
	if v := c.Net; v != "" {
		buf.WriteString("--net=" + v + " ")
	}
	if v := c.PID; v != "" {
		buf.WriteString("--pid=" + v + " ")
	}
	for _, v := range c.Ports {
		buf.WriteString("--publish=" + v + " ")
	}
	if c.Privileged {
		buf.WriteString("--privileged ")
	}
	if c.ReadOnly {
		buf.WriteString("--read-only ")
	}
	if Remove {
		buf.WriteString("--rm ")
	}
	if v := c.Restart; v != "" {
		buf.WriteString("--restart=" + v + " ")
	}
	for _, v := range c.SecurityOpt {
		buf.WriteString("--security-opt=" + v + " ")
	}
	if c.StdinOpen {
		buf.WriteString("--interactive ")
	}
	if c.TTY {
		buf.WriteString("--tty ")
	}
	if v := c.User; v != "" {
		buf.WriteString("--user=" + v + " ")
	}
	for _, v := range c.Volumes {
		buf.WriteString("--volume=" + v + " ")
	}
	for _, v := range c.VolumesFrom {
		buf.WriteString("--volumes-from=" + Prefix + v + " ")
	}
	if v := c.WorkingDir; v != "" {
		buf.WriteString("--workdir=" + v + " ")
	}

	// image
	buf.WriteString(c.Image + " ")

	// [command] [arg...]
	buf.WriteString(c.Command)

	return strings.TrimSpace(buf.String())
}

// ParseComposeFile returns a slice of Docker Compose services. The order of the
// services is preserved.
func ParseComposeFile(compPath string) ([]ComposeService, error) {
	bs, err := ioutil.ReadFile(compPath)
	if err != nil {
		return nil, err
	}

	var srvMap map[string]ComposeService
	if err := yaml.Unmarshal(bs, &srvMap); err != nil {
		return nil, err
	}

	var srvOrder yaml.MapSlice
	if err := yaml.Unmarshal(bs, &srvOrder); err != nil {
		return nil, err
	}

	services := make([]ComposeService, 0, len(srvOrder))
	for _, ms := range srvOrder {
		srvName := ms.Key.(string)
		srv := srvMap[srvName]
		srv.Name = srvName

		// required for valid docker-compose file
		if srv.Build == "" && srv.Image == "" {
			return nil, fmt.Errorf("%s must specify either image or build", srv.Name)
		}

		// filepath on host must be absolute path
		for i, vp := range srv.Volumes {
			// /var/lib/mysql
			if !strings.Contains(vp, ":") {
				continue
			}

			// cache/:/tmp/cache:ro
			splits := strings.Split(vp, ":")
			abs, err := filepath.Abs(splits[0])
			if err != nil {
				return nil, err
			}
			splits[0] = abs

			srv.Volumes[i] = strings.Join(splits, ":")
		}

		services = append(services, srv)
	}

	return services, nil
}
