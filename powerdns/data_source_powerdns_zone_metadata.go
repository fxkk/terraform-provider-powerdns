package powerdns

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePowerDNSZoneMetadata() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePowerDNSZoneMetadataRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePowerDNSZoneMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	zone := d.Get("zone").(string)

	// Validate that the zone exists
	exists, err := client.ZoneExists(zone)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("PowerDNS Zone %s does not exist", zone)
	}

	metadata, err := client.GetZoneMetadata(zone)
	if err != nil {
		return fmt.Errorf("error retrieving metadata for zone %s: %s", zone, err)
	}

	log.Printf("[DEBUG] Retrieved metadata for zone %s: %+v", zone, metadata)

	var m []map[string]interface{}

	for _, v := range metadata {
		metadataMap := map[string]interface{}{
			"kind":     v.Kind,
			"metadata": v.Metadata,
		}
		m = append(m, metadataMap)
	}

	err = d.Set("metadata", m)
	if err != nil {
		return fmt.Errorf("error setting metadata for zone %s: %s", zone, err)
	}

	d.SetId(zone)

	return nil
}
