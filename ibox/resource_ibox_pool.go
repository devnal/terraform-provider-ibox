package ibox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceIboxPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceIboxPoolCreate,
		Read:   resourceIboxPoolRead,
		Update: resourceIboxPoolUpdate,
		Delete: resourceIboxPoolDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Pool name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"virtual_capacity": {
				Description:  "Virtual capacity in bytes",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerGeqThan(1000000000000),
			},
			"physical_capacity": {
				Description:  "Physical capacity in bytes",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerGeqThan(1000000000000),
			},
			"max_extend": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"physical_capacity_critical": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 100),
			},
			"physical_capacity_warning": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 100),
			},
			"ssd_enabled": {
				Description: "Enable/Disable SSD read cache for pool",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"compression_enabled": {
				Description: "Enable/Disable compression for pool",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIboxPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	newPool := Pool{
		Name:                       d.Get("name").(string),
		Max_extend:                 d.Get("max_extend").(int),
		Physical_capacity_critical: d.Get("physical_capacity_critical").(int),
		Physical_capacity_warning:  d.Get("physical_capacity_warning").(int),
		Ssd_enabled:                d.Get("ssd_enabled").(bool),
		Compression_enabled:        d.Get("compression_enabled").(bool),
		Virtual_capacity:           d.Get("virtual_capacity").(int),
		Physical_capacity:          d.Get("physical_capacity").(int),
	}

	// virtual_capacity_size := d.Get("virtual_capacity").(int)
	// if virtual_capacity_size < pool_min_size {
	// 	return fmt.Errorf("[ERROR] Configured virtual_capacity pool size: %v bytes is less than: %v bytes", virtual_capacity_size, pool_min_size)
	// }
	// err := VerifyCapacity(virtual_capacity_size, unit_size)
	// if err != nil {
	// 	return err
	// } else {
	// 	newPool.Virtual_capacity = virtual_capacity_size
	// }

	// physical_capacity_size := d.Get("physical_capacity").(int)
	// if physical_capacity_size < pool_min_size {
	// 	return fmt.Errorf("[ERROR] Configured physical_capacity pool size: %v bytes is less than: %v bytes", physical_capacity_size, pool_min_size)
	// }
	// err = VerifyCapacity(physical_capacity_size, unit_size)
	// if err != nil {
	// 	return err
	// } else {
	// 	newPool.Physical_capacity = physical_capacity_size
	// }

	pool, err := client.CreatePool(newPool)
	if err != nil {
		return err
	} else {
		d.SetId(strconv.Itoa(pool.Id))
		return nil
	}
}

func resourceIboxPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	pool, err := client.ReadPool(d.Id())
	if err != nil {
		return err
	}
	if pool == nil {
		log.Printf("[WARN] Probably the pool was deleted out of band, removing it from state")
		d.SetId("")
		return nil
	}
	d.Set("name", pool.Name)

	return nil
}

func resourceIboxPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	err := client.DeletePool(d.Id())
	if err != nil {
		return err
	}
	return nil
}

func resourceIboxPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	pool_id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	d.Partial(true)

	for k, _ := range resourceIboxPool().Schema {

		var m map[string]interface{}
		m = make(map[string]interface{})

		if d.HasChange(k) {
			old_value, new_value := d.GetChange(k)
			log.Printf("[DEBUG] %v has changed from: %v to: %v", k, old_value, new_value)
			// if k == "physical_capacity" || k == "virtual_capacity" {
			// 	err := VerifyCapacity(d.Get(k).(int), unit_size)
			// 	if err != nil {
			// 		return fmt.Errorf("[ERROR] updating key: %v, error: %v", k, err)
			// 	}
			// }
			m[k] = d.Get(k)

			_, err := client.UpdatePool(m, pool_id)
			if err != nil {
				return err
			}
			d.SetPartial(k)
		}
	}
	d.Partial(false)
	return nil
}
