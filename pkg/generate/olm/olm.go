package olm

import (
	"fmt"
	"strings"
	ttpl "text/template"
	"time"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
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
	commonMeta ackmodel.OLMCommonMetadata,
	serviceMeta ackmodel.OLMMetadata,
	vers string,
	templateBasePath string,
) (*templateset.TemplateSet, error) {

	createdAt := time.Now().Format("2006-01-02 15:04:05")

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
		createdAt,
		g.MetaVars(),
		commonMeta,
		serviceMeta,
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
	Common ackmodel.OLMCommonMetadata
	ackmodel.OLMMetadata
	CRDs []*ackmodel.CRD
}
