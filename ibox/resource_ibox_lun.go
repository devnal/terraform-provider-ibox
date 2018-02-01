package ibox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceIboxLun() *schema.Resource {
	return &schema.Resource{
		Create: resourceIboxLunMap,
		Read:   resourceIboxLunQuery,
		Delete: resourceIboxLunUnmap,

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"host_cluster_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"host_id"},
			},
			"lun": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"clustered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIboxLunMap(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newLun := Lun{
		Volume_id: d.Get("volume_id").(int),
		// Host_id:   d.Get("host_id").(int),
		Lun: d.Get("lun").(int),
	}

	if v, ok := d.GetOk("host_cluster_id"); ok == true {
		newLun.Host_cluster_id = v.(int)
	} else if v, ok := d.GetOk("host_id"); ok == true {
		newLun.Host_id = v.(int)
	} else {
		return fmt.Errorf("[ERROR] either host_id or host_cluster_id must be set for lun: %v", newLun)
	}

	lun, err := client.LunMap(newLun)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(lun.Id))
	d.Set("host_id", lun.Host_id)
	d.Set("host_cluster_id", lun.Host_cluster_id)
	d.Set("clustered", lun.Clustered)

	return nil
}

func resourceIboxLunQuery(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newLun := Lun{
		Volume_id:       d.Get("volume_id").(int),
		Host_id:         d.Get("host_id").(int),
		Host_cluster_id: d.Get("host_cluster_id").(int),
		Lun:             d.Get("lun").(int),
		Id:              d.Get("id").(int),
		Clustered:       d.Get("clustered").(bool),
	}

	// host_id := d.Get("host_id").(int)
	// id, _ := strconv.Atoi(d.Id())
	lun, err := client.LunQuery(newLun)
	if err != nil {
		return err
	}
	if lun == nil {
		log.Printf("[WARN] Probably the LUN was deleted out of band, removing it from state")
		d.SetId("")
	}
	return nil
}

func resourceIboxLunUnmap(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newLun := Lun{
		Volume_id:       d.Get("volume_id").(int),
		Host_id:         d.Get("host_id").(int),
		Host_cluster_id: d.Get("host_cluster_id").(int),
		Lun:             d.Get("lun").(int),
		Id:              d.Get("id").(int),
		Clustered:       d.Get("clustered").(bool),
	}

	// host_id := d.Get("host_id").(int)
	// volume_id := d.Get("volume_id").(int)
	err := client.LunUnmap(newLun)
	if err != nil {
		return err
	}
	return nil
}
