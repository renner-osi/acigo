package aci

import (
	"bytes"
	"fmt"
)

func rnOut(out string) string {
	return "out-" + out
}

func dnL3ExtOut(tenant, out string) string {
	return rnTenant(tenant) + "/" + rnOut(out)
}

// L3ExtOutAdd creates a new external routed network in a tenant.
func (c *Client) L3ExtOutAdd(tenant, out, descr string) error {

	me := "L3ExtOutAdd"

	rn := rnOut(out)

	dn := dnL3ExtOut(tenant, out)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"l3extOut":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s","rn":"%s","status":"created"}}}`,
		dn, out, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// L3ExtOutDel deletes an external routed network from a tenant.
func (c *Client) L3ExtOutDel(tenant, out string) error {

	me := "L3ExtOutDel"

	rnT := rnTenant(tenant)

	dn := dnL3ExtOut(tenant, out)

	api := "/api/node/mo/uni/" + rnT + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvTenant":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"l3extOut":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rnT, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// L3ExtOutList retrieves the list of external routed networks from a tenant.
func (c *Client) L3ExtOutList(tenant string) ([]map[string]interface{}, error) {

	me := "L3ExtOutList"

	key := "l3extOut"

	t := rnTenant(tenant)

	api := "/api/node/mo/uni/" + t + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// L3ExtOutVrfSet defines the VRF for an external routed network.
func (c *Client) L3ExtOutVrfSet(tenant, out, vrf string) error {

	me := "L3ExtOutVrfSet"

	dn := dnL3ExtOut(tenant, out)

	api := "/api/node/mo/uni/" + dn + "/rsectx.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"l3extRsEctx":{"attributes":{"tnFvCtxName":"%s"}}}`,
		vrf)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// L3ExtOutVrfGet retrieves the VRF for an external routed network.
func (c *Client) L3ExtOutVrfGet(tenant, out string) (string, error) {

	me := "L3ExtOutVrfGet"

	key := "l3extRsEctx"

	dn := dnL3ExtOut(tenant, out)

	api := "/api/node/mo/uni/" + dn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, me)
	if errAttr != nil {
		return "", fmt.Errorf("%s: %v", me, errAttr)
	}

	if len(attrs) < 1 {
		return "", fmt.Errorf("%s: empty list of VRFs", me)
	}

	attr := attrs[0]

	v, found := attr["tnFvCtxName"]
	if !found {
		return "", fmt.Errorf("%s: VRF not found", me)
	}

	vrf, isStr := v.(string)
	if !isStr {
		return "", fmt.Errorf("%s: VRF is not a string", me)
	}

	if vrf == "" {
		return "", fmt.Errorf("%s: empty VRF name", me)
	}

	return vrf, nil
}
