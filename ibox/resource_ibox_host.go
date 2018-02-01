package ibox

import (
	// "log"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceIboxHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceIboxHostCreate,
		Read:   resourceIboxHostRead,
		Update: resourceIboxHostUpdate,
		Delete: resourceIboxHostDelete,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateStringInList([]string{
					"NONE",
					"CHAP",
					"MUTUAL_CHAP",
				}, false),
			},
			"security_chap_inbound_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_chap_inbound_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLenghtInRange(14, 255),
			},
			"security_chap_outbound_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_chap_outbound_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLenghtInRange(14, 255),
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ports": {
				Description: "FC or ISCSI port",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description: "IQN for ISCSI or WWN address for FC",
							Optional:    true,
							Type:        schema.TypeString,
						},
						"type": {
							Description: "Port type FC or ISCSI",
							Optional:    true,
							Type:        schema.TypeString,
							ValidateFunc: validateStringInList([]string{
								"FC",
								"ISCSI",
							}, false),
						},
						"host_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIboxHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newHost := Host{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("security_method"); ok {
		newHost.Security_method = v.(string)
	}
	if v, ok := d.GetOk("security_chap_inbound_username"); ok {
		newHost.Security_chap_inbound_username = v.(string)
	}

	if v, ok := d.GetOk("security_chap_inbound_secret"); ok {
		newHost.Security_chap_inbound_secret = v.(string)
	}

	if v, ok := d.GetOk("security_chap_outbound_username"); ok {
		newHost.Security_chap_outbound_username = v.(string)
	}

	if v, ok := d.GetOk("security_chap_outbound_secret"); ok {
		newHost.Security_chap_outbound_secret = v.(string)
	}

	host, err := client.CreateHost(newHost)
	if err != nil {
		return err
	} else {
		d.SetId(strconv.Itoa(host.Id))
		portsRaw := d.Get("ports").([]interface{})
		for _, port := range portsRaw {
			portmap := port.(map[string]interface{})
			portAdd := Port{Address: portmap["address"].(string), Type: portmap["type"].(string)}
			// Trying to add port to the new host, if one of the defined ports cannot be added, the new created host will be deleted.
			_, err := client.CreatePort(portAdd, host.Id)
			if err != nil {
				log.Printf("[ERROR] The new port: %v cannot be added, rolling back changes, deleting host_id: %v", portAdd, host.Id)
				client.DeleteHost(host.Id)
				return err
			}
		}
		resourceIboxHostRead(d, meta)
		return nil
	}
}

func resourceIboxHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	host_id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	host, err := client.ReadHost(host_id)
	if err != nil {
		return err
	}
	if host == nil {
		log.Printf("[WARN] Probably the host was deleted out of band, removing it from state")
		d.SetId("")
		return nil
	}
	portsRawNew := host.Ports
	ports := make([]map[string]interface{}, 0, len(portsRawNew))
	for _, port := range portsRawNew {

		portToSave := make(map[string]interface{})
		portToSave["address"] = port.Address
		portToSave["type"] = port.Type
		portToSave["host_id"] = port.Host_id
		ports = append(ports, portToSave)
	}

	log.Printf("[INFO] READER Setting ports %+v to tfstate", ports)
	err = d.Set("ports", ports)
	if err != nil {
		return fmt.Errorf("[ERROR] Error setting ports: %#v", err)
	}
	d.Set("name", host.Name)

	return nil
}

func resourceIboxHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	host_id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = client.DeleteHost(host_id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIboxHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	host_id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if d.HasChange("name") {
		hostToUpdate := new(Host)
		hostToUpdate.Name = d.Get("name").(string)
		log.Printf("[DEBUG] Host to update: %v", hostToUpdate)
		host, err := client.UpdateHost(*hostToUpdate, host_id)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Updated host: %v", host)
		d.SetId(strconv.Itoa(host.Id))
	}

	if d.HasChange("security_method") || d.HasChange("security_chap_inbound_username") || d.HasChange("security_chap_inbound_secret") || d.HasChange("security_chap_outbound_username") || d.HasChange("security_chap_outbound_secret") {
		var hostToUpdate Host
		_, newv := d.GetChange("security_method")
		if newv == "" {
			hostToUpdate.Security_method = "NONE"
		}
		hostToUpdate.Security_chap_inbound_username = d.Get("security_chap_inbound_username").(string)
		hostToUpdate.Security_chap_inbound_secret = d.Get("security_chap_inbound_secret").(string)
		hostToUpdate.Security_chap_outbound_username = d.Get("security_chap_outbound_username").(string)
		hostToUpdate.Security_chap_outbound_secret = d.Get("security_chap_outbound_secret").(string)

		log.Printf("[DEBUG] Host to update: %v", hostToUpdate)
		_, err := client.UpdateHost(hostToUpdate, host_id)
		if err != nil {
			return err
		}
	}

	if d.HasChange("ports") {
		old_value, new_value := d.GetChange("ports")
		portsRawOld := old_value.([]interface{})
		portsRawNew := new_value.([]interface{})
		for _, port := range portsRawOld {
			portmap := port.(map[string]interface{})
			portToDelete := Port{Address: portmap["address"].(string), Type: portmap["type"].(string)}
			_, err := client.DeletePort(host_id, portToDelete)
			if err != nil {
				return err
			}
		}
		for _, port := range portsRawNew {
			portmap := port.(map[string]interface{})
			portToAdd := Port{Address: portmap["address"].(string), Type: portmap["type"].(string)}
			_, err := client.CreatePort(portToAdd, host_id)
			if err != nil {
				return err
			}
		}
	}
	resourceIboxHostRead(d, meta)
	return nil
}
