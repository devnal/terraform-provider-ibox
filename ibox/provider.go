package ibox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IBOX_USERNAME", nil),
				Description: "iBox username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IBOX_PASSWORD", nil),
				Description: "iBox password",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IBOX_HOSTNAME", nil),
				Description: "iBox hostname",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ibox_host_cluster": resourceIboxHostCluster(),
			"ibox_host":         resourceIboxHost(),
			"ibox_pool":         resourceIboxPool(),
			"ibox_volume":       resourceIboxVolume(),
			"ibox_lun":          resourceIboxLun(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: data.Get("username").(string),
		Password: data.Get("password").(string),
		Hostname: data.Get("hostname").(string),
	}

	return config.Client()
}
