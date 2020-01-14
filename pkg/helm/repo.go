// Modifications copyright (C) 2019 Alibaba Group Holding Limited / Yuning Xie (xyn1016@gmail.com)

package helm

import (
	"fmt"
	urllib "net/url"
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	v2environment "k8s.io/helm/pkg/helm/environment"
	v2helmpath "k8s.io/helm/pkg/helm/helmpath"
)

type (
	// Repo represents a collection of parameters for chart repository
	Repo struct {
		*repo.ChartRepository
	}
)

// GetRepoByURL returns repository by url
func GetRepoByURL(url string) (*Repo, error) {
	r, err := repoFile()
	if err != nil {
		return nil, err
	}
	entry, exists := findRepoEntryByURL(url, r)
	if !exists {
		return nil, fmt.Errorf("no repo url %q found", url)
	}

	// FIXME
	// https://github.com/chartmuseum/helm-push/issues/54
	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		_ = os.Setenv("HELM_PLUGINS", os.Getenv("HELM_PLUGIN"))
	}

	settings := cli.New()
	getters := getter.All(settings)
	cr, err := repo.NewChartRepository(entry, getters)
	if err != nil {
		return nil, err
	}

	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		home := v2helmHome()
		cr.CachePath = filepath.Join(home.Repository(), "cache")
	}

	return &Repo{cr}, nil
}

// GetRepoByName returns repository by name
func GetRepoByName(name string) (*Repo, error) {
	r, err := repoFile()
	if err != nil {
		return nil, err
	}
	entry, exists := findRepoEntry(name, r)
	if !exists {
		return nil, fmt.Errorf("no repo named %q found", name)
	}

	// FIXME
	// https://github.com/chartmuseum/helm-push/issues/54
	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		_ = os.Setenv("HELM_PLUGINS", os.Getenv("HELM_PLUGIN"))
	}

	settings := cli.New()
	getters := getter.All(settings)
	cr, err := repo.NewChartRepository(entry, getters)
	if err != nil {
		return nil, err
	}

	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		home := v2helmHome()
		cr.CachePath = filepath.Join(home.Repository(), "cache")
	}

	return &Repo{cr}, nil
}

// TempRepoFromURL builds a temporary Repo from a given URL
func TempRepoFromURL(url string) (*Repo, error) {
	u, err := urllib.Parse(url)
	if err != nil {
		return nil, err
	}
	entry := &repo.Entry{}
	if u.User != nil {
		// remove the username/password section from URL
		pass, _ := u.User.Password()
		entry.URL = strings.Split(url, "://")[0] + "://" + strings.Split(url, fmt.Sprintf("%s@", pass))[1]
		entry.Username = u.User.Username()
		entry.Password = pass
	} else {
		entry.URL = url
	}
	cr, err := repo.NewChartRepository(entry, getter.All(cli.New()))
	if err != nil {
		return nil, err
	}
	return &Repo{cr}, nil
}

func repoFile() (*repo.File, error) {
	var repoFilePath string
	if HelmMajorVersionCurrent() == HelmMajorVersion2 {
		home := v2helmHome()
		repoFilePath = home.RepositoryFile()
	} else {
		settings := cli.New()
		repoFilePath = settings.RepositoryConfig
	}
	repoFile, err := repo.LoadFile(repoFilePath)
	return repoFile, err
}

func v2helmHome() v2helmpath.Home {
	var helmHomePath string
	if v, ok := os.LookupEnv("HELM_HOME"); ok {
		helmHomePath = v
	} else {
		helmHomePath = v2environment.DefaultHelmHome
	}
	return v2helmpath.Home(helmHomePath)
}

func findRepoEntry(name string, r *repo.File) (*repo.Entry, bool) {
	var entry *repo.Entry
	exists := false
	for _, re := range r.Repositories {
		if re.Name == name {
			entry = re
			exists = true
			break
		}
	}
	return entry, exists
}

func findRepoEntryByURL(url string, r *repo.File) (*repo.Entry, bool) {
	var entry *repo.Entry
	exists := false
	for _, re := range r.Repositories {
		if strings.TrimSuffix(re.URL, "/") == strings.TrimSuffix(url, "/") {
			entry = re
			exists = true
			break
		}
	}
	return entry, exists
}
