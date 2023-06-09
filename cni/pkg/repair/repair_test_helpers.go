// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repair

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/istio/tools/istio-iptables/pkg/constants"
)

type makePodArgs struct {
	PodName             string
	Namespace           string
	Labels              map[string]string
	Annotations         map[string]string
	InitContainerName   string
	InitContainerStatus *v1.ContainerStatus
	NodeName            string
}

func makePod(args makePodArgs) *v1.Pod {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        args.PodName,
			Namespace:   args.Namespace,
			Labels:      args.Labels,
			Annotations: args.Annotations,
		},
		Spec: v1.PodSpec{
			NodeName: args.NodeName,
			Volumes:  nil,
			InitContainers: []v1.Container{
				{
					Name: args.InitContainerName,
				},
			},
			Containers: []v1.Container{
				{
					Name: "payload-container",
				},
			},
		},
		Status: v1.PodStatus{
			InitContainerStatuses: []v1.ContainerStatus{
				*args.InitContainerStatus,
			},
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name: "payload-container",
					State: v1.ContainerState{
						Waiting: &v1.ContainerStateWaiting{
							Reason: "PodInitializing",
						},
					},
				},
			},
		},
	}
	return pod
}

// Container specs
var (
	brokenInitContainerWaiting = v1.ContainerStatus{
		Name: constants.ValidationContainerName,
		State: v1.ContainerState{
			Waiting: &v1.ContainerStateWaiting{
				Reason:  "CrashLoopBackOff",
				Message: "Back-off 5m0s restarting failed blah blah blah",
			},
		},
		LastTerminationState: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: constants.ValidationErrorCode,
				Reason:   "Error",
				Message:  "Died for some reason",
			},
		},
	}

	brokenInitContainerTerminating = v1.ContainerStatus{
		Name: constants.ValidationContainerName,
		State: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: constants.ValidationErrorCode,
				Reason:   "Error",
				Message:  "Died for some reason",
			},
		},
		LastTerminationState: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: constants.ValidationErrorCode,
				Reason:   "Error",
				Message:  "Died for some reason",
			},
		},
	}

	workingInitContainerDiedPreviously = v1.ContainerStatus{
		Name: constants.ValidationContainerName,
		State: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: 0,
				Reason:   "Completed",
			},
		},
		LastTerminationState: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: 126,
				Reason:   "Error",
				Message:  "Died for some reason",
			},
		},
	}

	workingInitContainer = v1.ContainerStatus{
		Name: constants.ValidationContainerName,
		State: v1.ContainerState{
			Terminated: &v1.ContainerStateTerminated{
				ExitCode: 0,
				Reason:   "Completed",
			},
		},
	}
)

// Pod specs
var (
	brokenPodTerminating = *makePod(makePodArgs{
		PodName: "BrokenPodTerminating",
		Annotations: map[string]string{
			"sidecar.istio.io/status": "something",
		},
		Labels: map[string]string{
			"testlabel": "true",
		},
		NodeName:            "TestNode",
		InitContainerStatus: &brokenInitContainerTerminating,
	})

	brokenPodWaiting = *makePod(makePodArgs{
		PodName: "BrokenPodWaiting",
		Annotations: map[string]string{
			"sidecar.istio.io/status": "something",
		},
		NodeName:            "TestNode",
		InitContainerStatus: &brokenInitContainerWaiting,
	})

	brokenPodNoAnnotation = *makePod(makePodArgs{
		PodName:             "BrokenPodNoAnnotation",
		InitContainerStatus: &brokenInitContainerWaiting,
	})

	workingPod = *makePod(makePodArgs{
		PodName: "WorkingPod",
		Annotations: map[string]string{
			"sidecar.istio.io/status": "something",
		},
		InitContainerStatus: &workingInitContainer,
	})

	workingPodDiedPreviously = *makePod(makePodArgs{
		PodName: "WorkingPodDiedPreviously",
		Annotations: map[string]string{
			"sidecar.istio.io/status": "something",
		},
		InitContainerStatus: &workingInitContainerDiedPreviously,
	})
)
