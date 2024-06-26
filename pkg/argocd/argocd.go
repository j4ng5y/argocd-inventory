package argocd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type (
	ArgoClient struct {
		outfile string
		logger  *zerolog.Logger
		url     *url.URL
		token   string
	}

	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token,omitempty"`
	}

	LoginResponse struct {
		Token   string `json:"token"`
		Code    int64  `json:"code"`
		Details []struct {
			TypeURL string `json:"type_url"`
			Value   string `json:"value"`
		} `json:"details"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	GetApplicationsResponse struct {
		Items []struct {
			Metadata struct {
				Annotations                map[string]string `json:"annotations"`
				CreationTimestamp          time.Time         `json:"creationTimestamp"`
				DeletionGracePeriodSeconds int64             `json:"deletionGracePeriodSeconds"`
				DeletionTimestamp          time.Time         `json:"deletionTimestamp"`
				Finalizers                 []string          `json:"finalizers"`
				GenerateName               string            `json:"generateName"`
				Generation                 int64             `json:"generation"`
				Labels                     map[string]string `json:"labels"`
				ManagedFields              []map[string]any  `json:"managedFields"`
				Name                       string            `json:"name"`
				Namespace                  string            `json:"namespace"`
				OwnerReference             []map[string]any  `json:"ownerReferences"`
				ResourceVersion            string            `json:"resourceVersion"`
				SelfLink                   string            `json:"selfLink"`
				UID                        string            `json:"uid"`
			} `json:"metadata"`
			Operation map[string]any `json:"operation"`
			Spec      map[string]any `json:"spec"`
			Status    struct {
				Conditions          []map[string]any `json:"conditions"`
				ControllerNamespace string           `json:"controllerNamespace"`
				Health              struct {
					Message string `json:"message"`
					Status  string `json:"status"`
				} `json:"health"`
				History              []map[string]any `json:"history"`
				ObservedAt           time.Time        `json:"observedAt"`
				OperationState       map[string]any   `json:"operationState"`
				ReconciledAt         time.Time        `json:"reconciledAt"`
				ResourceHealthSource string           `json:"resourceHealthSource"`
				Resources            []map[string]any `json:"resources"`
				SourceType           string           `json:"sourceType"`
				SourceTypes          []string         `json:"sourceTypes"`
				Summary              map[string]any   `json:"summary"`
				Sync                 struct {
					ComparedTo map[string]any `json:"comparedTo"`
					Revision   string
					Revisions  []string `json:"revisions"`
					Status     string   `json:"status"`
				} `json:"sync"`
			} `json:"status"`
		} `json:"items"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	GetResouceTreeResponse struct {
		Hosts []struct {
			Name          string           `json:"name"`
			ResourcesInfo []map[string]any `json:"resourcesInfo"`
			SystemInfo    map[string]any   `json:"systemInfo"`
		} `json:"hosts"`
		Nodes []struct {
			CreatedAt       time.Time        `json:"createdAt"`
			Health          map[string]any   `json:"health"`
			Images          []string         `json:"images"`
			Info            []map[string]any `json:"info"`
			NetworkingInfo  map[string]any   `json:"networkingInfo"`
			ParentRefs      []map[string]any `json:"parentRefs"`
			ResourceVersion string           `json:"resourceVersion"`
			Group           string           `json:"group"`
			Kind            string           `json:"kind"`
			Name            string           `json:"name"`
			Namespace       string           `json:"namespace"`
			UID             string           `json:"uid"`
			Version         string           `json:"version"`
		} `json:"nodes"`
		OrphanedNodes []struct {
			CreatedAt       time.Time        `json:"createdAt"`
			Health          map[string]any   `json:"health"`
			Images          []string         `json:"images"`
			Info            []map[string]any `json:"info"`
			NetworkingInfo  map[string]any   `json:"networkingInfo"`
			ParentRefs      []map[string]any `json:"parentRefs"`
			ResourceVersion string           `json:"resourceVersion"`
			Group           string           `json:"group"`
			Kind            string           `json:"kind"`
			Name            string           `json:"name"`
			Namespace       string           `json:"namespace"`
			UID             string           `json:"uid"`
			Version         string           `json:"version"`
		} `json:"orphanedNodes"`
	}
)

func NewArgoClient(url *url.URL, username, password, out string, logger *zerolog.Logger) (*ArgoClient, error) {
	client := &ArgoClient{
		outfile: out,
		logger:  logger,
		url:     url,
	}
	client.logger.Debug().Any("url", client.url).Msg("creating new argocd client")

	if err := client.login(username, password); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ArgoClient) FetchApplications() error {
	out, err := os.OpenFile(c.outfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	if err := writer.Write([]string{"Application", "Group", "Version", "Kind", "Namespace/Name"}); err != nil {
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
		if app.Status.Health.Status == "Healthy" || app.Status.Health.Status == "Progressing" {
			c.logger.Debug().Str("app", app.Metadata.Name).Msg("matched healthy/progressing app: fetching resource tree")
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
				if err := writer.Write([]string{app.Metadata.Name, node.Group, node.Version, node.Kind, node.Namespace + "/" + node.Name}); err != nil {
					return err
				}
			}

			c.logger.Debug().Str("app", app.Metadata.Name).Msg("fetched resource tree")
		} else {
			c.logger.Warn().Str("app", app.Metadata.Name).Str("status", app.Status.Health.Status).Msg("failed to match health status of healthy/progressing")
		}
	}

	return nil
}

func (c *ArgoClient) login(username, password string) error {
	c.logger.Debug().Any("url", c.url).Msg("logging into argocd")

	c.logger.Debug().Msg("marshalling login request")
	b, err := json.Marshal(&LoginRequest{
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
	r := new(LoginResponse)

	if err := json.Unmarshal(body, r); err != nil {
		return err
	}

	c.logger.Debug().Msg("setting token")
	c.token = r.Token

	c.logger.Debug().Msg("logged into argocd")

	return nil
}
