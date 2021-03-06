package bigip

import (
	"log"
	"strconv"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmPersistenceProfileDstAddr() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileDstAddrCreate,
		Read:   resourceBigipLtmPersistenceProfileDstAddrRead,
		Update: resourceBigipLtmPersistenceProfileDstAddrUpdate,
		Delete: resourceBigipLtmPersistenceProfileDstAddrDelete,
		Exists: resourceBigipLtmPersistenceProfileDstAddrExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipPersistenceProfileDstAddrImporter,
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

			// Specific to DestAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Specify the hash algorithm",
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileDstAddrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateDestAddrPersistenceProfile(
		name,
		parent,
	)

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileDstAddrUpdate(d, meta)
	if err != nil {
		client.DeleteDestAddrPersistenceProfile(name)
		return err
	}

	return resourceBigipLtmPersistenceProfileDstAddrRead(d, meta)

}

func resourceBigipLtmPersistenceProfileDstAddrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Destination Address Persistence Profile " + name)

	pp, err := client.GetDestAddrPersistenceProfile(name)
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

	// Specific to DestAddrPersistenceProfile
	d.Set("hash_algorithm", pp.HashAlgorithm)
	d.Set("mask", pp.Mask)

	return nil
}

func resourceBigipLtmPersistenceProfileDstAddrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.DestAddrPersistenceProfile{
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

		// Specific to DestAddrPersistenceProfile
		HashAlgorithm: d.Get("hash_algorithm").(string),
		Mask:          d.Get("mask").(string),
	}

	err := client.ModifyDestAddrPersistenceProfile(name, pp)
	if err != nil {
		return err
	}

	return resourceBigipLtmPersistenceProfileDstAddrRead(d, meta)
}

func resourceBigipLtmPersistenceProfileDstAddrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Destination Address Persistence Profile " + name)

	return client.DeleteDestAddrPersistenceProfile(name)
}

func resourceBigipLtmPersistenceProfileDstAddrExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Destination Address Persistence Profile " + name)

	pp, err := client.GetDestAddrPersistenceProfile(name)
	if err != nil {
		return false, err
	}

	if pp == nil {
		d.SetId("")
	}

	return pp != nil, nil
}

func resourceBigipPersistenceProfileDstAddrImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
