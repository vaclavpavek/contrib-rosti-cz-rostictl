package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rosti-cz/cli/src/parser"
	"github.com/rosti-cz/cli/src/rostiapi"
	"github.com/rosti-cz/cli/src/state"
	"github.com/urfave/cli/v2"
)

func createArchive(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

// Returns list of SSH key found in the current system.
// It looks for the keys in ~/.ssh which should be valid for Linux, Mac and possibly Windows.
// The function returns paths to private key, public key and error
func findSSHKey() (string, string, error) {
	keyFileNames := []string{
		// "id_ed25519", // Not supported by current dropbear version
		"id_rsa",
	}

	user, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("getting user info error: %w", err)
	}

	for _, keyFilename := range keyFileNames {
		privateKeyPath := path.Join(user.HomeDir, ".ssh", keyFilename)
		publicKeyPath := path.Join(user.HomeDir, ".ssh", keyFilename+".pub")

		_, errPrivate := os.Stat(privateKeyPath)
		_, errPublic := os.Stat(publicKeyPath)

		if !os.IsNotExist(errPrivate) && !os.IsNotExist(errPublic) {
			return privateKeyPath, publicKeyPath, nil
		}
	}

	return "", "", errors.New("no ssh key found")
}

func readLocalSSHPubKey() (string, error) {
	_, publicKeyPath, err := findSSHKey()
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// findCompany returns company ID based the environment
func findCompany(client *rostiapi.Client, appState *state.RostiState, c *cli.Context) (uint, error) {
	companies, err := client.GetCompanies()
	if err != nil {
		return 0, err
	}

	if len(companies) == 0 {
		return 0, errors.New("no company found")
	}

	companyIDFromFlag := uint(c.Int("company"))
	companyID := appState.CompanyID

	if companyIDFromFlag != 0 {
		companyID = companyIDFromFlag
	} else if companyID == 0 {
		if len(companies) == 1 {
			companyID = companies[0].ID
		} else if len(companies) > 1 {
			fmt.Println("You have access to multiple companies, pick one of the companies below and use -c COMPANY_ID flag to call this command.")
			fmt.Println("")
			fmt.Printf("  %6s  Company name\n", "ID")
			fmt.Printf("  %6s  ------------\n", "------")
			for _, company := range companies {
				fmt.Printf("  %6s  %s\n", strconv.Itoa(int(company.ID)), company.Name)
			}
			fmt.Println("")
			return companyID, nil
		} else {
			return companyID, errors.New("no company found")
		}
	}

	var found bool
	for _, company := range companies {
		if company.ID == companyID {
			found = true
			break
		}
	}
	if !found {
		return companyID, errors.New("selected company (" + strconv.Itoa(int(companyIDFromFlag)) + ") not found")
	}

	return companyID, nil
}

// Selects plan based on Rostifile or default settings
func selectPlan(client *rostiapi.Client, rostifile *parser.Rostifile) (uint, error) {
	// TODO: implements something like default plan loaded from the API (needs support in the API)
	rostifile.Plan = "start"

	fmt.Println(".. loading list of available plans")
	plans, err := client.GetPlans()
	if err != nil {
		return 0, err
	}

	var planID uint
	for _, plan := range plans {
		if strings.ToLower(plan.Name) == strings.ToLower(rostifile.Plan) {
			planID = plan.ID
		}
	}

	return planID, nil
}

// Selects runtime image based on rostifile
func selectRuntime(client *rostiapi.Client, rostifile *parser.Rostifile) (string, error) {
	fmt.Println(".. loading list of available runtimes")
	runtimes, err := client.GetRuntimes()
	if err != nil {
		return "", err
	}

	var selectedRuntime string
	var lastRuntime string

	if len(runtimes) == 0 {
		return selectedRuntime, errors.New("no runtime available")
	}

	for _, runtime := range runtimes {
		if runtime.Image == rostifile.Runtime {
			selectedRuntime = rostifile.Runtime
			break
		}
		lastRuntime = runtime.Image
	}

	if selectedRuntime == "" {
		selectedRuntime = lastRuntime
	}

	return selectedRuntime, nil
}