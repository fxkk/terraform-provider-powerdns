package powerdns

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePDNSZoneMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourcePDNSZoneMetadataUpdate,
		Read:   resourcePDNSZoneMetadataRead,
		Update: resourcePDNSZoneMetadataUpdate,
		Delete: resourcePDNSZoneMetadataDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
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

func resourcePDNSZoneMetadataUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	zone := d.Get("zone").(string)
	metadata := d.Get("metadata").(*schema.Set)

	var metadataList []ZoneMetadata

	for _, item := range metadata.List() {
		m := item.(map[string]interface{})
		kind := m["kind"].(string)
		valuesRaw := m["values"].([]interface{})
		values := make([]string, len(valuesRaw))
		for i, v := range valuesRaw {
			values[i] = v.(string)
		}

		log.Printf("[DEBUG] Creating metadata for zone %s: kind=%s, values=%v", zone, kind, values)

		metadataList = append(metadataList, ZoneMetadata{
			Kind:     kind,
			Metadata: values,
		})

	}

	err := client.UpdateZoneMetadata(zone, metadataList)
	if err != nil {
		log.Printf("[ERROR] Failed to create metadata for zone %s: %s", zone, err)
		return err
	}

	d.SetId(zone)

	return nil
}

func resourcePDNSZoneMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	zone := d.Id()

	log.Printf("[DEBUG] Reading metadata for zone: %s", zone)
	metadata, err := client.GetZoneMetadata(zone)
	if err != nil {
		return fmt.Errorf("Couldn't fetch metadata for PowerDNS Zone %s: %s", zone, err)
	}

	var metadataList []map[string]interface{}

	for _, m := range d.Get("metadata").(*schema.Set).List() {
		// Convert the State schema.Set item to a map[string]interface{}
		item := m.(map[string]interface{})
		kind := item["kind"].(string)

		// Ony include metadata kinds that are present in the resource block
		for _, md := range metadata {
			if md.Kind == kind {
				log.Printf("[DEBUG] Zone %s metadata: kind=%s, values=%v", zone, md.Kind, md.Metadata)
				metadataList = append(metadataList, map[string]interface{}{
					"kind":   md.Kind,
					"values": md.Metadata,
				})
			} else {
				log.Printf("[DEBUG] Skipping metadata kind %s not present in state", md.Kind)
			}
		}
	}

	d.Set("zone", zone)
	d.Set("metadata", metadataList)

	return nil
}

func resourcePDNSZoneMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting metadata for zone: %s", d.Id())
	client := meta.(*Client)
	zone := d.Get("zone").(string)

	metadata := d.Get("metadata").(*schema.Set)
	var metadataList []ZoneMetadata
	for _, item := range metadata.List() {
		m := item.(map[string]interface{})
		kind := m["kind"].(string)
		valuesRaw := m["values"].([]interface{})
		values := make([]string, len(valuesRaw))
		for i, v := range valuesRaw {
			values[i] = v.(string)
		}
		metadataList = append(metadataList, ZoneMetadata{
			Kind:     kind,
			Metadata: values,
		})
	}

	err := client.DeleteZoneMetadata(zone, metadataList)
	if err != nil {
		log.Printf("[ERROR] Failed to delete metadata for zone %s: %s", zone, err)
		return err
	}

	d.SetId(zone)

	return nil
}
