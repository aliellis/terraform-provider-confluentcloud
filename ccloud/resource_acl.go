package ccloud

import (
	"context"
	"errors"
	"fmt"
	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func aclResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: aclCreate,
		ReadContext:   aclRead,
		DeleteContext: aclDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "",
			},
			"bootstrap_servers": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL resource type",
			},
			"pattern_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL pattern type",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL name",
			},
			"principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "???",
			},
			"operation": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL operation type: ???",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL host: ???",
			},
			"permission_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL permission type: ???",
			},
		},
	}
}

func aclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	resourceType := d.Get("resource_type").(string)
	patternType := d.Get("pattern_type").(string)
	name := d.Get("name").(string)
	principal := d.Get("principal").(string)
	operation := d.Get("operation").(string)
	host := d.Get("host").(string)
	permissionType := d.Get("permission_type").(string)

	req := ccloud.ACLCreateRequest{
		Pattern: &ccloud.Pattern{
			ResourceType: resourceType,
			PatternType:  patternType,
			Name:         name},
		Entry: &ccloud.Entry{
			Principal:      principal,
			Operation:      operation,
			Host:           host,
			PermissionType: permissionType},
	}
	reqs := ccloud.ACLCreateRequestW{req}

	clusterId := d.Get("cluster_id").(string)
	bootstrapServers := d.Get("bootstrap_servers").(url.URL)
	err := c.CreateACLs(&bootstrapServers, clusterId, &reqs)
	if err == nil {
		// we probably need a unique id here
		d.SetId(fmt.Sprintf("#{name}"))

		// might need to set more attributes to diag?
		err = d.Set("name", name)
		if err != nil {
			diag.FromErr(err)
		}
	} else {
		log.Printf("[ERROR] Could not create ACL: #{err}")
	}

	return diag.FromErr(err)
}

func aclRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	// todo lets try not to use name as ID for the final version
	clusterId := d.Get("cluster_id").(string)
	bootstrapServers := d.Get("bootstrap_servers").(url.URL)
	acl, err := getACL(c, d.Id(), bootstrapServers, clusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("resource_type", acl.Pattern.ResourceType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("pattern_type", acl.Pattern.PatternType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name", acl.Pattern.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("principal", acl.Entry.Principal)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("operation", acl.Entry.Operation)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("host", acl.Entry.Host)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("permission_type", acl.Entry.PermissionType)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getACL(client *ccloud.Client, name string, clusterEndpoint url.URL, clusterID string) (*ccloud.ACL, error) {
	resourceType := "ANY"
	patternType := "LITERAL"
	operation := "ANY"
	host := "*"
	permissionType := "ANY"

	req := ccloud.ACLRequest{
		PatternFilter: &ccloud.ListPatternFilter{ResourceType: resourceType, PatternType: patternType},
		EntryFilter:   &ccloud.ListEntryFilter{Operation: operation, Host: host, PermissionType: permissionType},
	}

	// todo again it would be good to have an id here, name is dangerous
	acls, err := client.ListACLs(&clusterEndpoint, clusterID, &req)
	if err != nil {
		return nil, err
	}

	for _, acl := range acls {
		if acl.Pattern.Name == name {
			return &acl, nil
		}
	}

	return nil, errors.New("unable to find ACL")
}

func aclDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	// todo lets try not to use name as ID for the final version
	ID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not parse ACL ID %s to int", d.Id())
		return diag.FromErr(err)
	}

	resourceType := d.Get("resource_type").(string)
	patternType := d.Get("pattern_type").(string)
	name := d.Get("name").(string)
	principal := d.Get("principal").(string)
	operation := d.Get("operation").(string)
	host := d.Get("host").(string)
	permissionType := d.Get("permission_type").(string)

	req := ccloud.ACLDeleteRequest{
		PatternFilter: &ccloud.DeletePatternFilter{
			ResourceType: resourceType,
			PatternType:  patternType,
			Name:         name},
		EntryFilter: &ccloud.DeleteEntryFilter{
			Principal:      principal,
			Operation:      operation,
			Host:           host,
			PermissionType: permissionType},
	}
	reqs := ccloud.ACLDeleteRequestW{req}

	clusterId := d.Get("cluster_id").(string)
	bootstrapServers := d.Get("bootstrap_servers").(url.URL)
	err = c.DeleteACLs(&bootstrapServers, clusterId, &reqs)
	if err != nil {
		log.Printf("[ERROR] ACL can not be deleted: %d", ID)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] ACL deleted: %d", ID)

	return nil
}
