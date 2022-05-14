package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

// Tunnel is the struct definition of a tunnel.
type Tunnel struct {
	ID             string             `json:"id,omitempty"`
	Name           string             `json:"name,omitempty"`
	Secret         string             `json:"tunnel_secret,omitempty"`
	CreatedAt      *time.Time         `json:"created_at,omitempty"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty"`
	Connections    []TunnelConnection `json:"connections,omitempty"`
	ConnsActiveAt  *time.Time         `json:"conns_active_at,omitempty"`
	ConnInactiveAt *time.Time         `json:"conns_inactive_at,omitempty"`
}

// TunnelConnection represents the connections associated with a tunnel.
type TunnelConnection struct {
	ColoName           string `json:"colo_name"`
	ID                 string `json:"id"`
	IsPendingReconnect bool   `json:"is_pending_reconnect"`
	ClientID           string `json:"client_id"`
	ClientVersion      string `json:"client_version"`
	OpenedAt           string `json:"opened_at"`
	OriginIP           string `json:"origin_ip"`
}
type TunnelConfiguration struct {
	Ingress     []TunnelConfigurationIngress `json:"ingress"`
	WarpRouting struct {
		Enabled bool `json:"enabled"`
	} `json:"warp-routing"`
}

type TunnelConfigurationIngress struct {
	Service       string                                  `json:"service"`
	Hostname      string                                  `json:"hostname,omitempty"`
	OriginRequest TunnelConfigurationIngressOriginRequest `json:"originRequest,omitempty"`
}

type TunnelConfigurationIngressOriginRequest struct {
	HTTPHostHeader string `json:"httpHostHeader"`
}

type TunnelConfigurationDetail struct {
	Config    TunnelConfiguration `json:"config"`
	TunnelID  string              `json:"tunnel_id"`
	Version   int                 `json:"version"`
	CreatedAt time.Time           `json:"created_at"`
}

// TunnelsDetailResponse is used for representing the API response payload for
// multiple tunnels.
type TunnelsDetailResponse struct {
	Result []Tunnel `json:"result"`
	Response
}

// TunnelDetailResponse is used for representing the API response payload for
// a single tunnel.
type TunnelDetailResponse struct {
	Result Tunnel `json:"result"`
	Response
}

// TunnelTokenResponse is  the API response for a tunnel token.
type TunnelTokenResponse struct {
	Result string `json:"result"`
	Response
}

// ConfiguratonDetailResponse ist the api response for multiple configuration
// detail requests
type ConfiguratonDetailResponse struct {
	Result TunnelConfigurationDetail `json:"result"`
	Response
}

type TunnelParams struct {
	AccountID string
	ID        string
}

type TunnelCreateParams struct {
	AccountID string `json:"-"`
	Name      string `json:"name,omitempty"`
	Secret    string `json:"tunnel_secret,omitempty"`
}

type TunnelUpdateParams struct {
	AccountID string `json:"-"`
	Name      string `json:"name,omitempty"`
	Secret    string `json:"tunnel_secret,omitempty"`
}

type TunnelDeleteParams struct {
	AccountID string
	ID        string
}

type TunnelCleanupParams struct {
	AccountID string
	ID        string
}

type TunnelTokenParams struct {
	AccountID string
	ID        string
}

type TunnelListParams struct {
	AccountID string
	Name      string     `url:"name,omitempty"`
	UUID      string     `url:"uuid,omitempty"` // the tunnel ID
	IsDeleted bool       `url:"is_deleted,omitempty"`
	ExistedAt *time.Time `url:"existed_at,omitempty"`
}

type TunnelConfigurationUpdateParams struct {
	AccountID string              `json:"-"`
	TunnelID  string              `json:"tunnel_id"`
	Config    TunnelConfiguration `json:"config"`
}

type TunnelConfigurationParams struct {
	AccountID string `json:"-"`
	TunnelID  string `json:"tunnel_id"`
}

// Tunnels lists all tunnels.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-list-cloudflare-tunnels
func (api *API) Tunnels(ctx context.Context, params TunnelListParams) ([]Tunnel, error) {
	if params.AccountID == "" {
		return []Tunnel{}, ErrMissingAccountID
	}

	v, _ := query.Values(params)
	queryParams := v.Encode()
	if queryParams != "" {
		queryParams = "?" + queryParams
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel", params.AccountID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodGet, uri+queryParams, nil, nil)
	if err != nil {
		return []Tunnel{}, err
	}

	var argoDetailsResponse TunnelsDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return []Tunnel{}, errors.Wrap(err, errUnmarshalError)
	}
	return argoDetailsResponse.Result, nil
}

// Tunnel returns a single Argo tunnel.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-get-cloudflare-tunnel
func (api *API) Tunnel(ctx context.Context, params TunnelParams) (Tunnel, error) {
	if params.AccountID == "" {
		return Tunnel{}, ErrMissingAccountID
	}

	if params.ID == "" {
		return Tunnel{}, errors.New("missing tunnel ID")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s", params.AccountID, params.ID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return Tunnel{}, err
	}

	var argoDetailsResponse TunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, errUnmarshalError)
	}
	return argoDetailsResponse.Result, nil
}

// CreateTunnel creates a new tunnel for the account.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-create-cloudflare-tunnel
func (api *API) CreateTunnel(ctx context.Context, params TunnelCreateParams) (Tunnel, error) {
	if params.AccountID == "" {
		return Tunnel{}, ErrMissingAccountID
	}

	if params.Name == "" {
		return Tunnel{}, errors.New("missing tunnel name")
	}

	if params.Secret == "" {
		return Tunnel{}, errors.New("missing tunnel secret")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel", params.AccountID)

	tunnel := Tunnel{Name: params.Name, Secret: params.Secret}

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodPost, uri, tunnel, nil)
	if err != nil {
		return Tunnel{}, err
	}

	var argoDetailsResponse TunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, errUnmarshalError)
	}

	return argoDetailsResponse.Result, nil
}

// UpdateTunnel updates an existing tunnel for the account.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-update-cloudflare-tunnel
func (api *API) UpdateTunnel(ctx context.Context, params TunnelUpdateParams) (Tunnel, error) {
	if params.AccountID == "" {
		return Tunnel{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel", params.AccountID)

	var tunnel Tunnel

	if params.Name != "" {
		tunnel.Name = params.Name
	}

	if params.Secret != "" {
		tunnel.Secret = params.Secret
	}

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodPatch, uri, tunnel, nil)
	if err != nil {
		return Tunnel{}, err
	}

	var argoDetailsResponse TunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, errUnmarshalError)
	}

	return argoDetailsResponse.Result, nil
}

// DeleteTunnel removes a single Argo tunnel.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-delete-cloudflare-tunnel
func (api *API) DeleteTunnel(ctx context.Context, params TunnelDeleteParams) error {
	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s", params.AccountID, params.ID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodDelete, uri, nil, nil)
	if err != nil {
		return err
	}

	var argoDetailsResponse TunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}

	return nil
}

// CleanupTunnelConnections deletes any inactive connections on a tunnel.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-clean-up-cloudflare-tunnel-connections
func (api *API) CleanupTunnelConnections(ctx context.Context, params TunnelCleanupParams) error {
	if params.AccountID == "" {
		return ErrMissingAccountID
	}

	if params.ID == "" {
		return errors.New("missing tunnel ID")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/connections", params.AccountID, params.ID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodDelete, uri, nil, nil)
	if err != nil {
		return err
	}

	var argoDetailsResponse TunnelDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}

	return nil
}

// TunnelToken that allows to run a tunnel.
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-get-cloudflare-tunnel-token
func (api *API) TunnelToken(ctx context.Context, params TunnelTokenParams) (string, error) {
	if params.AccountID == "" {
		return "", ErrMissingAccountID
	}

	if params.ID == "" {
		return "", errors.New("missing tunnel ID")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/token", params.AccountID, params.ID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return "", err
	}

	var tunnelTokenResponse TunnelTokenResponse
	err = json.Unmarshal(res, &tunnelTokenResponse)
	if err != nil {
		return "", errors.Wrap(err, errUnmarshalError)
	}

	return tunnelTokenResponse.Result, nil
}

// PutTunnelConfiguration to add or update a configuration for a cloudflared tunnel
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-configuration-put-configuration
func (api *API) TunnelConfigurationUpdate(ctx context.Context, params TunnelConfigurationUpdateParams) (TunnelConfiguration, error) {
	if params.AccountID == "" {
		return TunnelConfiguration{}, ErrMissingAccountID
	}

	if params.TunnelID == "" {
		return TunnelConfiguration{}, errors.New("missing tunnel id")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/configurations", params.AccountID, params.TunnelID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodPut, uri, params, nil)
	if err != nil {
		return TunnelConfiguration{}, err
	}

	var argoDetailsResponse ConfiguratonDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		return TunnelConfiguration{}, errors.Wrap(err, errUnmarshalError)
	}

	return argoDetailsResponse.Result.Config, nil
}

// GetTunnelConfiguration to add or update a configuration for a cloudflared tunnel
//
// API reference: https://api.cloudflare.com/#cloudflare-tunnel-configuration-get-configuration
func (api *API) TunnelConfiguration(ctx context.Context, params TunnelConfigurationParams) (TunnelConfiguration, error) {
	if params.AccountID == "" {
		return TunnelConfiguration{}, ErrMissingAccountID
	}

	if params.TunnelID == "" {
		return TunnelConfiguration{}, errors.New("missing tunnel id")
	}

	uri := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/configurations", params.AccountID, params.TunnelID)

	res, err := api.makeRequestContextWithHeaders(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return TunnelConfiguration{}, err
	}

	var argoDetailsResponse ConfiguratonDetailResponse
	err = json.Unmarshal(res, &argoDetailsResponse)
	if err != nil {
		fmt.Println(string(res))
		return TunnelConfiguration{}, errors.Wrap(err, errUnmarshalError)
	}

	return argoDetailsResponse.Result.Config, nil
}
