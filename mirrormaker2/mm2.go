package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

// RVER is set during build
var RVER string = "<not set>"

var logger *logrus.Logger

// MirrorMakerEnv map for the mm2.properties file
type MirrorMakerEnv map[string]string

// ConfigProps for populating mm2.properties
type ConfigProps struct {
	Name  string
	Value string
}

func readPropertyFile(file string) (MirrorMakerEnv, error) {
	mirrorMakerEnv := MirrorMakerEnv{}
	if file == "" {
		return mirrorMakerEnv, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if line == "\n" || strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		eqIndex := strings.Index(line, "=")
		if eqIndex <= 0 {
			return nil, fmt.Errorf("line [%s] has an invalid format", strings.TrimSpace(line))
		}
		key := strings.TrimSpace(line[:eqIndex])
		if len(line) <= eqIndex {
			continue
		}
		if len(strings.TrimSpace(line[eqIndex+1:])) < 1 {
			mirrorMakerEnv[key] = ""
			continue
		}
		mirrorMakerEnv[key] = strings.TrimSpace(line[eqIndex+1:])

	}
	return mirrorMakerEnv, nil
}

// GenerateConfigFile for parsing env and populating mm2.properties
func GenerateConfigFile(file string, mirrorMakerEnv map[string]string) error {
	opt := ConfigProps{}
	const mirrorMakerEnvPropsTempl = `{{.Name}}={{.Value}}
`
	t := template.Must(template.New("ConfigProps").Parse(mirrorMakerEnvPropsTempl))
	outLogFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outLogFile.Close()

	for n, v := range mirrorMakerEnv {
		opt.Name = n
		opt.Value = v

		err := t.Execute(outLogFile, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupKafkaMirrorMaker Get mm2 configuration from env
func SetupKafkaMirrorMaker(source, target string) error {
	var key string

	// --------- LOADING ENV ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka MirrorMaker",
		"Stage":     "Environment",
	}).Info("Loading and parsing environment")

	mirrorMakerEnv, err := readPropertyFile(source)
	if err != nil {
		return err
	}
	osEnviron := os.Environ()
	kmm2Prefix := "KMM2_"

	for _, b := range osEnviron {
		if strings.HasPrefix(b, kmm2Prefix) {
			pair := strings.SplitN(b, "=", 2)
			key = strings.TrimPrefix(pair[0], kmm2Prefix)
			key = strings.ToLower(key)
			key = strings.Replace(key, "_", ".", -1)
			mirrorMakerEnv[key] = pair[1]
		}
	}

	err = GenerateConfigFile(
		target,
		mirrorMakerEnv,
	)
	if err != nil {
		return err
	}

	// --------- End Configuring System Items ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka MirrorMaker",
		"Stage":     "-",
	}).Info("Successfully configured kafka MirrorMaker")

	return nil
}

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	source := flag.String("source", "", "path to mm2.properties to read")
	target := flag.String("target", "/etc/kafka/kafka-mm2.properties", "mm2.properties files to be created")

	flag.Parse()

	if *target == "" {
		flag.Usage()
		return
	}

	logger = logrus.New()
	logger.WithField("Kafka MirrorMaker2", RVER).Infof("Initiated")

	cfgDir := filepath.Dir(*target)
	if cfgDir != "/" {
		cfgDir = strings.TrimSuffix(cfgDir, "/")
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			err = os.MkdirAll(cfgDir, 0700)
			if err != nil {
				logger.Fatal(err)
			}
		}
	}

	fileExists := false
	f, err := os.Stat(*source)
	if os.IsNotExist(err) {
		fileExists = false
	} else if f.IsDir() {
		fileExists = false
	}
	if !fileExists {
		logger.Info("Could not locate file [%s]")
	}

	err = SetupKafkaMirrorMaker(*source, *target)
	if err != nil {
		logger.Fatal(err)
	}
}
