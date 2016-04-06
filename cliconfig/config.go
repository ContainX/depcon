// Provides storage and retrieval of user preferences and cluster environment configuration
package cliconfig

// Some functions have been imported for docker/cliconfig

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gondor/depcon/pkg/userdir"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigFileName = "config.json"
	DotConfigDir   = ".depcon"
	TypeMarathon   = "marathon"
	TypeKubernetes = "kubernetes"
	TypeECS        = "ecs"
)

var (
	configDir      = os.Getenv("DEPCON_CONFIG")
	ErrEnvNotFound = errors.New("Specified environment was not found")
)

func init() {
	if configDir == "" {
		configDir = filepath.Join(userdir.Get(), DotConfigDir)
	}
}

// ConfigDir returns the directory the configuration file is stored in
func ConfigDir() string {
	return configDir
}

// SetConfigDir sets the directory the configuration file is stored in
func SetConfigDir(dir string) {
	configDir = dir
}

type ConfigFile struct {
	Format       string                        `json:"format,omitempty"`
	RootService  bool                          `json:"rootservice"`
	Environments map[string]*ConfigEnvironment `json:"environments,omitempty"`
	DefaultEnv   string                        `json:"default,omitempty"`
	filename     string                        // not serialized
}

type ConfigEnvironment struct {
	// currently only supporting marathon as initial release
	Marathon *ServiceConfig `json:"marathon,omitempty"`
}

type ServiceConfig struct {
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	HostUrl  string            `json:"serveraddress,omitempty"`
	Features map[string]string `json:"features,omitempty"`
	Name     string            `json:"-"`
}

func HasExistingConfig() (*ConfigFile, bool) {
	configFile, err := Load("")
	return configFile, err == nil
}

func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}
	var err error
	for _, configEnv := range configFile.Environments {
		configEnv.Marathon.Password, err = DecodePassword(configEnv.Marathon.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func Load(configDir string) (*ConfigFile, error) {
	if configDir == "" {
		configDir = ConfigDir()
	}

	configFile := ConfigFile{
		Format:       "column",
		Environments: make(map[string]*ConfigEnvironment),
		filename:     filepath.Join(configDir, ConfigFileName),
	}

	_, err := os.Stat(configFile.filename)
	if err == nil {
		file, err := os.Open(configFile.filename)
		if err != nil {
			return &configFile, err
		}

		defer file.Close()
		err = configFile.LoadFromReader(file)
		return &configFile, err
	}
	return &configFile, err
}

// Determines if the configuration has only a single environment defined and the user prefers a rooted service
// Returns the environment type and true or else "" and false
func (configFile *ConfigFile) DetermineIfServiceIsRooted() (string, bool) {
	if len(configFile.Environments) > 1 {
		return "", false
	}

	// FIXME: Once we support more than Marathon (post initial release) - determine type
	return TypeMarathon, configFile.RootService
}

func (configEnv *ConfigEnvironment) EnvironmentType() string {
	// FIXME: Once we support more than Marathon (post initial release) - determine type
	return TypeMarathon
}

func (configFile *ConfigFile) AddEnvironment() {
	serviceEnv := createEnvironment()
	configEnv := &ConfigEnvironment{
		Marathon: serviceEnv,
	}
	configFile.Environments[serviceEnv.Name] = configEnv
	configFile.Save()
}

// Removes the specified environment from the configuration
// {name}  - name of the environment
// {force} - if true will not prompt for confirmation
//
// Will return ErrEnvNotFound if the environment could not be found
func (configFile *ConfigFile) RemoveEnvironment(name string, force bool) error {
	configEnv := configFile.Environments[name]
	if configEnv == nil {
		return ErrEnvNotFound
	}
	if !force {
		if !getBoolAnswer(fmt.Sprintf("Are you sure you would like to remove '%s'", name), true) {
			return nil
		}
	}
	if configFile.DefaultEnv == name {
		configFile.DefaultEnv = ""
	}
	delete(configFile.Environments, name)
	configFile.Save()

	return nil
}

// Sets the default environment to use.  This is used by other parts of the application to eliminate the user always
// specifying and environment
// {name} - the environment name
//
// Will return ErrEnvNOtFound if the environment could not be found
func (configFile *ConfigFile) SetDefaultEnvironment(name string) error {
	configEnv := configFile.Environments[name]
	if configEnv == nil {
		return ErrEnvNotFound
	}
	configFile.DefaultEnv = name
	configFile.Save()
	return nil
}

// Renames an environment and updates the default environment if it matches the current old
// {oldName} - the old environment name
// {newName} - the new environment name
//
// Will return ErrEnvNOtFound if the old environment could not be found
func (configFile *ConfigFile) RenameEnvironment(oldName, newName string) error {
	configEnv := configFile.Environments[oldName]
	if configEnv == nil {
		return ErrEnvNotFound
	}
	delete(configFile.Environments, oldName)
	configFile.Environments[newName] = configEnv

	if configFile.DefaultEnv == oldName {
		configFile.DefaultEnv = newName
	}
	configFile.Save()
	return nil
}

// Returns the Configuration for the specified environment.  If the environment
// is not found then
func (configFile *ConfigFile) GetEnvironment(name string) (*ConfigEnvironment, error) {
	configEnv := configFile.Environments[name]
	if configEnv == nil {
		return nil, ErrEnvNotFound
	}
	return configEnv, nil
}

func (configFile *ConfigFile) GetEnvironments() []string {
	keys := make([]string, 0, len(configFile.Environments))
	for k := range configFile.Environments {
		keys = append(keys, k)
	}
	return keys
}

func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {
	tmpEnvConfigs := make(map[string]*ConfigEnvironment, len(configFile.Environments))
	for k, configEnv := range configFile.Environments {
		configEnvCopy := configEnv

		if configEnvCopy.Marathon != nil {
			configEnvCopy.Marathon.Password = EncodePassword(configEnvCopy.Marathon)
			configEnvCopy.Marathon.Name = ""
		}
		tmpEnvConfigs[k] = configEnvCopy
	}
	saveEnvConfigs := configFile.Environments
	configFile.Environments = tmpEnvConfigs

	defer func() { configFile.Environments = saveEnvConfigs }()

	data, err := json.MarshalIndent(configFile, "", "\t")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (configFile *ConfigFile) Save() error {
	if configFile.Filename() == "" {
		configFile.filename = filepath.Join(configDir, ConfigFileName)
	}

	if err := os.MkdirAll(filepath.Dir(configFile.filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return configFile.SaveToWriter(f)
}

// EncodePassword creates a base64 encoded string using the authorization info.  If the username
// is not specified then an empty password is returned since we need both
func EncodePassword(auth *ServiceConfig) string {
	if auth.Username == "" {
		return ""
	}
	authStr := auth.Username + ":" + auth.Password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

func DecodePassword(authStr string) (string, error) {

	if authStr == "" {
		return "", nil
	}

	decLen := base64.StdEncoding.DecodedLen(len(authStr))
	decoded := make([]byte, decLen)
	authByte := []byte(authStr)
	n, err := base64.StdEncoding.Decode(decoded, authByte)
	if err != nil {
		return "", err
	}
	if n > decLen {
		return "", fmt.Errorf("Something went wrong decoding service authentication")
	}
	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", fmt.Errorf("Invalid auth configuration file")
	}
	password := strings.Trim(arr[1], "\x00")
	return password, nil
}

// Filename returns the name of the configuration file
func (configFile *ConfigFile) Filename() string {
	return configFile.filename
}
