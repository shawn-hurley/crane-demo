package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"github.com/sirupsen/logrus"
)

var logger logrus.FieldLogger

const Version = "v0.0.1"
const (
	// flags
	DeploymentVolumePath = "volume-data-path"

	volumePatch      = `[{"op": "add", "path": "/spec/template/spec/volumes/-", "value": {"name": "%v", "emptyDir": {}}}]`
	volumeMountPatch = `[{"op": "add", "path": "/spec/template/spec/containers/0/volumeMounts/-", "value": {"name": "%v", "mountPath": "%v"}}]`
)

func main() {
	fields := []transform.OptionalFields{
		{
			FlagName: DeploymentVolumePath,
			Help:     "the path that will be used in for the mount",
			Example:  `carts-db="/data"`,
		},
	}
	logger = logrus.New()

	cli.RunAndExit(cli.NewCustomPlugin("VolumeMountPlugin", Version, fields, Run))
}

func Run(request transform.PluginRequest) (transform.PluginResponse, error) {
	u := request.Unstructured
	response := transform.PluginResponse{
		Version:    string(transform.V1),
		IsWhiteOut: false,
		Patches:    jsonpatch.Patch{},
	}

	//	deploymentVolumeDataPath := transform.ParseOptionalFieldMapVal(request.Extras[DeploymentVolumePath])
	deploymentVolumeDataPath := map[string]string{
		"carts-db":  "/data/db",
		"orders-db": "/data/db",
		"user-db":   "/data/db-users",
	}

	if u.GetKind() != "Deployment" {
		return response, nil
	}

	var dataPath string
	var ok bool
	if dataPath, ok = deploymentVolumeDataPath[u.GetName()]; !ok {
		return response, nil
	}

	vp := fmt.Sprintf(volumePatch, "data-volume")
	vpatch, err := jsonpatch.DecodePatch([]byte(vp))
	if err != nil {
		return response, err
	}
	vmp := fmt.Sprintf(volumeMountPatch, "data-volume", dataPath)
	vmpatch, err := jsonpatch.DecodePatch([]byte(vmp))
	if err != nil {
		fmt.Printf("here")
		return response, err
	}

	vpatch = append(vpatch, vmpatch...)

	response.Patches = vpatch

	return response, nil
}
