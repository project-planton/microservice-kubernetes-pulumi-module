package data

import (
	"fmt"
	"reflect"
	"testing"

	code2cloudv1deploymsistackk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
)

func TestSecretStringData(t *testing.T) {
	testCases := []struct {
		name           string
		input          *code2cloudv1deploymsistackk8smodel.KubernetesImagePullSecretInput
		expectedOutput map[string]string
		expectedErr    error
	}{{
		name: "valid input",
		input: &code2cloudv1deploymsistackk8smodel.KubernetesImagePullSecretInput{
			ArtifactReaderGsaKeyBase64: "c29tZWdpYmJlcmlzaAo=",
			DockerRepoHostname:         "us-central1-docker.pkg.dev",
		},
		expectedOutput: map[string]string{Key: fmt.Sprintf(DockerConfigTemplateFormatString, "us-central1-docker.pkg.dev", "X2pzb25fa2V5OnNvbWVnaWJiZXJpc2gK")},
		expectedErr:    nil,
	}}
	t.Run("docker config secret data", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := Get(tc.input)
				if err != tc.expectedErr {
					t.Errorf("expected: error %v got: %v", tc.expectedErr, err)
				}
				if !reflect.DeepEqual(result, tc.expectedOutput) {
					t.Errorf("expected: %v got: %v", tc.expectedOutput, result)
				}
			})
		}
	})
}

func TestDockerConfigAuth(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedErr    error
	}{{
		name:           "valid",
		input:          "c29tZWdpYmJlcmlzaAo=",
		expectedOutput: "X2pzb25fa2V5OnNvbWVnaWJiZXJpc2gK",
		expectedErr:    nil,
	}}
	t.Run("docker config auth", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := getDockerConfigAuth(tc.input)
				if err != tc.expectedErr {
					t.Errorf("expected: error %v got: %v", tc.expectedErr, err)
				}
				if !reflect.DeepEqual(result, tc.expectedOutput) {
					t.Errorf("expected: %v got: %v", tc.expectedOutput, result)
				}
			})
		}
	})
}
