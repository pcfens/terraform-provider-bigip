package bigip

import (
	"log"
	"strconv"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmPersistenceProfileSrcAddr() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileSrcAddrCreate,
		Read:   resourceBigipLtmPersistenceProfileSrcAddrRead,
		Update: resourceBigipLtmPersistenceProfileSrcAddrUpdate,
		Delete: resourceBigipLtmPersistenceProfileSrcAddrDelete,
		Exists: resourceBigipLtmPersistenceProfileSrcAddrExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipPersistenceProfileSrcAddrImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the persistence profile",
				ValidateFunc: validateF5Name,
			},

			"app_service": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},

			"defaults_from": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Inherit defaults from parent profile",
				ValidateFunc: validateF5Name,
			},

			"match_across_pools": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across pools with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"match_across_services": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across services with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"match_across_virtuals": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across services with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"mirror": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable",
				ValidateFunc: validateEnabledDisabled,
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for persistence of the session",
			},

			"override_conn_limit": {
				Type:         schema.TypeString,
				Default:      false,
				Optional:     true,
				Description:  "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
				ValidateFunc: validateEnabledDisabled,
			},

			// Specific to SourceAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Specify the hash algorithm",
			},

			"map_proxies": {
				Type:         schema.TypeString,
				Default:      true,
				Optional:     true,
				Description:  "To enable _ disable directs all to the same single pool member",
				ValidateFunc: validateEnabledDisabled,
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileSrcAddrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateSourceAddrPersistenceProfile(
		name,
		parent,
	)

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileSrcAddrUpdate(d, meta)
	if err != nil {
		client.DeleteSourceAddrPersistenceProfile(name)
		return err
	}

	return resourceBigipLtmPersistenceProfileSrcAddrRead(d, meta)

}

func resourceBigipLtmPersistenceProfileSrcAddrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		return err
	}

	d.Set("name", name)
	d.Set("app_service", pp.AppService)
	d.Set("defaults_from", pp.DefaultsFrom)
	d.Set("match_across_pools", pp.MatchAcrossPools)
	d.Set("match_across_services", pp.MatchAcrossServices)
	d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	d.Set("mirror", pp.Mirror)
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	// Specific to SourceAddrPersistenceProfile
	d.Set("hash_algorithm", pp.HashAlgorithm)
	d.Set("map_proxies", pp.MapProxies)
	d.Set("mask", pp.Mask)

	return nil
}

func resourceBigipLtmPersistenceProfileSrcAddrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.SourceAddrPersistenceProfile{
		PersistenceProfile: bigip.PersistenceProfile{
			AppService:              d.Get("app_service").(string),
			DefaultsFrom:            d.Get("defaults_from").(string),
			MatchAcrossPools:        d.Get("match_across_pools").(string),
			MatchAcrossServices:     d.Get("match_across_services").(string),
			MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
			Mirror:                  d.Get("mirror").(string),
			OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			Timeout:                 strconv.Itoa(d.Get("timeout").(int)),
		},

		// Specific to SourceAddrPersistenceProfile
		HashAlgorithm: d.Get("hash_algorithm").(string),
		MapProxies:    d.Get("map_proxies").(string),
		Mask:          d.Get("mask").(string),
	}

	err := client.ModifySourceAddrPersistenceProfile(name, pp)
	if err != nil {
		return err
	}

	return resourceBigipLtmPersistenceProfileSrcAddrRead(d, meta)
}

func resourceBigipLtmPersistenceProfileSrcAddrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Source Address Persistence Profile " + name)

	return client.DeleteSourceAddrPersistenceProfile(name)
}

func resourceBigipLtmPersistenceProfileSrcAddrExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		return false, err
	}

	if pp == nil {
		d.SetId("")
	}

	return pp != nil, nil
}

func resourceBigipPersistenceProfileSrcAddrImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
