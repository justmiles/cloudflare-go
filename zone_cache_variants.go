package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ZoneCacheVariants struct {
	ID         string              `json:"id"`
	ModifiedOn string              `json:"modified_on,omitempty"`
	Value      map[string][]string `json:"value"`
}

type ZoneCacheVariantsSingleResponse struct {
	Response
	Result ZoneCacheVariants `json:"result"`
}

// ZoneCacheVariants returns information about the current cache variants
//
// API reference: https://api.cloudflare.com/#zone-cache-settings-get-variants-setting
func (api *API) ZoneCacheVariants(ctx context.Context, zoneID string) (ZoneCacheVariants, error) {
	uri := fmt.Sprintf("/zones/%s/cache/variants", zoneID)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return ZoneCacheVariants{}, err
	}
	var r ZoneCacheVariantsSingleResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return ZoneCacheVariants{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// ZoneCacheVariants updates the cache variants for a given zone.
//
// API reference: https://api.cloudflare.com/#zone-cache-settings-change-variants-setting
func (api *API) UpdateZoneCacheVariants(ctx context.Context, zoneID string, setting ZoneCacheVariants) (ZoneCacheVariants, error) {
	uri := fmt.Sprintf("/zones/%s/cache/variants", zoneID)
	res, err := api.makeRequestContext(ctx, http.MethodPatch, uri, setting)
	if err != nil {
		return ZoneCacheVariants{}, err
	}

	response := &ZoneCacheVariantsSingleResponse{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		return ZoneCacheVariants{}, errors.Wrap(err, errUnmarshalError)
	}

	return response.Result, nil
}

// DeleteZoneCacheVariants deletes cache variants for a given zone.
//
// API reference: https://api.cloudflare.com/#zone-cache-settings-delete-variants-setting
func (api *API) DeleteZoneCacheVariants(ctx context.Context, zoneID string) error {
	uri := fmt.Sprintf("/zones/%s/cache/variants", zoneID)
	_, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	return nil
}