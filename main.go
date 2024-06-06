package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	outfile, logLevel, argoUrl, argoUsername, argoPassword string
)

type (
	ArgoClient struct {
		outfile string
		logger *zerolog.Logger
		url    *url.URL
		token  string
	}
)

func init() {
	flag.StringVar(&argoUsername, "argo-username", "", "The ArgoCD username to use.")
	flag.StringVar(&argoPassword, "argo-password", "", "The ArgoCD password to use.")
	flag.StringVar(&argoUrl, "argo-url", "", "The ArgoCD URL to use.")
	flag.StringVar(&logLevel, "log-level", "info", "The log level to use.")
	flag.StringVar(&outfile, "out", "report.csv", "The output file to write the report to.")
	flag.Parse()
}

func main() {
	var (
		argourl *url.URL
		err     error
	)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	loglvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		logger.Warn().Err(err).Msg("log level failed to parse, defaulting to info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(loglvl)
	}

	if argoUrl == "" {
		logger.Fatal().Msg("argo-url flag was not set")
	} else {
		argourl, err = url.Parse(argoUrl)
		if err != nil {
			logger.Fatal().Err(err).Send()
		} else {
			logger.Debug().Msgf("argo-url set to '%s'", argourl.String())
		}
	}
	if argoUsername == "" {
		logger.Fatal().Msg("argo-username flag was not set")
	} else {
		logger.Debug().Msgf("argo-username set to '%s'", argoUsername)
	}
	if argoPassword == "" {
		logger.Fatal().Msg("argo-password flag was not set")
	} else {
		logger.Debug().Msgf("argo-password is set: %s", strings.Repeat("*", len(argoPassword)))
	}

	if argourl == nil {
		logger.Fatal().Msg("argourl parsed to nil")
	}

	client, err := NewArgoClient(argourl, argoUsername, argoPassword, outfile, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	if err := client.FetchApplications(); err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("done")
}

func NewArgoClient(url *url.URL, username, password, out string, logger *zerolog.Logger) (*ArgoClient, error) {
	client := &ArgoClient{
		outfile: out,
		logger: logger,
		url:    url,
	}
	client.logger.Debug().Any("url", client.url).Msg("creating new argocd client")

	if err := client.login(username, password); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ArgoClient) FetchApplications() error {
	type (
		GetApplicationsResponse struct {
			Items []struct {
				Metadata struct {
					Annotations                map[string]string `json:"annotations"`
					CreationTimestamp          time.Time         `json:"creationTimestamp"`
					DeletionGracePeriodSeconds int64 					 `json:"deletionGracePeriodSeconds"`
					DeletionTimestamp          time.Time				 `json:"deletionTimestamp"`
					Finalizers                 []string					 `json:"finalizers"`
					GenerateName               string						 `json:"generateName"`
					Generation                 int64							 `json:"generation"`
					Labels                     map[string]string `json:"labels"`
					ManagedFields              []struct {
						APIVersion  string `json:"apiVersion"`
						FieldsType  string `json:"fieldsType"`
						FieldsV1    struct{
							Raw string `json:"Raw"`
						} `json:"fieldsV1"`
						Manager     string `json:"manager"`
						Operation   string `json:"operation"`
						Subresource string `json:"subresource"`
						Time        time.Time `json:"time"`
					} `json:"managedFields"`
					Name           string `json:"name"`
					Namespace      string `json:"namespace"`
					OwnerReference []struct {
						APIVersion         string `json:"apiVersion"`
						BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
						Controller         bool   `json:"controller"`
						Kind               string `json:"kind"`
						Name               string `json:"name"`
						UID                string `json:"uid"`
					} `json:"ownerReferences"`
					ResourceVersion string `json:"resourceVersion"`
					SelfLink        string `json:"selfLink"`
					UID             string `json:"uid"`
				} `json:"metadata"`
				Operation map[string]interface{} `json:"operation"`
				Spec      map[string]interface{} `json:"spec"`
				Status    map[string]interface{} `json:"status"`
			} `json:"items"`
			Metadata struct {
				Continue           string `json:"continue"`
				RemainingItemCount int64  `json:"remainingItemCount"`
				ResourceVersios    string `json:"resourceVersion"`
				SelfLink           string `json:"selfLink"`
			} `json:"metadata"`
		}
		GetResouceTreeResponse struct {
			Hosts []struct {
				Name          string                   `json:"name"`
				ResourcesInfo []map[string]interface{} `json:"resourcesInfo"`
				SystemInfo    map[string]interface{}   `json:"systemInfo"`
			} `json:"hosts"`
			Nodes []struct {
				CreatedAt       time.Time                `json:"createdAt"`
				Health          map[string]interface{}   `json:"health"`
				Images          []string                 `json:"images"`
				Info            []map[string]interface{} `json:"info"`
				NetworkingInfo  map[string]interface{}   `json:"networkingInfo"`
				ParentRefs      []map[string]interface{} `json:"parentRefs"`
				ResourceVersion string                   `json:"resourceVersion"`
				Group           string                   `json:"group"`
				Kind            string                   `json:"kind"`
				Name            string                   `json:"name"`
				Namespace       string                   `json:"namespace"`
				UID             string                   `json:"uid"`
				Version         string                   `json:"version"`
			} `json:"nodes"`
			OrphanedNodes []struct {
				CreatedAt       time.Time                `json:"createdAt"`
				Health          map[string]interface{}   `json:"health"`
				Images          []string                 `json:"images"`
				Info            []map[string]interface{} `json:"info"`
				NetworkingInfo  map[string]interface{}   `json:"networkingInfo"`
				ParentRefs      []map[string]interface{} `json:"parentRefs"`
				ResourceVersion string                   `json:"resourceVersion"`
				Group           string                   `json:"group"`
				Kind            string                   `json:"kind"`
				Name            string                   `json:"name"`
				Namespace       string                   `json:"namespace"`
				UID             string                   `json:"uid"`
				Version         string                   `json:"version"`
			} `json:"orphanedNodes"`
		}
	)

	out, err := os.OpenFile(c.outfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	if err := writer.Write([]string{"Application", "Group", "Version", "Kind", "Name"}); err != nil {
		return err
	}

	getAppsReq, err := http.NewRequest(http.MethodGet, c.url.JoinPath("/api/v1/applications").String(), nil)
	if err != nil {
		return err
	}

	getAppsReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	getAppsResp, err := http.DefaultClient.Do(getAppsReq)
	if err != nil {
		return err
	}

	if getAppsResp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code recieved while fetching applications: %d", getAppsResp.StatusCode)
	}

	body, err := io.ReadAll(getAppsResp.Body)
	if err != nil {
		return err
	}

	r := new(GetApplicationsResponse)
	if err := json.Unmarshal(body, r); err != nil {
		return err
	}

	for _, app := range r.Items {
		c.logger.Debug().Str("app", app.Metadata.Name).Msg("fetching resource tree")
		getResourceTreeReq, err := http.NewRequest(http.MethodGet, c.url.JoinPath("/api/v1/applications", app.Metadata.Name, "resource-tree").String(), nil)
		if err != nil {
			return err
		}

		getResourceTreeReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

		getResourceTreeResp, err := http.DefaultClient.Do(getResourceTreeReq)
		if err != nil {
			return err
		}

		if getResourceTreeResp.StatusCode != http.StatusOK {
			return fmt.Errorf("invalid status code recieved while fetching resource tree for application %s: %d", app.Metadata.Name, getResourceTreeResp.StatusCode)
		}

		body, err := io.ReadAll(getResourceTreeResp.Body)
		if err != nil {
			return err
		}

		r := new(GetResouceTreeResponse)
		if err := json.Unmarshal(body, r); err != nil {
			return err
		}
     
		for _, node := range r.Nodes {
			if err := writer.Write([]string{app.Metadata.Name, node.Group, node.Version, node.Kind, node.Name}); err != nil {
				return err
			}
		}

		c.logger.Debug().Str("app", app.Metadata.Name).Msg("fetched resource tree")
	}

	return nil
}

func (c *ArgoClient) login(username, password string) error {
	c.logger.Debug().Any("url", c.url).Msg("logging into argocd")
	type (
		Request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		Response struct {
			Token   string `json:"token"`
			Code    int64  `json:"code"`
			Details []struct {
				TypeURL string `json:"type_url"`
				Value   string `json:"value"`
			} `json:"details"`
			Error   string `json:"error"`
			Message string `json:"message"`
		}
	)

	c.logger.Debug().Msg("marshalling login request")
	b, err := json.Marshal(&Request{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	c.logger.Debug().Str("url", c.url.JoinPath("api/v1/session").String()).Msg("sending login request")
	req, err := http.NewRequest(http.MethodPost, c.url.JoinPath("api/v1/session").String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.logger.Debug().Msg("recieved response from argocd")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code recieved while logging into argocd: %d", resp.StatusCode)
	}

	c.logger.Debug().Msg("reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.logger.Debug().Msg("unmarshalling response body")
	r := new(Response)

	if err := json.Unmarshal(body, r); err != nil {
		return err
	}

	c.logger.Debug().Msg("setting token")
	c.token = r.Token

	c.logger.Debug().Msg("logged into argocd")

	return nil
}
