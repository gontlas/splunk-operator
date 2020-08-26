// Copyright (c) 2018-2020 Splunk Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enterprise

import (
	"encoding/json"
	"fmt"
	"testing"

	enterprisev1 "github.com/splunk/splunk-operator/pkg/apis/enterprise/v1alpha3"
	splcommon "github.com/splunk/splunk-operator/pkg/splunk/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func configTester(t *testing.T, method string, f func() (interface{}, error), want string) {
	result, err := f()
	if err != nil {
		t.Errorf("%s returned error: %v", method, err)
	}
	got, err := json.Marshal(result)
	if err != nil {
		t.Errorf("%s failed to marshall: %v", method, err)
	}
	if string(got) != want {
		t.Errorf("%s = %s;\nwant %s", method, got, want)
	}
}

func TestGetSplunkService(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
	}

	test := func(instanceType InstanceType, isHeadless bool, want string) {
		f := func() (interface{}, error) {
			return getSplunkService(&cr, &cr.Spec.CommonSplunkSpec, instanceType, isHeadless), nil
		}
		configTester(t, fmt.Sprintf("getSplunkService(\"%s\",%t)", instanceType, isHeadless), f, want)
	}

	test(SplunkIndexer, false, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-indexer-service","namespace":"test","creationTimestamp":null,"labels":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"hec","protocol":"TCP","port":8088,"targetPort":8088},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"s2s","protocol":"TCP","port":9997,"targetPort":9997}],"selector":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"}},"status":{"loadBalancer":{}}}`)
	test(SplunkIndexer, true, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-indexer-headless","namespace":"test","creationTimestamp":null,"labels":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"hec","protocol":"TCP","port":8088,"targetPort":8088},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"s2s","protocol":"TCP","port":9997,"targetPort":9997}],"selector":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"},"clusterIP":"None","type":"ClusterIP"},"status":{"loadBalancer":{}}}`)
	// Multipart IndexerCluster - test part-of and instance labels for child part
	cr.Spec.IndexerClusterRef.Name = "cluster1"
	test(SplunkIndexer, false, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-indexer-service","namespace":"test","creationTimestamp":null,"labels":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/instance":"splunk-stack1-indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-cluster1-indexer"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"hec","protocol":"TCP","port":8088,"targetPort":8088},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"s2s","protocol":"TCP","port":9997,"targetPort":9997}],"selector":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/instance":"splunk-stack1-indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-cluster1-indexer"}},"status":{"loadBalancer":{}}}`)
	cr.Spec.IndexerClusterRef.Name = ""

	cr.Spec.ServiceTemplate.Spec.Type = "LoadBalancer"
	cr.Spec.ServiceTemplate.ObjectMeta.Labels = map[string]string{"1": "2"}
	cr.ObjectMeta.Labels = map[string]string{"one": "two"}
	cr.ObjectMeta.Annotations = map[string]string{"a": "b"}

	test(SplunkSearchHead, false, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-search-head-service","namespace":"test","creationTimestamp":null,"labels":{"1":"2","app.kubernetes.io/component":"search-head","app.kubernetes.io/instance":"splunk-stack1-search-head","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"search-head","app.kubernetes.io/part-of":"splunk-stack1-search-head","one":"two"},"annotations":{"a":"b"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"dfsmaster","protocol":"TCP","port":9000,"targetPort":9000},{"name":"dfccontrol","protocol":"TCP","port":17000,"targetPort":17000},{"name":"datareceive","protocol":"TCP","port":19000,"targetPort":19000}],"selector":{"app.kubernetes.io/component":"search-head","app.kubernetes.io/instance":"splunk-stack1-search-head","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"search-head","app.kubernetes.io/part-of":"splunk-stack1-search-head"},"type":"LoadBalancer"},"status":{"loadBalancer":{}}}`)
	test(SplunkSearchHead, true, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-search-head-headless","namespace":"test","creationTimestamp":null,"labels":{"app.kubernetes.io/component":"search-head","app.kubernetes.io/instance":"splunk-stack1-search-head","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"search-head","app.kubernetes.io/part-of":"splunk-stack1-search-head","one":"two"},"annotations":{"a":"b"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"dfsmaster","protocol":"TCP","port":9000,"targetPort":9000},{"name":"dfccontrol","protocol":"TCP","port":17000,"targetPort":17000},{"name":"datareceive","protocol":"TCP","port":19000,"targetPort":19000}],"selector":{"app.kubernetes.io/component":"search-head","app.kubernetes.io/instance":"splunk-stack1-search-head","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"search-head","app.kubernetes.io/part-of":"splunk-stack1-search-head"},"clusterIP":"None","type":"ClusterIP","publishNotReadyAddresses":true},"status":{"loadBalancer":{}}}`)
}

func TestGetSplunkDefaults(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
		Spec: enterprisev1.IndexerClusterSpec{
			CommonSplunkSpec: enterprisev1.CommonSplunkSpec{Defaults: "defaults_string"},
		},
	}

	test := func(want string) {
		f := func() (interface{}, error) {
			return getSplunkDefaults(cr.GetName(), cr.GetNamespace(), SplunkIndexer, cr.Spec.Defaults), nil
		}
		configTester(t, "getSplunkDefaults()", f, want)
	}

	test(`{"metadata":{"name":"splunk-stack1-indexer-defaults","namespace":"test","creationTimestamp":null},"data":{"default.yml":"defaults_string"}}`)
}

func TestGetService(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
		Spec: enterprisev1.IndexerClusterSpec{
			Replicas: 3,
			CommonSplunkSpec: enterprisev1.CommonSplunkSpec{
				Spec: splcommon.Spec{
					ServiceTemplate: corev1.Service{
						Spec: corev1.ServiceSpec{
							Ports: []corev1.ServicePort{{Name: "user-defined", Port: 32000, TargetPort: intstr.FromInt(6443)}},
						},
					},
				},
			},
		},
	}

	test := func(instanceType InstanceType, want string) {
		f := func() (interface{}, error) {
			return getSplunkService(&cr, &cr.Spec.CommonSplunkSpec, instanceType, false), nil
		}
		configTester(t, "getSplunkService()", f, want)
	}

	test(SplunkIndexer, `{"kind":"Service","apiVersion":"v1","metadata":{"name":"splunk-stack1-indexer-service","namespace":"test","creationTimestamp":null,"labels":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"},"ownerReferences":[{"apiVersion":"","kind":"","name":"stack1","uid":"","controller":true}]},"spec":{"ports":[{"name":"user-defined","port":32000,"targetPort":6443},{"name":"splunkweb","protocol":"TCP","port":8000,"targetPort":8000},{"name":"hec","protocol":"TCP","port":8088,"targetPort":8088},{"name":"splunkd","protocol":"TCP","port":8089,"targetPort":8089},{"name":"s2s","protocol":"TCP","port":9997,"targetPort":9997}],"selector":{"app.kubernetes.io/component":"indexer","app.kubernetes.io/managed-by":"splunk-operator","app.kubernetes.io/name":"indexer","app.kubernetes.io/part-of":"splunk-stack1-indexer"}},"status":{"loadBalancer":{}}}`)
}

func TestSetVolumeDefault(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
	}

	setVolumeDefaults(&cr.Spec.CommonSplunkSpec)
	if cr.Spec.CommonSplunkSpec.Volumes == nil {
		t.Errorf("setVolumeDefaults() returns nil for Volumes")
	}

	mode := int32(644)
	cr.Spec.CommonSplunkSpec.Volumes = []corev1.Volume{
		{
			Name: "vol1",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "top-secret1",
				},
			},
		},
		{
			Name: "vol2",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "top-secret2",
					DefaultMode: &mode,
				},
			},
		},
		{
			Name: "vol3",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "config3",
					},
				},
			},
		},
	}

	// Make sure the default mode is set correcty
	setVolumeDefaults(&cr.Spec.CommonSplunkSpec)
	if cr.Spec.CommonSplunkSpec.Volumes == nil {
		t.Errorf("setVolumeDefaults() returns nil for Volumes")
	}

	for _, v := range cr.Spec.CommonSplunkSpec.Volumes {
		if v.Name == "vol1" {
			if *v.Secret.DefaultMode != int32(corev1.SecretVolumeSourceDefaultMode) {
				t.Errorf("setVolumeDefaults() did not set defaultMode correctly. Want %d, Got %d", int32(corev1.SecretVolumeSourceDefaultMode), *v.Secret.DefaultMode)
			}
		} else if v.Name == "vol2" {
			if *v.Secret.DefaultMode != mode {
				t.Errorf("setVolumeDefaults() did not set defaultMode correctly. Want %d, Got %d", mode, *v.Secret.DefaultMode)
			}
		} else if v.Name == "vol3" {
			if *v.ConfigMap.DefaultMode != int32(corev1.ConfigMapVolumeSourceDefaultMode) {
				t.Errorf("setVolumeDefaults() did not set defaultMode correctly. Want %d, Got %d", int32(corev1.ConfigMapVolumeSourceDefaultMode), *v.ConfigMap.DefaultMode)
			}
		}
	}
}

func TestSmartstoreApplyIndexerClusterFailsOnInvalidSmartStoreConfig(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "idxCluster",
			Namespace: "test",
		},
		Spec: enterprisev1.IndexerClusterSpec{
			Replicas: 1,
			SmartStore: enterprisev1.SmartStoreSpec{
				VolList: []enterprisev1.VolumeSpec{
					{Name: "msos_s2s3_vol", Endpoint: ""},
				},

				IndexList: []enterprisev1.IndexSpec{
					{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
				},
			},
		},
	}

	var client splcommon.ControllerClient

	_, err := ApplyIndexerCluster(client, &cr)
	if err == nil {
		t.Errorf("ApplyIndexerCluster should fail on invalid smartstore config")
	}
}

func TestSmartstoreApplyStandaloneFailsOnInvalidSmartStoreConfig(t *testing.T) {
	cr := enterprisev1.Standalone{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "standalone",
			Namespace: "test",
		},
		Spec: enterprisev1.StandaloneSpec{
			Replicas: 1,
			SmartStore: enterprisev1.SmartStoreSpec{
				VolList: []enterprisev1.VolumeSpec{
					{Name: "msos_s2s3_vol", Endpoint: ""},
				},

				IndexList: []enterprisev1.IndexSpec{
					{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
				},
			},
		},
	}

	var client splcommon.ControllerClient

	_, err := ApplyStandalone(client, &cr)
	if err == nil {
		t.Errorf("ApplyStandalone should fail on invalid smartstore config")
	}
}

func TestSmartStoreConfigDoesNotFailOnIndexerClusterCRForCM(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "idxc_CM",
			Namespace: "test",
		},
		Spec: enterprisev1.IndexerClusterSpec{
			Replicas: 3,
			SmartStore: enterprisev1.SmartStoreSpec{
				VolList: []enterprisev1.VolumeSpec{
					{Name: "msos_s2s3_vol", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
				},

				IndexList: []enterprisev1.IndexSpec{
					{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
				},
			},
		},
	}

	err := validateIndexerClusterSpec(&cr)

	if err != nil {
		t.Errorf("Smartstore configuration should not fail on IndexerCluster CR with CM: %v", err)
	}
}

func TestSmartStoreConfigFailsOnIndexerClusterCRForIndexers(t *testing.T) {
	cr := enterprisev1.IndexerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "idxc",
			Namespace: "test",
		},
		Spec: enterprisev1.IndexerClusterSpec{
			Replicas: 3,
			SmartStore: enterprisev1.SmartStoreSpec{
				VolList: []enterprisev1.VolumeSpec{
					{Name: "msos_s2s3_vol", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
				},

				IndexList: []enterprisev1.IndexSpec{
					{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
					{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
				},
			},
		},
	}

	cr.Spec.IndexerClusterRef.Name = "testRefWithCM"

	err := validateIndexerClusterSpec(&cr)

	if err == nil {
		t.Errorf("Indexer Cluster Custom Resource for indexers should not allow Smartstore configuration")
	}
}

func TestValidateSplunkSmartstoreSpec(t *testing.T) {
	var err error

	// Valid smartstore config
	SmartStore := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "msos_s2s3_vol_2", Endpoint: "https://s3-us-west-2.amazonaws.com"},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	err = ValidateSplunkSmartstoreSpec(&SmartStore)
	if err != nil {
		t.Errorf("Valid Smartstore configuration should not cause error: %v", err)
	}

	// Only one remote volume is allowed
	SmartStoreMultipleVolumes := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "msos_s2s3_vol_1", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
			{Name: "msos_s2s3_vol_2", Endpoint: "https://s3-us-west-2.amazonaws.com"},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	err = ValidateSplunkSmartstoreSpec(&SmartStoreMultipleVolumes)
	if err == nil {
		t.Errorf("Multiple Smartstore volume configurations should error out")
	}

	// Smartstore config with missing endpoint for the volume
	SmartStoreVolumeWithNoRemoteEndPoint := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "msos_s2s3_vol_1", Endpoint: ""},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	err = ValidateSplunkSmartstoreSpec(&SmartStoreVolumeWithNoRemoteEndPoint)
	if err == nil {
		t.Errorf("Should not accept a volume with missing Endpoint")
	}

	// Smartstore config with missing remote name for the volume
	SmartStoreWithVolumeNameMissing := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	err = ValidateSplunkSmartstoreSpec(&SmartStoreWithVolumeNameMissing)
	if err == nil {
		t.Errorf("Should not accept a volume with missing Remotename")
	}

	// Smartstore config with missing index name
	SmartStoreWithMissingIndexName := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "msos_s2s3_vol_1", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	//Smartstore config with missing remote index location
	err = ValidateSplunkSmartstoreSpec(&SmartStoreWithMissingIndexName)
	if err == nil {
		t.Errorf("Should not accept an Index with missing indexname ")
	}

	SmartStoreWithMissingIndexLocation := enterprisev1.SmartStoreSpec{
		VolList: []enterprisev1.VolumeSpec{
			{Name: "msos_s2s3_vol_1", Endpoint: "https://s3-eu-west-2.amazonaws.com"},
		},

		IndexList: []enterprisev1.IndexSpec{
			{Name: "salesdata1", RemoteLocation: "testbucket-rs-london/$_index_name"},
			{Name: "salesdata2", RemoteLocation: ""},
			{Name: "salesdata3", RemoteLocation: "testbucket-rs-london/$_index_name"},
		},
	}

	err = ValidateSplunkSmartstoreSpec(&SmartStoreWithMissingIndexLocation)
	if err == nil {
		t.Errorf("Should not accept an Index with missing index location")
	}

	// Empty smartstore config
	err = ValidateSplunkSmartstoreSpec(nil)
	if err != nil {
		t.Errorf("Smartstore config is optional, should not cause an error")
	}
}
