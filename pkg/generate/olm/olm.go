package olm

import (
	"fmt"
	"strings"
	ttpl "text/template"
	"time"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	opsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
)

var (
	csvTemplatePaths = []string{
		"config/controller/user-env.yaml.tpl",
		"config/samples/sample.yaml.tpl",
	}

	csvIncludePaths = []string{
		"config/controller/kustomization_def.yaml.tpl",
	}

	csvCopyPaths = []string{
		"config/manifests/kustomization.yaml",
		"config/scorecard/bases/config.yaml",
		"config/scorecard/kustomization.yaml",
		"config/scorecard/patches/basic.config.yaml",
		"config/scorecard/patches/olm.config.yaml",
		"config/samples/kustomization.yaml",
	}

	csvFuncMap = ttpl.FuncMap{
		"ToLower": strings.ToLower,
	}
)

// BundleAssets generates the assets necessary to generate
// a bundle used for deploying a service via OLM.
func BundleAssets(
	g *generate.Generator,
	commonMeta CommonMetadata,
	serviceConfig ServiceConfig,
	vers string,
	templateBasePath string,
) (*templateset.TemplateSet, error) {

	ts := templateset.New(
		templateBasePath,
		csvIncludePaths,
		csvCopyPaths,
		csvFuncMap,
	)

	crds, err := g.GetCRDs()
	if err != nil {
		return nil, err
	}

	olmVars := templateOLMVars{
		vers,
		time.Now().Format("2006-01-02 15:04:05"),
		g.MetaVars(),
		commonMeta,
		serviceConfig,
		crds,
	}

	for _, path := range csvTemplatePaths {
		outPath := strings.TrimSuffix(path, ".tpl")
		if err := ts.Add(outPath, path, olmVars); err != nil {
			return nil, err
		}
	}

	if err := ts.Add("config/controller/kustomization.yaml", "config/controller/olm-kustomization.yaml.tpl", olmVars); err != nil {
		return nil, err
	}

	csvBaseOutPath := fmt.Sprintf(
		"config/manifests/bases/ack-%s-controller.clusterserviceversion.yaml",
		g.MetaVars().ServiceIDClean)
	if err := ts.Add(csvBaseOutPath, "config/manifests/bases/clusterserviceversion.yaml.tpl", olmVars); err != nil {
		return nil, err
	}

	return ts, nil
}

type templateOLMVars struct {
	Version   string
	CreatedAt string
	templateset.MetaVars
	Common CommonMetadata
	ServiceConfig
	CRDs []*ackmodel.CRD
}

// DefaultServiceConfig returns a default representation of ServiceConfig to be
// used as a base. The various values are expected to be overridden (if
// needed) by an input YAML manifest for the given service controller.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Annotations{
			Repository:         "https://github.com/aws-controllers-k8s",
			SuggestedNamespace: "ack-system",
			ShortDescription:   "This is the placeholder short description for the default configuration",
			IsCertified:        false,
			ContainerImage:     "public.ecr.aws/aws-controllers-k8s/controller",
			Support:            "Community",
		},
		[]Sample{},
		opsv1alpha1.ClusterServiceVersionSpec{
			Maturity: "alpha",
			Icon: []opsv1alpha1.Icon{
				{
					// This is the AWS logo
					Data:      "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCEtLSBHZW5lcmF0b3I6IEFkb2JlIElsbHVzdHJhdG9yIDE5LjAuMSwgU1ZHIEV4cG9ydCBQbHVnLUluIC4gU1ZHIFZlcnNpb246IDYuMDAgQnVpbGQgMCkgIC0tPgo8c3ZnIHZlcnNpb249IjEuMSIgaWQ9IkxheWVyXzEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMzA0IDE4MiIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgMzA0IDE4MjsiIHhtbDpzcGFjZT0icHJlc2VydmUiPgo8c3R5bGUgdHlwZT0idGV4dC9jc3MiPgoJLnN0MHtmaWxsOiMyNTJGM0U7fQoJLnN0MXtmaWxsLXJ1bGU6ZXZlbm9kZDtjbGlwLXJ1bGU6ZXZlbm9kZDtmaWxsOiNGRjk5MDA7fQo8L3N0eWxlPgo8Zz4KCTxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik04Ni40LDY2LjRjMCwzLjcsMC40LDYuNywxLjEsOC45YzAuOCwyLjIsMS44LDQuNiwzLjIsNy4yYzAuNSwwLjgsMC43LDEuNiwwLjcsMi4zYzAsMS0wLjYsMi0xLjksM2wtNi4zLDQuMiAgIGMtMC45LDAuNi0xLjgsMC45LTIuNiwwLjljLTEsMC0yLTAuNS0zLTEuNEM3Ni4yLDkwLDc1LDg4LjQsNzQsODYuOGMtMS0xLjctMi0zLjYtMy4xLTUuOWMtNy44LDkuMi0xNy42LDEzLjgtMjkuNCwxMy44ICAgYy04LjQsMC0xNS4xLTIuNC0yMC03LjJjLTQuOS00LjgtNy40LTExLjItNy40LTE5LjJjMC04LjUsMy0xNS40LDkuMS0yMC42YzYuMS01LjIsMTQuMi03LjgsMjQuNS03LjhjMy40LDAsNi45LDAuMywxMC42LDAuOCAgIGMzLjcsMC41LDcuNSwxLjMsMTEuNSwyLjJ2LTcuM2MwLTcuNi0xLjYtMTIuOS00LjctMTZjLTMuMi0zLjEtOC42LTQuNi0xNi4zLTQuNmMtMy41LDAtNy4xLDAuNC0xMC44LDEuM2MtMy43LDAuOS03LjMsMi0xMC44LDMuNCAgIGMtMS42LDAuNy0yLjgsMS4xLTMuNSwxLjNjLTAuNywwLjItMS4yLDAuMy0xLjYsMC4zYy0xLjQsMC0yLjEtMS0yLjEtMy4xdi00LjljMC0xLjYsMC4yLTIuOCwwLjctMy41YzAuNS0wLjcsMS40LTEuNCwyLjgtMi4xICAgYzMuNS0xLjgsNy43LTMuMywxMi42LTQuNWM0LjktMS4zLDEwLjEtMS45LDE1LjYtMS45YzExLjksMCwyMC42LDIuNywyNi4yLDguMWM1LjUsNS40LDguMywxMy42LDguMywyNC42VjY2LjR6IE00NS44LDgxLjYgICBjMy4zLDAsNi43LTAuNiwxMC4zLTEuOGMzLjYtMS4yLDYuOC0zLjQsOS41LTYuNGMxLjYtMS45LDIuOC00LDMuNC02LjRjMC42LTIuNCwxLTUuMywxLTguN3YtNC4yYy0yLjktMC43LTYtMS4zLTkuMi0xLjcgICBjLTMuMi0wLjQtNi4zLTAuNi05LjQtMC42Yy02LjcsMC0xMS42LDEuMy0xNC45LDRjLTMuMywyLjctNC45LDYuNS00LjksMTEuNWMwLDQuNywxLjIsOC4yLDMuNywxMC42ICAgQzM3LjcsODAuNCw0MS4yLDgxLjYsNDUuOCw4MS42eiBNMTI2LjEsOTIuNGMtMS44LDAtMy0wLjMtMy44LTFjLTAuOC0wLjYtMS41LTItMi4xLTMuOUw5Ni43LDEwLjJjLTAuNi0yLTAuOS0zLjMtMC45LTQgICBjMC0xLjYsMC44LTIuNSwyLjQtMi41aDkuOGMxLjksMCwzLjIsMC4zLDMuOSwxYzAuOCwwLjYsMS40LDIsMiwzLjlsMTYuOCw2Ni4ybDE1LjYtNjYuMmMwLjUtMiwxLjEtMy4zLDEuOS0zLjljMC44LTAuNiwyLjItMSw0LTEgICBoOGMxLjksMCwzLjIsMC4zLDQsMWMwLjgsMC42LDEuNSwyLDEuOSwzLjlsMTUuOCw2N2wxNy4zLTY3YzAuNi0yLDEuMy0zLjMsMi0zLjljMC44LTAuNiwyLjEtMSwzLjktMWg5LjNjMS42LDAsMi41LDAuOCwyLjUsMi41ICAgYzAsMC41LTAuMSwxLTAuMiwxLjZjLTAuMSwwLjYtMC4zLDEuNC0wLjcsMi41bC0yNC4xLDc3LjNjLTAuNiwyLTEuMywzLjMtMi4xLDMuOWMtMC44LDAuNi0yLjEsMS0zLjgsMWgtOC42Yy0xLjksMC0zLjItMC4zLTQtMSAgIGMtMC44LTAuNy0xLjUtMi0xLjktNEwxNTYsMjNsLTE1LjQsNjQuNGMtMC41LDItMS4xLDMuMy0xLjksNGMtMC44LDAuNy0yLjIsMS00LDFIMTI2LjF6IE0yNTQuNiw5NS4xYy01LjIsMC0xMC40LTAuNi0xNS40LTEuOCAgIGMtNS0xLjItOC45LTIuNS0xMS41LTRjLTEuNi0wLjktMi43LTEuOS0zLjEtMi44Yy0wLjQtMC45LTAuNi0xLjktMC42LTIuOHYtNS4xYzAtMi4xLDAuOC0zLjEsMi4zLTMuMWMwLjYsMCwxLjIsMC4xLDEuOCwwLjMgICBjMC42LDAuMiwxLjUsMC42LDIuNSwxYzMuNCwxLjUsNy4xLDIuNywxMSwzLjVjNCwwLjgsNy45LDEuMiwxMS45LDEuMmM2LjMsMCwxMS4yLTEuMSwxNC42LTMuM2MzLjQtMi4yLDUuMi01LjQsNS4yLTkuNSAgIGMwLTIuOC0wLjktNS4xLTIuNy03Yy0xLjgtMS45LTUuMi0zLjYtMTAuMS01LjJMMjQ2LDUyYy03LjMtMi4zLTEyLjctNS43LTE2LTEwLjJjLTMuMy00LjQtNS05LjMtNS0xNC41YzAtNC4yLDAuOS03LjksMi43LTExLjEgICBjMS44LTMuMiw0LjItNiw3LjItOC4yYzMtMi4zLDYuNC00LDEwLjQtNS4yYzQtMS4yLDguMi0xLjcsMTIuNi0xLjdjMi4yLDAsNC41LDAuMSw2LjcsMC40YzIuMywwLjMsNC40LDAuNyw2LjUsMS4xICAgYzIsMC41LDMuOSwxLDUuNywxLjZjMS44LDAuNiwzLjIsMS4yLDQuMiwxLjhjMS40LDAuOCwyLjQsMS42LDMsMi41YzAuNiwwLjgsMC45LDEuOSwwLjksMy4zdjQuN2MwLDIuMS0wLjgsMy4yLTIuMywzLjIgICBjLTAuOCwwLTIuMS0wLjQtMy44LTEuMmMtNS43LTIuNi0xMi4xLTMuOS0xOS4yLTMuOWMtNS43LDAtMTAuMiwwLjktMTMuMywyLjhjLTMuMSwxLjktNC43LDQuOC00LjcsOC45YzAsMi44LDEsNS4yLDMsNy4xICAgYzIsMS45LDUuNywzLjgsMTEsNS41bDE0LjIsNC41YzcuMiwyLjMsMTIuNCw1LjUsMTUuNSw5LjZjMy4xLDQuMSw0LjYsOC44LDQuNiwxNGMwLDQuMy0wLjksOC4yLTIuNiwxMS42ICAgYy0xLjgsMy40LTQuMiw2LjQtNy4zLDguOGMtMy4xLDIuNS02LjgsNC4zLTExLjEsNS42QzI2NC40LDk0LjQsMjU5LjcsOTUuMSwyNTQuNiw5NS4xeiIvPgoJPGc+CgkJPHBhdGggY2xhc3M9InN0MSIgZD0iTTI3My41LDE0My43Yy0zMi45LDI0LjMtODAuNywzNy4yLTEyMS44LDM3LjJjLTU3LjYsMC0xMDkuNS0yMS4zLTE0OC43LTU2LjdjLTMuMS0yLjgtMC4zLTYuNiwzLjQtNC40ICAgIGM0Mi40LDI0LjYsOTQuNywzOS41LDE0OC44LDM5LjVjMzYuNSwwLDc2LjYtNy42LDExMy41LTIzLjJDMjc0LjIsMTMzLjYsMjc4LjksMTM5LjcsMjczLjUsMTQzLjd6Ii8+CgkJPHBhdGggY2xhc3M9InN0MSIgZD0iTTI4Ny4yLDEyOC4xYy00LjItNS40LTI3LjgtMi42LTM4LjUtMS4zYy0zLjIsMC40LTMuNy0yLjQtMC44LTQuNWMxOC44LTEzLjIsNDkuNy05LjQsNTMuMy01ICAgIGMzLjYsNC41LTEsMzUuNC0xOC42LDUwLjJjLTIuNywyLjMtNS4zLDEuMS00LjEtMS45QzI4Mi41LDE1NS43LDI5MS40LDEzMy40LDI4Ny4yLDEyOC4xeiIvPgoJPC9nPgo8L2c+Cjwvc3ZnPg==",
					MediaType: mediatypeSVG,
				},
			},
			Provider: opsv1alpha1.AppLink{
				Name: "Amazon, Inc.",
				URL:  "https://aws.amazon.com",
			},

			InstallModes: []opsv1alpha1.InstallMode{
				{
					Type:      opsv1alpha1.InstallModeTypeAllNamespaces,
					Supported: true,
				},
			},
		},
	}
}
