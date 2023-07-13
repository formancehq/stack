package fctl

import (
	"encoding/json"
	"flag"

	"github.com/pkg/errors"
)

type persistedConfig struct {
	CurrentProfile string              `json:"currentProfile"`
	Profiles       map[string]*Profile `json:"profiles"`
	UniqueID       string              `json:"uniqueID"`
}

type Config struct {
	currentProfile string
	uniqueID       string
	profiles       map[string]*Profile
	manager        *ConfigManager
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(persistedConfig{
		CurrentProfile: c.currentProfile,
		Profiles:       c.profiles,
		UniqueID:       c.uniqueID,
	})
}

func (c *Config) UnmarshalJSON(data []byte) error {
	cfg := &persistedConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}
	*c = Config{
		currentProfile: cfg.CurrentProfile,
		profiles:       cfg.Profiles,
		uniqueID:       cfg.UniqueID,
	}
	return nil
}

func (c *Config) GetProfile(name string) *Profile {
	p := c.profiles[name]
	if p != nil {
		p.config = c
	}
	return p
}

func (c *Config) GetProfileOrDefault(name string, membershipUri string) *Profile {
	p := c.GetProfile(name)
	if p == nil {
		if c.profiles == nil {
			c.profiles = map[string]*Profile{}
		}
		f := &Profile{
			membershipURI: membershipUri,
			config:        c,
		}
		c.profiles[name] = f
		return f
	}
	return p
}

func (c *Config) DeleteProfile(s string) error {
	_, ok := c.profiles[s]
	if !ok {
		return errors.New("not found")
	}
	delete(c.profiles, s)
	return nil
}

func (c *Config) Persist() error {
	return c.manager.UpdateConfig(c)
}

func (c *Config) SetCurrentProfile(name string, profile *Profile) {
	c.profiles[name] = profile
	c.currentProfile = name
}

func (c *Config) SetUniqueID(id string) {
	c.uniqueID = id
}

func (c *Config) SetProfile(name string, profile *Profile) {
	c.profiles[name] = profile
}

func (c *Config) GetUniqueID() string {
	return c.uniqueID
}

func (c *Config) GetProfiles() map[string]*Profile {
	return c.profiles
}

func (c *Config) GetCurrentProfileName() string {
	return c.currentProfile
}

func (c *Config) SetCurrentProfileName(s string) {
	c.currentProfile = s
}

func GetConfig(flagSet *flag.FlagSet) (*Config, error) {
	return GetConfigManager(flagSet).Load()
}

func GetConfigManager(flagSet *flag.FlagSet) *ConfigManager {
	return NewConfigManager(GetString(flagSet, FileFlag))
}

func GetCurrentProfileName(flags *flag.FlagSet, config *Config) string {
	if profile := GetString(flags, ProfileFlag); profile != "" {
		return profile
	}
	currentProfileName := config.GetCurrentProfileName()
	if currentProfileName == "" {
		currentProfileName = "default"
	}
	return currentProfileName
}

func GetCurrentProfile(flags *flag.FlagSet, cfg *Config) *Profile {
	return cfg.GetProfileOrDefault(GetCurrentProfileName(flags, cfg), DefaultMembershipURI)
}
