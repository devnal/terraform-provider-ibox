package ibox

import (
	// "fmt"
	// "github.com/adam-hanna/arrayOperations"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceIboxHostCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceIboxHostClusterCreate,
		Read:   resourceIboxHostClusterRead,
		Delete: resourceIboxHostClusterDelete,
		Update: resourceIboxHostClusterUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hosts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIboxHostClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newhostCluster := Host_cluster{
		Name: d.Get("name").(string),
	}

	hostCluster, err := client.CreateHostCluster(newhostCluster)
	if err != nil {
		return err
	}

	hosts := d.Get("hosts").([]interface{})
	for _, host_id := range hosts {
		log.Printf("[DEBUG] configured host_id: %v in host_cluster config", host_id)
		client.AddHostToHostCluster(hostCluster.Id, host_id.(int))
	}

	d.SetId(strconv.Itoa(hostCluster.Id))
	// d.Set("host_id", lun.Host_id)
	// d.Set("host_cluster_id", lun.Host_cluster_id)

	return nil
}

func resourceIboxHostClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	host_cluster_id := d.Get("id").(int)

	host_cluster, err := client.ReadHostCluster(host_cluster_id)
	if err != nil {
		return err
	}
	if host_cluster == nil {
		log.Printf("[WARN] Probably the host cluster was deleted out of band, removing it from state")
		d.SetId("")
	}
	return nil
}

func resourceIboxHostClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	host_cluster_id := d.Get("id").(int)
	err := client.DeleteHostCluster(host_cluster_id)
	if err != nil {
		return err
	}
	return nil
}

func resourceIboxHostClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	// d.Partial(true)

	host_cluster_id := d.Get("id").(int)
	if d.HasChange("hosts") {
		oldv, newv := d.GetChange("hosts")
		// o, n := d.GetChange("hosts.#")
		oldList := oldv.([]interface{})

		for _, host_id := range oldList {
			log.Printf("[INFO] Going to remove the following host id: %v to cluster id: %v", host_id, host_cluster_id)
			_, err := client.RemoveHostFromHostCluster(host_cluster_id, host_id.(int))
			if err != nil {
				return err
			}
		}

		newList := newv.([]interface{})

		for _, host_id := range newList {
			log.Printf("[INFO] Going to add the following host id: %v to cluster id: %v", host_id, host_cluster_id)
			_, err := client.AddHostToHostCluster(host_cluster_id, host_id.(int))
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("name") {
		var m map[string]interface{}
		m = make(map[string]interface{})
		m["name"] = d.Get("name").(string)

		_, err := client.UpdateHostCluster(m, host_cluster_id)
		if err != nil {
			return err
		}

	}
	log.Printf("[DEBUG] %v", client)
	return nil
}
