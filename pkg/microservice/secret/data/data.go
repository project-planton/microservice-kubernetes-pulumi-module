package data

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	code2cloudv1deploymsistackk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
)

const (
	Key = ".dockerconfigjson"
	//DockerConfigTemplateFormatString takes following inputs in the same order
	// 1. docker registry hostname
	// 2. base64 encoded docker config auth
	DockerConfigTemplateFormatString = `
{
  "auths": {
    "%s": {
      "username": "_json_key",
      "auth": "%s"
    }
  }
}
`
)

func Get(input *code2cloudv1deploymsistackk8smodel.KubernetesImagePullSecretInput) (map[string]string, error) {
	dockerConfigAuth, err := getDockerConfigAuth(input.ArtifactReaderGsaKeyBase64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get docker auth config")
	}
	return map[string]string{Key: fmt.Sprintf(DockerConfigTemplateFormatString, input.DockerRepoHostname, dockerConfigAuth)}, nil
}

// getDockerConfigAuth creates base64 encoded docker config auth
func getDockerConfigAuth(gsaKeyBase64Encoded string) (string, error) {
	decodedStringBytes, err := b64.StdEncoding.DecodeString(gsaKeyBase64Encoded)
	if err != nil {
		return "", errors.Wrap(err, "failed to base64 decode gsa key")
	}
	dockerConfigAuth := fmt.Sprintf("_json_key:%s", string(decodedStringBytes))
	return b64.StdEncoding.EncodeToString([]byte(dockerConfigAuth)), nil
}
