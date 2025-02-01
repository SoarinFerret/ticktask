package profile

import "github.com/spf13/viper"

// Profile is a struct that holds the profile information
type Profile struct {
	Name     string
	Projects []string
	Contexts []string
	Active   bool
}

func GetProfiles() []Profile {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	return profiles
}

func GetActiveProfile() Profile {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	for _, p := range profiles {
		if p.Active {
			return p
		}
	}

	return Profile{}
}

func GetActiveProfileFilter() (projects []string, contexts []string) {
	// get active profile
	activeProfile := GetActiveProfile()

	if activeProfile.Name == "" {
		return
	}

	// get projects and contexts from active profile
	projects = activeProfile.Projects
	contexts = activeProfile.Contexts

	return
}

func SetActiveProfile(name string, commit bool) error {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	// set active profile
	found := false
	for i, p := range profiles {
		if p.Name == name {
			profiles[i].Active = true
			found = true
		} else {
			profiles[i].Active = false
		}
	}

	// if no profile match name, return error
	if !found {
		return &ProfileNotFoundError{name}
	}

	viper.Set("profiles", profiles)
	if commit {
		viper.WriteConfig()
	}
	return nil
}

func UnsetActiveProfile(commit bool) error {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	// unset active profile
	for i := range profiles {
		profiles[i].Active = false
	}

	viper.Set("profiles", profiles)
	if commit {
		viper.WriteConfig()
	}
	return nil
}

func AddProfile(name string, active bool, projects []string, contexts []string) *Profile {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	// add new profile
	profiles = append(profiles, Profile{
		Name:     name,
		Projects: projects,
		Contexts: contexts,
		Active:   active,
	})

	// if active, set all other profiles to inactive
	if active {
		for i := range profiles {
			if profiles[i].Name != name {
				profiles[i].Active = false
			}
		}
	}

	viper.Set("profiles", profiles)
	viper.WriteConfig()

	return &profiles[len(profiles)-1]
}

func EditProfile(name string, projects []string, contexts []string) (*Profile, error) {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	// edit profile
	found := -1
	for i, p := range profiles {
		if p.Name == name {
			// if not empty, replace projects and contexts
			if len(projects) > 0 {
				profiles[i].Projects = projects
			}
			if len(contexts) > 0 {
				profiles[i].Contexts = contexts
			}

			found = i
			continue
		}
	}

	// if no profile match name, return error
	if found == -1 {
		return nil, &ProfileNotFoundError{name}
	}

	viper.Set("profiles", profiles)
	viper.WriteConfig()
	return &profiles[found], nil
}

func RemoveProfile(name string) {
	// get profiles from config
	var profiles []Profile
	viper.UnmarshalKey("profiles", &profiles)

	// remove profile
	for i, p := range profiles {
		if p.Name == name {
			profiles = append(profiles[:i], profiles[i+1:]...)
		}
	}

	viper.Set("profiles", profiles)
	viper.WriteConfig()
}

// profile not found error definition
type ProfileNotFoundError struct {
	Name string
}

func (e *ProfileNotFoundError) Error() string {
	return "Profile not found: " + e.Name
}
