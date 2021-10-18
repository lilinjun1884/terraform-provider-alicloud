package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcTrafficMirrorFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcTrafficMirrorFiltersRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Created", "Creating", "Deleting", "Modifying"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"traffic_mirror_filter_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"egress_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_rule_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_rule_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ingress_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_rule_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"traffic_mirror_filter_rule_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcTrafficMirrorFiltersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTrafficMirrorFilters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("traffic_mirror_filter_name"); ok {
		request["TrafficMirrorFilterName"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var trafficMirrorFilterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		trafficMirrorFilterNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_traffic_mirror_filters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TrafficMirrorFilters", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrafficMirrorFilters", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if trafficMirrorFilterNameRegex != nil && !trafficMirrorFilterNameRegex.MatchString(fmt.Sprint(item["TrafficMirrorFilterName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TrafficMirrorFilterId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["TrafficMirrorFilterStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"status":                            object["TrafficMirrorFilterStatus"],
			"traffic_mirror_filter_description": object["TrafficMirrorFilterDescription"],
			"id":                                fmt.Sprint(object["TrafficMirrorFilterId"]),
			"traffic_mirror_filter_id":          fmt.Sprint(object["TrafficMirrorFilterId"]),
			"traffic_mirror_filter_name":        object["TrafficMirrorFilterName"],
		}

		egressRules := make([]map[string]interface{}, 0)
		if egressRulesList, ok := object["EgressRules"].([]interface{}); ok {
			for _, v := range egressRulesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"destination_cidr_block":            m1["DestinationCidrBlock"],
						"destination_port_range":            m1["DestinationPortRange"],
						"priority":                          formatInt(m1["Priority"]),
						"protocol":                          m1["Protocol"],
						"source_cidr_block":                 m1["SourceCidrBlock"],
						"source_port_range":                 m1["SourcePortRange"],
						"traffic_direction":                 m1["TrafficDirection"],
						"traffic_mirror_filter_id":          m1["TrafficMirrorFilterId"],
						"traffic_mirror_filter_rule_id":     m1["TrafficMirrorFilterRuleId"],
						"traffic_mirror_filter_rule_status": m1["TrafficMirrorFilterRuleStatus"],
					}
					egressRules = append(egressRules, temp1)
				}
			}
		}
		mapping["egress_rules"] = egressRules

		ingressRules := make([]map[string]interface{}, 0)
		if ingressRulesList, ok := object["IngressRules"].([]interface{}); ok {
			for _, v := range ingressRulesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"destination_cidr_block":            m1["DestinationCidrBlock"],
						"destination_port_range":            m1["DestinationPortRange"],
						"priority":                          formatInt(m1["Priority"]),
						"protocol":                          m1["Protocol"],
						"source_cidr_block":                 m1["SourceCidrBlock"],
						"source_port_range":                 m1["SourcePortRange"],
						"traffic_direction":                 m1["TrafficDirection"],
						"traffic_mirror_filter_id":          m1["TrafficMirrorFilterId"],
						"traffic_mirror_filter_rule_id":     m1["TrafficMirrorFilterRuleId"],
						"traffic_mirror_filter_rule_status": m1["TrafficMirrorFilterRuleStatus"],
					}
					ingressRules = append(ingressRules, temp1)
				}
			}
		}
		mapping["ingress_rules"] = ingressRules
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TrafficMirrorFilterName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("filters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}