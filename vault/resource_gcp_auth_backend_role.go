package vault

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/hashicorp/vault/api"
)

func gcpAuthBackendRoleResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: gcpAuthResourceWrite,
		Update: gcpAuthResourceUpdate,
		Read:   gcpAuthResourceRead,
		Delete: gcpAuthResourceDelete,

		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policies": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"bound_service_accounts": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"backend": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "gcp",
				StateFunc: func(v interface{}) string {
					return strings.Trim(v.(string), "/")
				},
			},
		},
	}
}

func gcpRoleResourcePath(backend, role string) string {
	return "auth/" + strings.Trim(backend, "/") + "/role/" + strings.Trim(role, "/")
}

func gcpAuthResourceWrite(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	backend := d.Get("backend").(string)
	role := d.Get("role").(string)

	path := gcpRoleResourcePath(backend, role)

	data := map[string]interface{}{}

	if v, ok := d.GetOk("type"); ok {
		data["type"] = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		data["project_id"] = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		data["ttl"] = v.(string)
	}

	if v, ok := d.GetOk("max_ttl"); ok {
		data["max_ttl"] = v.(string)
	}

	if v, ok := d.GetOk("period"); ok {
		data["period"] = v.(string)
	}

	if v, ok := d.GetOk("policies"); ok {
		data["policies"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("bound_service_accounts"); ok {
		data["bound_service_accounts"] = v.(*schema.Set).List()
	}

	log.Printf("[DEBUG] Writing role %q to GCP auth backend", path)
	_, err := client.Logical().Write(path, data)
	d.SetId(path)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error writing GCP auth role %q: %s", path, err)
	}
	log.Printf("[DEBUG] Wrote role %q to GCP auth backend", path)

	return gcpAuthResourceRead(d, meta)
}

func gcpAuthResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	path := d.Id()

	data := map[string]interface{}{}

	if v, ok := d.GetOk("ttl"); ok {
		data["ttl"] = v.(string)
	}

	if v, ok := d.GetOk("max_ttl"); ok {
		data["max_ttl"] = v.(string)
	}

	if v, ok := d.GetOk("period"); ok {
		data["period"] = v.(string)
	}

	if v, ok := d.GetOk("policies"); ok {
		data["policies"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("bound_service_accounts"); ok {
		data["bound_service_accounts"] = v.(*schema.Set).List()
	}

	log.Printf("[DEBUG] Updating role %q in GCP auth backend", path)
	_, err := client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("Error updating GCP auth role %q: %s", path, err)
	}
	log.Printf("[DEBUG] Updated role %q to GCP auth backend", path)

	return gcpAuthResourceRead(d, meta)
}

func gcpAuthResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	path := d.Id()

	log.Printf("[DEBUG] Reading GCP role %q", path)
	resp, err := client.Logical().Read(path)
	if err != nil {
		return fmt.Errorf("Error reading GCP role %q: %s", path, err)
	}
	log.Printf("[DEBUG] Read GCP role %q", path)

	if resp == nil {
		log.Printf("[WARN] GCP role %q not found, removing from state", path)
		d.SetId("")
		return nil
	}

	d.Set("ttl", resp.Data["ttl"])
	d.Set("max_ttl", resp.Data["max_ttl"])
	d.Set("type", resp.Data["role_type"])
	d.Set("project_id", resp.Data["project_id"])
	d.Set("period", resp.Data["period"])

	d.Set("policies",
		schema.NewSet(
			schema.HashString, resp.Data["policies"].([]interface{})))

	d.Set("bound_service_accounts",
		schema.NewSet(
			schema.HashString, resp.Data["bound_service_accounts"].([]interface{})))

	return nil
}

func gcpAuthResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	path := d.Id()

	log.Printf("[DEBUG] Deleting GCP role %q", path)
	_, err := client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("Error deleting GCP role %q", path)
	}
	log.Printf("[DEBUG] Deleted GCP role %q", path)

	return nil
}
