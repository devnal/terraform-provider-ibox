package ibox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceIboxVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceIboxVolumeCreate,
		Read:   resourceIboxVolumeRead,
		Update: resourceIboxVolumeUpdate,
		Delete: resourceIboxVolumeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"size": {
				Description: "Volume size in bytes",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"provtype": {
				Description: "Provision type THIN/THICK",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validateStringInList([]string{
					"THIN",
					"THICK",
				}, false),
			},
			"ssd_enabled": {
				Description: "Enable/Disable SSD read cache for volume",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"compression_enabled": {
				Description: "Enable/Disable compression for volume",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceIboxVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newVolume := Volume{
		Name:                d.Get("name").(string),
		Pool_id:             d.Get("pool_id").(int),
		Size:                d.Get("size").(int),
		Provtype:            d.Get("provtype").(string),
		Ssd_enabled:         d.Get("ssd_enabled").(bool),
		Compression_enabled: d.Get("compression_enabled").(bool),
	}

	if newVolume.Size < 1000000000 {
		return fmt.Errorf("[ERROR] Volume size should be at least 1000000000 bytes")
	}

	volume, err := client.CreateVolume(newVolume)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(volume.Id))

	return nil
}

func resourceIboxVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	volume, err := client.ReadVolume(d.Id())
	if err != nil {
		return err
	}
	if volume == nil {
		log.Printf("[WARN] Probably the volume was deleted out of band, removing it from state")
		d.SetId("")
	}
	d.Set("name", volume.Name)
	return nil
}

func resourceIboxVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	err := client.DeleteVolume(d.Id())
	if err != nil {
		return err
	}
	return nil
}

func resourceIboxVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	volume_id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	d.Partial(true)

	for k, _ := range resourceIboxVolume().Schema {

		var m map[string]interface{}
		m = make(map[string]interface{})

		if d.HasChange(k) {
			old_value, new_value := d.GetChange(k)
			log.Printf("[DEBUG] %v has changed from: %v to: %v", k, old_value, new_value)
			if k == "pool_id" {
				m["pool_id"] = d.Get(k).(int)
				m["with_capacity"] = false

				_, err = client.MoveVolume(m, volume_id)
				if err != nil {
					return err
				}
			} else {
				m[k] = d.Get(k)

				_, err := client.UpdateVolume(m, volume_id)
				if err != nil {
					return err
				}
			}
			d.SetPartial(k)
		}
	}
	d.Partial(false)
	return nil
}
