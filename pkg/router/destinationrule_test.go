package router

import (
	v1alpha32 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"testing"
)

var fakeIstioClient IstioClientInterface

func TestMain(m *testing.M) {
	// discard stdout logs if not being run with '-v' flag
	log.SetOutput(ioutil.Discard)
	result := m.Run()
	os.Exit(result)
}

func TestValidateDestinationRuleList_Unit(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: []v1alpha32.VirtualService{
				{},
			},
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: []v1alpha32.DestinationRule{
				{},
			},
		},
	}

	err := ValidateDestinationRuleList(&irl)
	assert.NoError(t, err)
}

func TestValidateDestinationRuleList_Unit_EmptyItems(t *testing.T) {
	irl := IstioRouteList{
		VList: &v1alpha32.VirtualServiceList{
			Items: nil,
		},
		DList: &v1alpha32.DestinationRuleList{
			Items: nil,
		},
	}

	err := ValidateDestinationRuleList(&irl)
	assert.EqualError(t, err, "empty destinationRules")
}

func TestDestinationRule_Validate_Unit(t *testing.T) {
	fakeIstioClient = &fake.Clientset{}

	cases := []struct {
		dr    DestinationRule
		shift Shift
		want  string
	}{
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: nil,
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty label-selector",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     0,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty port",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     1000,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     66000,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"port not in range 1024 - 65535",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-testing",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic:  Traffic{},
			},
			"empty pod selector",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "",
			Namespace:  "arrow",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'name' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "",
			Build:      1,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'namespace' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "arrow",
			Build:      0,
			Istio:      fakeIstioClient,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"empty 'build' attribute",
		},
		{DestinationRule{
			TrackingId: "unit-testing-uuid",
			Name:       "api-test",
			Namespace:  "arrow",
			Build:      1,
			Istio:      nil,
		},
			Shift{
				Port:     8080,
				Hostname: "api-domain",
				Selector: map[string]string{"app": "api-domain"},
				Traffic: Traffic{
					PodSelector: map[string]string{"version": "1.2.3"},
				},
			},
			"nil istioClient object",
		},
	}

	for _, tt := range cases {
		err := tt.dr.Validate(tt.shift)
		assert.EqualError(t, err, tt.want)
	}
}

func TestDestinationRule_Create_Integrated(t *testing.T) {
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Name:       "api-testing",
		Namespace:  "arrow",
		Build:      10000,
		Istio:      fakeIstioClient,
	}

	shift := Shift{
		Traffic: Traffic{
			PodSelector: map[string]string{
				"environment": "test",
				"app":         "api-testing",
			},
		},
	}

	irl, err := dr.Create(shift)
	assert.NotNil(t, irl)
	assert.NoError(t, err)
	assert.Equal(t, "api-testing-10000-arrow", irl.Subset.Name)
}

func TestDestinationRule_Clear(t *testing.T) {
	fakeIstioClient = &fake.Clientset{}
	dr := DestinationRule{
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	shift := Shift{}

	err := dr.Clear(shift)
	assert.NoError(t, err)
}

func TestDestinationRule_Update_Integrated(t *testing.T) {
	fakeIstioClient = fake.NewSimpleClientset()
	dr := DestinationRule{
		Namespace:  "integration",
		TrackingId: "unit-testing-tracking-id",
		Istio:      fakeIstioClient,
	}

	// create a destinationRule object in memory
	tdr := v1alpha32.DestinationRule{
		Spec: v1alpha32.DestinationRuleSpec{},
	}

	tdr.Name = "integration-testing-dr"
	tdr.Namespace = dr.Namespace
	labelSelector := map[string]string{
		"app":         "api-test",
		"environment": "integration-tests",
	}
	tdr.Labels = labelSelector

	_, err := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).Create(&tdr)

	shift := Shift{
		Port:     8080,
		Hostname: "api-domain",
		Selector: labelSelector,
		Traffic: Traffic{
			PodSelector: map[string]string{"version": "1.2.3"},
		},
	}

	err = dr.Update(shift)

	v, _ := fakeIstioClient.NetworkingV1alpha3().DestinationRules(dr.Namespace).List(v1.ListOptions{})
	mockedDr := v.Items[0]

	assert.NoError(t, err)
	assert.Equal(t, 1, len(v.Items))
	assert.Equal(t, "integration-testing-dr", mockedDr.Name)
	assert.Equal(t, "integration", mockedDr.Namespace)
	assert.Equal(t, "integration-tests", mockedDr.Labels["environment"])
}
