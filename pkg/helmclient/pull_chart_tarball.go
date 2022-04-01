package helmclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
)

const (
	OCIScheme = "oci"
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

	if u.Scheme == OCIScheme {
		chartTarballPath, err = c.doFileOCI(ctx)
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

func (c *Client) doFileOCI(ctx context.Context) (string, error) {
	// TODO(kuba): Download tarball from OCI registry with retries, then return
	// its path.
	// Utilize ORAS project for the request, reimplementing
	// https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/pkg/registry/client.go#L256-L414.
	// I think we can ignore provenance layer for now, adding it to the implementation only if it becomes necessary.
	return "", nil
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
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				return microerror.Mask(err)
			}

			// Github Pages 404 produces full HTML page which obscures the logs.
			if resp.StatusCode == http.StatusNotFound {
				return backoff.Permanent(microerror.Maskf(pullChartNotFoundError, fmt.Sprintf("got StatusCode %d for url %#q", resp.StatusCode, req.URL.String())))
			}

			// Github Pages 503 produces full HTML page which obscures the logs.
			if resp.StatusCode == http.StatusServiceUnavailable {
				return backoff.Permanent(microerror.Maskf(pullChartFailedError, fmt.Sprintf("got StatusCode %d for url %#q", resp.StatusCode, req.URL.String())))
			}

			return microerror.Maskf(executionFailedError, fmt.Sprintf("got StatusCode %d for url %#q with body %s", resp.StatusCode, req.URL.String(), buf.String()))
		}

		tmpfile, err := afero.TempFile(c.fs, "", "chart-tarball")
		if err != nil {
			return microerror.Mask(err)
		}
		defer tmpfile.Close()

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
