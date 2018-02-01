package ibox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type Client struct {
	Username string
	Password string
	Hostname string
	Http     *http.Client
	// AuthToken string
}

type ApiError struct {
	Code      string      `json:"code"`
	Data      string      `json:"data"`
	Is_remote bool        `json:"is_remote"`
	Message   string      `json:"message"`
	Reasons   interface{} `json:"reasons"`
	Severity  string      `json:"severity"`
}

type ApiMetadata struct {
	Number_of_objects int  `json:"number_of_objects,omitempty"`
	Page              int  `json:"page,omitempty"`
	Page_size         int  `json:"page_size,omitempty"`
	Pages_total       int  `json:"pages_total,omitempty"`
	Ready             bool `json:"ready,omitempty"`
}

type ApiResult struct {
	Error    *ApiError        `json:"error,omitempty"`
	Metadata *ApiMetadata     `json:"metadata,omitempty"`
	Result   *json.RawMessage `json:"result,omitempty"`
}

type Qos_policy struct {
	Burst_duration_seconds int     `json:"burst_duration_seconds,omitempty"`
	Burst_enabled          bool    `json:"burst_enabled,omitempty"`
	Burst_factor           float32 `json:"burst_factor,omitempty"`
	Id                     int     `json:"id,omitempty"`
	Max_bps                int     `json:"max_bps,omitempty"`
	Max_ops                int     `json:"max_ops,omitempty"`
	Name                   string  `json:"name,omitempty"`
	Type                   string  `json:"type,omitempty"`
}

type Lun struct {
	Clustered       bool `json:"clustered,omitempty"`
	Host_cluster_id int  `json:"host_cluster_id,omitempty"`
	Host_id         int  `json:"host_id,omitempty"`
	Id              int  `json:"id,omitempty"`
	Lun             int  `json:"lun,omitempty"`
	Volume_id       int  `json:"volume_id,omitempty"`
}

type Port struct {
	Address string `json:"address,omitempty"`
	Host_id int    `json:"host_id,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Host struct {
	Id                                int    `json:"id,omitempty"`
	Name                              string `json:"name,omitempty"`
	Host_type                         string `json:"host_type,omitempty"`
	San_client_type                   string `json:"san_client_type,omitempty"`
	Host_cluster_id                   int    `json:"host_cluster_id,omitempty"`
	Security_chap_has_inbound_secret  bool   `json:"security_chap_has_inbound_secret,omitempty"`
	Security_chap_has_outbound_secret bool   `json:"security_chap_has_outbound_secret,omitempty"`
	Security_chap_inbound_username    string `json:"security_chap_inbound_username,omitempty"`
	Security_chap_inbound_secret      string `json:"security_chap_inbound_secret,omitempty"`
	Security_chap_outbound_username   string `json:"security_chap_outbound_username,omitempty"`
	Security_chap_outbound_secret     string `json:"security_chap_outbound_secret,omitempty"`
	Security_method                   string `json:"security_method,omitempty"`
	Luns                              []Lun  `json:"luns,omitempty"`
	Ports                             []Port `json:"ports,omitempty"`
}

type Host_cluster struct {
	San_client_type string `json:"san_client_type,omitempty"`
	Host_type       string `json:"host_type,omitempty"`
	Luns            []Lun  `json:"luns,omitempty"`
	Hosts           []Host `json:"hosts,omitempty"`
	Name            string `json:"name,omitempty"`
	Id              int    `json:"id,omitempty"`
}

type Pool struct {
	Id                          int          `json:"id,omitempty"`
	Name                        string       `json:"name,omitempty"`
	Virtual_capacity            int          `json:"virtual_capacity,omitempty"`
	Physical_capacity           int          `json:"physical_capacity,omitempty"`
	Allocated_physical_capacity int          `json:"allocated_physical_capacity,omitempty"`
	Physical_capacity_critical  int          `json:"physical_capacity_critical,omitempty"`
	Physical_capacity_warning   int          `json:"physical_capacity_warning,omitempty"`
	Reserved_capacity           int          `json:"reserved_capacity,omitempty"`
	Ssd_enabled                 bool         `json:"ssd_enabled,omitempty"`
	Compression_enabled         bool         `json:"compression_enabled,omitempty"`
	Max_extend                  int          `json:"max_extend,omitempty"`
	State                       string       `json:"state,omitempty"`
	Volumes_count               int          `json:"volumes_count,omitempty"`
	Snapshots_count             int          `json:"snapshots_count,omitempty"`
	Filesystems_count           int          `json:"filesystems_count,omitempty"`
	Filesystem_snapshots_count  int          `json:"filesystem_snapshots_count,omitempty"`
	Entities_count              int          `json:"entities_count,omitempty"`
	Owners                      []int        `json:"owners,omitempty"`
	Qos_policies                []Qos_policy `json:"qos_policies,omitempty"`
}

type Volume struct {
	Cg_id                  int    `json:"cg_id,omitempty"`
	Compression_enabled    bool   `json:"compression_enabled,omitempty"`
	Compression_suppressed bool   `json:"compression_suppressed,omitempty"`
	Data_snapshot_guid     string `json:"data_snapshot_guid,omitempty"`
	Dataset_type           string `json:"dataset_type,omitempty"`
	Depth                  int    `json:"depth,omitempty"`
	Family_id              int    `json:"family_id,omitempty"`
	Has_children           bool   `json:"has_children,omitempty"`
	Id                     int    `json:"id,omitempty"`
	Mapped                 bool   `json:"mapped,omitempty"`
	Name                   string `json:"name,omitempty"`
	Num_blocks             int    `json:"num_blocks,omitempty"`
	Parent_id              int    `json:"parent_id,omitempty"`
	Pool_id                int    `json:"pool_id,omitempty"`
	Provtype               string `json:"provtype,omitempty"`
	Qos_policy_id          int    `json:"qos_policy_id,omitempty"`
	Qos_policy_name        string `json:"qos_policy_name,omitempty"`
	Qos_shared_policy_id   int    `json:"qos_shared_policy_id,omitempty"`
	Qos_shared_policy_name string `json:"qos_shared_policy_name,omitempty"`
	Rmr_snapshot_guid      string `json:"rmr_snapshot_guid,omitempty"`
	Rmr_source             bool   `json:"rmr_source,omitempty"`
	Rmr_target             bool   `json:"rmr_target,omitempty"`
	Serial                 string `json:"serial,omitempty"`
	Size                   int    `json:"size,omitempty"`
	Ssd_enabled            bool   `json:"ssd_enabled,omitempty"`
	Tree_allocated         int    `json:"tree_allocated,omitempty"`
	Type                   string `json:"type,omitempty"`
	Used                   int    `json:"used,omitempty"`
	Write_protected        bool   `json:"write_protected,omitempty"`
}

// NewClient returns a new iBox API client
func NewClient(username string, password string, hostname string) (*Client, error) {
	client := Client{
		Username: username,
		Password: password,
		Hostname: hostname,
		Http:     cleanhttp.DefaultClient(),
	}

	return &client, nil
}

// Creates a new request with necessary headers
func (c *Client) newRequest(method string, endpoint string, body []byte) (*http.Request, error) {

	var urlStr string

	urlStr = "http://" + c.Hostname + "/api/rest" + endpoint

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)

	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)

	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	if method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (client *Client) apiCall(method string, endpoint string, data []byte) (*ApiResult, *http.Response, error) {

	req, err := client.newRequest(method, endpoint, data)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] %v", err)
	}

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("[DEBUG] HTTP REQUEST: \n%v", string(requestDump))

	resp, err := client.Http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] %v", err)
	}

	defer resp.Body.Close()

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Errorf("[ERROR] dumping HTTP response: %v", err)
	}
	log.Printf("[DEBUG] HTTP RESPONSE: \n%v", string(responseDump))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] %v", err)
	}

	var apiresult ApiResult
	err = json.Unmarshal(body, &apiresult)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] %v", err)
	}
	return &apiresult, resp, nil
}

func (client *Client) CreateHost(host Host) (*Host, error) {

	reqBody, err := json.MarshalIndent(host, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting host record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/hosts/", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var myhost Host
		json.Unmarshal(*apiresult.Result, &myhost)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		out, _ := json.MarshalIndent(myhost, "", "    ")
		log.Printf("[INFO] Succesfully added new host: %v\n", string(out))
		return &myhost, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to create host record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) ReadHost(host_id int) (*Host, error) {

	apiresult, resp, err := client.apiCall("GET", "/hosts/"+strconv.Itoa(host_id), nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myhost Host
		json.Unmarshal(*apiresult.Result, &myhost)
		log.Printf("[INFO] succesfully fetched host: %v", myhost.Name)
		return &myhost, nil
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] the host with id: %v doesn't exists", host_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}

func (client *Client) DeleteHost(host_id int) error {

	apiresult, resp, err := client.apiCall("DELETE", "/hosts/"+strconv.Itoa(host_id), nil)
	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Succesfully deleted host with id: %v", host_id)
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] The host with id: %v doesn't exists", host_id)
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return fmt.Errorf("[ERROR] %v", string(out))
	}
	return nil
}

func (client *Client) UpdateHost(host Host, host_id int) (*Host, error) {

	reqBody, err := json.MarshalIndent(host, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting host record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("PUT", "/hosts/"+strconv.Itoa(host_id)+"?approved=true", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myhost Host
		json.Unmarshal(*apiresult.Result, &myhost)
		out, _ := json.MarshalIndent(myhost, "", "    ")
		log.Printf("[INFO] Succesfully updated host with id: %v to:\n %v", host_id, string(out))
		return &myhost, nil
	} else {
		out, err := json.MarshalIndent(apiresult.Error, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		return nil, fmt.Errorf("[ERROR] to update host record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) CreatePool(pool Pool) (*Pool, error) {

	reqBody, err := json.MarshalIndent(pool, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting pool record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/pools/", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var mypool Pool
		json.Unmarshal(*apiresult.Result, &mypool)

		out, _ := json.MarshalIndent(mypool, "", "    ")
		log.Printf("[INFO] Succesfully added new pool:\n %v\n", string(out))
		return &mypool, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to create pool record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) ReadPool(pool_id string) (*Pool, error) {

	apiresult, resp, err := client.apiCall("GET", "/pools/"+pool_id, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var mypool Pool
		json.Unmarshal(*apiresult.Result, &mypool)
		log.Printf("[INFO] succesfully fetched volume: %v", mypool.Name)
		return &mypool, nil
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] the volume with id: %v doesn't exists", pool_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}
func (client *Client) DeletePool(pool_id string) error {

	apiresult, resp, err := client.apiCall("DELETE", "/pools/"+pool_id+"?approved=true", nil)
	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Succesfully deleted pool with id: %v", pool_id)
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] The pool with id: %v doesn't exists", pool_id)
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return fmt.Errorf("[ERROR] %v", string(out))
	}
	return nil
}
func (client *Client) UpdatePool(kv map[string]interface{}, pool_id int) (*Pool, error) {

	reqBody, err := json.MarshalIndent(kv, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting pool key/value pair to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("PUT", "/pools/"+strconv.Itoa(pool_id), reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var mypool Pool
		json.Unmarshal(*apiresult.Result, &mypool)
		out, _ := json.MarshalIndent(mypool, "", "    ")
		log.Printf("[INFO] Succesfully updated volume with id: %v to:\n %v", pool_id, string(out))
		return &mypool, nil
	} else {
		out, err := json.MarshalIndent(apiresult.Error, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		return nil, fmt.Errorf("[ERROR] to update pool record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) FindPoolByName(pool_name string) (*Pool, error) {

	apiresult, resp, err := client.apiCall("GET", "/pools?name=eq:"+pool_name, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {

		if apiresult.Metadata.Number_of_objects == 0 {
			return nil, fmt.Errorf("[ERROR] Unable to find pool object name: %v", pool_name)
		}

		var mypools []Pool

		if err := json.Unmarshal(*apiresult.Result, &mypools); err != nil {
			return nil, err
		}
		fmt.Printf("[INFO] succesfully found pool")
		return &mypools[0], nil
	} else {
		return nil, fmt.Errorf("[ERROR]", apiresult.Error)
	}
}

func (client *Client) CreateVolume(volume Volume) (*Volume, error) {

	reqBody, err := json.MarshalIndent(volume, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting volume record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/volumes/", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var myvolume Volume
		json.Unmarshal(*apiresult.Result, &myvolume)

		out, _ := json.MarshalIndent(myvolume, "", "    ")
		log.Printf("[INFO] Succesfully added new volume:\n %v\n", string(out))
		return &myvolume, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to create volume record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) ReadVolume(volume_id string) (*Volume, error) {

	apiresult, resp, err := client.apiCall("GET", "/volumes/"+volume_id, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myvolume Volume
		json.Unmarshal(*apiresult.Result, &myvolume)
		log.Printf("[INFO] succesfully fetched volume: %v", myvolume.Name)
		return &myvolume, nil
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] the volume with id: %v doesn't exists", volume_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}

func (client *Client) DeleteVolume(volume_id string) error {

	apiresult, resp, err := client.apiCall("DELETE", "/volumes/"+volume_id+"?approved=true", nil)
	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Succesfully deleted volume with id: %v", volume_id)
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] The volume with id: %v doesn't exists", volume_id)
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return fmt.Errorf("[ERROR] %v", string(out))
	}
	return nil
}

func (client *Client) UpdateVolume(kv map[string]interface{}, volume_id int) (*Volume, error) {

	reqBody, err := json.MarshalIndent(kv, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting volume key/value pair to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("PUT", "/volumes/"+strconv.Itoa(volume_id), reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myvolume Volume
		json.Unmarshal(*apiresult.Result, &myvolume)
		out, _ := json.MarshalIndent(myvolume, "", "    ")
		log.Printf("[INFO] Succesfully updated volume with id: %v to:\n %v", volume_id, string(out))
		return &myvolume, nil
	} else {
		out, err := json.MarshalIndent(apiresult.Error, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		return nil, fmt.Errorf("[ERROR] to update volume record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) MoveVolume(kv map[string]interface{}, volume_id int) (*Volume, error) {

	reqBody, err := json.MarshalIndent(kv, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting volume key/value pair to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/volumes/"+strconv.Itoa(volume_id)+"/move", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myvolume Volume
		json.Unmarshal(*apiresult.Result, &myvolume)
		out, _ := json.MarshalIndent(myvolume, "", "    ")
		log.Printf("[INFO] Succesfully moved volume id: %v: %v\n", volume_id, string(out))
		return &myvolume, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] to move volume %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) LunMap(lun Lun) (*Lun, error) {

	reqBody, err := json.MarshalIndent(lun, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting lun record to json object: %v", err)
	}

	var url string

	if lun.Host_id != 0 {
		url = "/hosts/" + strconv.Itoa(lun.Host_id) + "/luns"
		fmt.Printf(url)
	} else if lun.Host_cluster_id != 0 {
		url = "/clusters/" + strconv.Itoa(lun.Host_cluster_id) + "/luns"
	} else {
		return nil, fmt.Errorf("[ERROR] either host_id or host cluster_id should be present")
	}

	apiresult, resp, err := client.apiCall("POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var mylun Lun
		json.Unmarshal(*apiresult.Result, &mylun)

		out, _ := json.MarshalIndent(mylun, "", "    ")
		log.Printf("[INFO] Succesfully mapped new lun:\n %v\n to host id: %v\n", string(out), lun.Host_id)
		return &mylun, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to create lun record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) LunQuery(lun Lun) (*Lun, error) {

	var url string
	if lun.Clustered {
		url = "/clusters/" + strconv.Itoa(lun.Host_cluster_id) + "/luns"
		fmt.Printf("[DEBUG] Clustered LUN")
	} else {
		url = "/hosts/" + strconv.Itoa(lun.Host_id) + "/luns"
		fmt.Printf("[DEBUG] Unclustered LUN")
	}

	apiresult, resp, err := client.apiCall("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myluns []Lun
		json.Unmarshal(*apiresult.Result, &myluns)
		log.Printf("[DEBUG] succesfully fetched list of luns for host or custer")
		if apiresult.Metadata.Number_of_objects > 0 {
			log.Printf("[DEBUG] the list is not empty, looking for LUN id: %v", lun.Id)
			for _, mylun := range myluns {
				log.Printf("[DEBUG] mapped LUN: %v", mylun.Id)
				if mylun.Id == lun.Id {
					log.Printf("[INFO] found mapping of LUN id: %v", mylun.Id)
					return &mylun, nil
				}
			}
		}
		log.Printf("[WARN] Unable to find mapping of LUN id: %v", lun.Id)
		return nil, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}

func (client *Client) LunUnmap(lun Lun) error {

	var url string

	if lun.Clustered {
		url = "/clusters/" + strconv.Itoa(lun.Host_cluster_id) + "/luns/volume_id/" + strconv.Itoa(lun.Volume_id) + "?approved=true"
		fmt.Printf("[DEBUG] Clustered LUN")
	} else {
		url = "/hosts/" + strconv.Itoa(lun.Host_id) + "/luns/volume_id/" + strconv.Itoa(lun.Volume_id) + "?approved=true"
		fmt.Printf("[DEBUG] Unclustered LUN")
	}

	apiresult, resp, err := client.apiCall("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}
	if resp.StatusCode == 200 {
		if lun.Clustered {
			log.Printf("[INFO] Succesfully unmapped volume id: %v from host cluster id: %v", lun.Volume_id, lun.Host_cluster_id)
		} else {
			log.Printf("[INFO] Succesfully unmapped volume id: %v from host id: %v", lun.Volume_id, lun.Host_id)
		}
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] Host or Cluster host does not exist")
		return nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return fmt.Errorf("[ERROR] %v", string(out))
	}
	return nil
}

func (client *Client) CreateHostCluster(hostCluster Host_cluster) (*Host_cluster, error) {

	reqBody, err := json.MarshalIndent(hostCluster, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting host cluster record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/clusters/", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var myhostCluster Host_cluster
		json.Unmarshal(*apiresult.Result, &myhostCluster)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		out, _ := json.MarshalIndent(myhostCluster, "", "    ")
		log.Printf("[INFO] Succesfully added new host cluster: %v\n", string(out))
		return &myhostCluster, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to create host cluster record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) ReadHostCluster(host_cluster_id int) (*Host_cluster, error) {

	apiresult, resp, err := client.apiCall("GET", "/clusters/"+strconv.Itoa(host_cluster_id), nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myhostCluster Host_cluster
		json.Unmarshal(*apiresult.Result, &myhostCluster)
		log.Printf("[INFO] succesfully fetched host cluster: %v", myhostCluster.Name)
		return &myhostCluster, nil
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] the host cluster with id: %v doesn't exists", host_cluster_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}

func (client *Client) DeleteHostCluster(host_cluster_id int) error {

	apiresult, resp, err := client.apiCall("DELETE", "/clusters/"+strconv.Itoa(host_cluster_id)+"?approved=true", nil)
	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Succesfully deleted host cluster with id: %v", host_cluster_id)
	} else if resp.StatusCode == 404 {
		log.Printf("[WARN] The host cluster with id: %v doesn't exists", host_cluster_id)
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return fmt.Errorf("[ERROR] %v", string(out))
	}
	return nil
}

func (client *Client) UpdateHostCluster(kv map[string]interface{}, host_cluster_id int) (*Host_cluster, error) {

	reqBody, err := json.MarshalIndent(kv, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting host_cluster key/value pair to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("PUT", "/clusters/"+strconv.Itoa(host_cluster_id), reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myhostCluster Host_cluster
		json.Unmarshal(*apiresult.Result, &myhostCluster)
		out, _ := json.MarshalIndent(myhostCluster, "", "    ")
		log.Printf("[INFO] Succesfully updated host_cluster with id: %v to:\n %v", host_cluster_id, string(out))
		return &myhostCluster, nil
	} else {
		out, err := json.MarshalIndent(apiresult.Error, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		return nil, fmt.Errorf("[ERROR] to update host_cluster record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) AddHostToHostCluster(host_cluster_id int, host_id int) (*Host_cluster, error) {

	var host_id_map map[string]interface{}
	host_id_map = make(map[string]interface{})
	host_id_map["id"] = host_id

	reqBody, err := json.MarshalIndent(host_id_map, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting host_id_map record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/clusters/"+strconv.Itoa(host_cluster_id)+"/hosts", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var myhostCluster Host_cluster
		json.Unmarshal(*apiresult.Result, &myhostCluster)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		out, _ := json.MarshalIndent(myhostCluster, "", "    ")
		log.Printf("[INFO] Succesfully added host_id: %v to host_cluster: %v\n", host_id, string(out))
		return &myhostCluster, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to add host_id: %v to host_cluster_id: %v\n API response: %v", host_id, host_cluster_id, string(out))
	}
}

func (client *Client) RemoveHostFromHostCluster(host_cluster_id int, host_id int) (*Host_cluster, error) {

	apiresult, resp, err := client.apiCall("DELETE", "/clusters/"+strconv.Itoa(host_cluster_id)+"/hosts/"+strconv.Itoa(host_id)+"?approved=true", nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myhostCluster Host_cluster
		json.Unmarshal(*apiresult.Result, &myhostCluster)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		out, _ := json.MarshalIndent(myhostCluster, "", "    ")
		log.Printf("[INFO] Succesfully removed host_id: %v from host_cluster: %v\n", host_id, string(out))
		return &myhostCluster, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to remove host_id: %v from host_cluster_id: %v\n API response: %v", host_id, host_cluster_id, string(out))
	}
}

func (client *Client) CreatePort(port Port, host_id int) (*Port, error) {

	reqBody, err := json.MarshalIndent(port, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Converting port record to json object: %v", err)
	}
	apiresult, resp, err := client.apiCall("POST", "/hosts/"+strconv.Itoa(host_id)+"/ports", reqBody)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 201 {
		var myport Port
		json.Unmarshal(*apiresult.Result, &myport)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] %v", err)
		}
		out, _ := json.MarshalIndent(myport, "", "    ")
		log.Printf("[INFO] Succesfully added new port: %v to host id: %v\n", string(out), host_id)
		return &myport, nil

	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] failed to add port record: %v\n API response: %v", string(reqBody), string(out))
	}
}

func (client *Client) ReadPort(host_id int, port_address string) (*Port, error) {

	apiresult, resp, err := client.apiCall("GET", "/hosts/"+strconv.Itoa(host_id)+"/ports/?address=eq:"+port_address, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myports []Port
		json.Unmarshal(*apiresult.Result, &myports)
		if apiresult.Metadata.Number_of_objects > 0 {
			log.Printf("[INFO] succesfully fetched port: %v", myports[0])
			return &myports[0], nil
		}
		log.Printf("[WARN] Port address: %v was not found in host id: %v", port_address, host_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}

func (client *Client) DeletePort(host_id int, port Port) (*Port, error) {

	apiresult, resp, err := client.apiCall("DELETE", "/hosts/"+strconv.Itoa(host_id)+"/ports/"+port.Type+"/"+port.Address+"?approved=true", nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v", err)
	}

	if resp.StatusCode == 200 {
		var myports []Port
		json.Unmarshal(*apiresult.Result, &myports)
		if apiresult.Metadata.Number_of_objects > 0 {
			log.Printf("[INFO] succesfully deleted port: %v", myports[0])
			return &myports[0], nil
		}
		log.Printf("[WARN] Port address: %v was not found in host id: %v", port.Address, host_id)
		return nil, nil
	} else {
		out, _ := json.MarshalIndent(apiresult.Error, "", "    ")
		return nil, fmt.Errorf("[ERROR] %v", string(out))
	}
}
