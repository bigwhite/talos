/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package userdata

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
)

// Kubeadm describes the set of configuration options available for kubeadm.
type Kubeadm struct {
	CommonServiceOptions `yaml:",inline"`

	// ConfigurationStr is converted to Configuration and back in Marshal/UnmarshalYAML
	Configuration    runtime.Object `yaml:"-"`
	ConfigurationStr string         `yaml:"configuration"`

	ExtraArgs             []string `yaml:"extraArgs,omitempty"`
	CertificateKey        string   `yaml:"certificateKey,omitempty"`
	IgnorePreflightErrors []string `yaml:"ignorePreflightErrors,omitempty"`
	bootstrap             bool
	controlPlane          bool
}

// MarshalYAML implements the yaml.Marshaler interface.
func (kdm *Kubeadm) MarshalYAML() (interface{}, error) {
	b, err := configutil.MarshalKubeadmConfigObject(kdm.Configuration)
	if err != nil {
		return nil, err
	}

	gvks, err := kubeadmutil.GroupVersionKindsFromBytes(b)
	if err != nil {
		return nil, err
	}

	if kubeadmutil.GroupVersionKindsHasInitConfiguration(gvks...) {
		kdm.bootstrap = true
	}
	if kubeadmutil.GroupVersionKindsHasJoinConfiguration(gvks...) {
		kdm.bootstrap = false
	}

	kdm.ConfigurationStr = string(b)

	type KubeadmAlias Kubeadm

	return (*KubeadmAlias)(kdm), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (kdm *Kubeadm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type KubeadmAlias Kubeadm

	if err := unmarshal((*KubeadmAlias)(kdm)); err != nil {
		return err
	}

	b := []byte(kdm.ConfigurationStr)

	gvks, err := kubeadmutil.GroupVersionKindsFromBytes(b)
	if err != nil {
		return err
	}

	if kubeadmutil.GroupVersionKindsHasInitConfiguration(gvks...) {
		// Since the ClusterConfiguration is embedded in the InitConfiguration
		// struct, it is required to (un)marshal it a special way. The kubeadm
		// API exposes one function (MarshalKubeadmConfigObject) to handle the
		// marshaling, but does not yet have that convenience for
		// unmarshaling.
		cfg, err := configutil.BytesToInitConfiguration(b)
		if err != nil {
			return err
		}
		kdm.Configuration = cfg
		kdm.bootstrap = true
	}
	if kubeadmutil.GroupVersionKindsHasJoinConfiguration(gvks...) {
		cfg, err := kubeadmutil.UnmarshalFromYamlForCodecs(b, kubeadmapi.SchemeGroupVersion, kubeadmscheme.Codecs)
		if err != nil {
			return err
		}
		kdm.Configuration = cfg
		kdm.bootstrap = false
		joinConfiguration, ok := cfg.(*kubeadm.JoinConfiguration)
		if !ok {
			return errors.New("expected JoinConfiguration")
		}
		if joinConfiguration.ControlPlane == nil {
			kdm.controlPlane = false
		} else {
			kdm.controlPlane = true
		}
	}

	return nil
}
