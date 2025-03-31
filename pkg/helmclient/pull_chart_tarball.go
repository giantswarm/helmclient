package helmclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
	helmregistry "helm.sh/helm/v3/pkg/registry"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"oras.land/oras-go/pkg/registry"
)

// PullChartTarball downloads a tarball from the provided tarball URL,
// returning the file path.
func (c *Client) PullChartTarball(ctx context.Context, tarballURL string) (string, error) {
	eventName := "pull_chart_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	chartTarballPath, err := c.pullChartTarball(ctx, tarballURL)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return "", microerror.Mask(err)
	}

	return chartTarballPath, nil
}

func (c *Client) pullChartTarball(ctx context.Context, tarballURL string) (string, error) {
	u, err := url.Parse(tarballURL)
	if err != nil {
		return "", microerror.Mask(err)
	}

	var chartTarballPath string

	if u.Scheme == helmregistry.OCIScheme {
		chartTarballPath, err = c.doFileOCI(ctx, tarballURL)
		if err != nil {
			return "", microerror.Mask(err)
		}
	} else {
		req, err := c.newRequest("GET", tarballURL)
		if err != nil {
			return "", microerror.Mask(err)
		}

		// Set host header to prevent 404 responses from GitHub Pages.
		req.Host = u.Host

		chartTarballPath, err = c.doFileHTTP(ctx, req)
		if err != nil {
			return "", microerror.Mask(err)
		}
	}

	return chartTarballPath, nil
}

func (c *Client) doFileOCI(ctx context.Context, url string) (string, error) {
	var tmpFileName string

	o := func() error {
		// We utilize 'oci://' scheme to recognize OCI registries, but
		// registry.ParseReference has a strict regex. Let's get rid of
		// protocol prefix.
		url = strings.TrimPrefix(url, helmregistry.OCIScheme+"://")
		ref, err := registry.ParseReference(url)
		if err != nil {
			return microerror.Maskf(pullChartFailedError, "error parsing url: %s", err)
		}
		memoryStore := content.NewMemory()
		// We accept Config layer (required), Chart layer (needed), and
		// Provenance layer (ignored).
		allowedMediaTypes := []string{
			helmregistry.ConfigMediaType,
			helmregistry.ChartLayerMediaType,
			helmregistry.LegacyChartLayerMediaType,
			helmregistry.ProvLayerMediaType,
		}

		var registryStore content.Registry
		{
			// We make an assumption that every registry we pull from is
			// public. Configuration provided to Client can override that, but
			// is optional.
			resolver, err := content.NewRegistry(c.registryOptions)
			if err != nil {
				return microerror.Maskf(pullChartFailedError, "error creating registry resolver: %s", err)
			}
			registryStore = content.Registry{Resolver: resolver}
		}

		var descriptors, layers []ocispec.Descriptor
		manifest, err := oras.Copy(ctx, registryStore, ref.String(), memoryStore, "",
			oras.WithPullEmptyNameAllowed(),
			oras.WithAllowedMediaTypes(allowedMediaTypes),
			oras.WithLayerDescriptors(func(l []ocispec.Descriptor) {
				layers = l
			}))
		if err != nil {
			return microerror.Maskf(pullChartFailedError, "error copying manifests: %s", err)
		}

		descriptors = append(descriptors, manifest)
		descriptors = append(descriptors, layers...)
		// We expect at least a config layer and a chart layer. Count may be higher
		// if the provenance layer is present.
		if len(descriptors) < 2 {
			return microerror.Maskf(pullChartFailedError,
				"manifest does not contain minimum number of descriptors (2), descriptors found: %d",
				len(descriptors),
			)
		}
		var configDescriptor, chartDescriptor *ocispec.Descriptor
		for _, descriptor := range descriptors {
			d := descriptor
			switch d.MediaType {
			case helmregistry.ConfigMediaType:
				configDescriptor = &d
			case helmregistry.ChartLayerMediaType:
				chartDescriptor = &d
			case helmregistry.LegacyChartLayerMediaType:
				chartDescriptor = &d
			}
		}
		if configDescriptor == nil {
			// configDescriptor is required as proof of successful chart pull,
			// although not used in further code. It contains chart metadata,
			// which might prove useful some day.
			return microerror.Maskf(pullChartFailedError, "could not load config with mediatype %s", helmregistry.ConfigMediaType)
		}
		if chartDescriptor == nil {
			return microerror.Maskf(pullChartFailedError, "manifest does not contain a layer with mediatype %s", helmregistry.ChartLayerMediaType)
		}

		_, chartData, ok := memoryStore.Get(*chartDescriptor)
		if !ok {
			return microerror.Maskf(pullChartFailedError, "unable to retrieve blob with digest %s", chartDescriptor.Digest)
		}

		tmpfile, err := afero.TempFile(c.fs, "", "chart-tarball")
		if err != nil {
			return microerror.Mask(err)
		}
		defer func() { _ = tmpfile.Close() }()

		buf := bytes.NewBuffer(chartData)
		_, err = io.Copy(tmpfile, buf)
		if err != nil {
			return microerror.Mask(err)
		}

		tmpFileName = tmpfile.Name()

		return nil
	}

	b := backoff.NewMaxRetries(3, 5*time.Second)
	n := backoff.NewNotifier(c.logger, ctx)

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return tmpFileName, nil
}

func (c *Client) doFileHTTP(ctx context.Context, req *http.Request) (string, error) {
	var tmpFileName string

	req = req.WithContext(ctx)

	o := func() error {
		resp, err := c.httpClient.Do(req)
		if isNoSuchHostError(err) {
			return backoff.Permanent(microerror.Maskf(pullChartFailedError, "no such host %#q", req.Host))
		} else if IsPullChartTimeout(err) {
			return backoff.Permanent(microerror.Maskf(pullChartTimeoutError, "%#q timeout for %#q", req.Method, req.URL.String()))
		} else if err != nil {
			return microerror.Mask(err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				return microerror.Mask(err)
			}

			// Github Pages 404 produces full HTML page which obscures the logs.
			if resp.StatusCode == http.StatusNotFound {
				return backoff.Permanent(microerror.Maskf(pullChartNotFoundError, "got StatusCode %d for url %#q", resp.StatusCode, req.URL.String()))
			}

			// Github Pages 503 produces full HTML page which obscures the logs.
			if resp.StatusCode == http.StatusServiceUnavailable {
				return backoff.Permanent(microerror.Maskf(pullChartFailedError, "got StatusCode %d for url %#q", resp.StatusCode, req.URL.String()))
			}

			return microerror.Maskf(executionFailedError, "got StatusCode %d for url %#q with body %s", resp.StatusCode, req.URL.String(), buf.String())
		}

		tmpfile, err := afero.TempFile(c.fs, "", "chart-tarball")
		if err != nil {
			return microerror.Mask(err)
		}
		defer func() { _ = tmpfile.Close() }()

		_, err = io.Copy(tmpfile, resp.Body)
		if err != nil {
			return microerror.Mask(err)
		}

		tmpFileName = tmpfile.Name()

		return nil
	}

	b := backoff.NewMaxRetries(3, 5*time.Second)
	n := backoff.NewNotifier(c.logger, ctx)

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return tmpFileName, nil
}

func (c *Client) newRequest(method, url string) (*http.Request, error) {
	var buf io.Reader

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	return req, nil
}
