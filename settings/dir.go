package settings

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

var (
	invalidFilenameChars = regexp.MustCompile(`[^0-9A-Za-z@_\-.]`)

	dashes = regexp.MustCompile(`[\-]+`)
)

// MakeUserDir makes the user directory according to settings.
func (settings *Settings) MakeUserDir(username, userScope, serverRoot string) (string, error) {
	var err error
	userScope = strings.TrimSpace(userScope)
	if userScope == "" || userScope == "./" {
		userScope = "."
	}

	if !settings.CreateUserDir {
		return userScope, nil
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), serverRoot)

	// Use the default auto create logic only if specific scope is not the default scope
	if userScope != settings.Defaults.Scope {
		// Try create the dir, for example: settings.Defaults.Scope == "." and userScope == "./foo"
		if userScope != "." {
			err = fs.MkdirAll(userScope, os.ModePerm)
			if err != nil {
				log.Printf("create user: failed to mkdir user home dir: [%s]", userScope)
			}
		}
		return userScope, err
	}

	// Clean username first
	username = cleanUsername(username)
	if username == "" || username == "-" || username == "." {
		log.Printf("create user: invalid user for home dir creation: [%s]", username)
		return "", errors.New("invalid user for home dir creation")
	}

	// Create default user dir
	userHomeBase := settings.Defaults.Scope + string(os.PathSeparator) + "users"
	userHome := userHomeBase + string(os.PathSeparator) + username
	err = fs.MkdirAll(userHome, os.ModePerm)
	if err != nil {
		log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
	} else {
		log.Printf("create user: mkdir user home dir: [%s] successfully.", userHome)
	}
	return userHome, err
}

func cleanUsername(s string) string {
	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")
	s = strings.Replace(s, "..", "", -1)

	// Replace all characters which not in the list `0-9A-Za-z@_\-.` with a dash
	s = invalidFilenameChars.ReplaceAllString(s, "-")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")
	return s
}

// MakeSpaceDir makes the space directory according to settings.
func (settings *Settings) MakeSpaceDir(userSpace, serverRoot string) (string, error) {
	var err error
	userSpace = strings.TrimSpace(userSpace)
	if userSpace == "" || userSpace == "./" {
		userSpace = "."
	}

	//if !settings.CreateUserDir {
	//	return userScope, nil
	//}

	fs := afero.NewBasePathFs(afero.NewOsFs(), serverRoot)
	afs := &afero.Afero{Fs: fs}

	// Clean username first
	spacename := cleanSpacename(userSpace)
	if spacename == "" || spacename == "-" || spacename == "." {
		log.Printf("create space: invalid space dir creation: [%s]", spacename)
		return "", errors.New("invalid space dir creation")
	}

	// Create default user dir
	spaceBase := settings.Defaults.Space + string(os.PathSeparator) + "spaces"
	space := spaceBase + string(os.PathSeparator) + spacename
	dirCheck, err := afs.DirExists(space)
	if dirCheck {
		log.Printf("create space dir: space dir is already exists: [%s]", space)
	} else {
		err = fs.MkdirAll(space, os.ModePerm)
		if err != nil {
			log.Printf("create space dir: failed to mkdir space dir: [%s]", space)
		} else {
			log.Printf("create space dir: mkdir space dir: [%s] successfully.", space)
		}
	}
	return space, err
}

func cleanSpacename(s string) string {
	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")
	s = strings.Replace(s, "..", "", -1)

	// Replace all characters which not in the list `0-9A-Za-z@_\-.` with a dash
	s = invalidFilenameChars.ReplaceAllString(s, "-")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")
	return s
}
