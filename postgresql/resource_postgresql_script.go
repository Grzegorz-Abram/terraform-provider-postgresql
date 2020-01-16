package postgresql

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	sqlScript = "sql_script"
)

func resourcePostgreSQLScript() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgreSQLRoleCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			sqlScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL script",
			},
		},
	}
}

func resourcePostgreSQLRoleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	c.catalogLock.Lock()
	defer c.catalogLock.Unlock()

	txn, err := c.DB().Begin()
	if err != nil {
		return err
	}
	defer deferredRollback(txn)

	script := d.Get(sqlScript).(string)

	sql := fmt.Sprintf("%s", script)
	if _, err := txn.Exec(sql); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error running script %s: {{err}}", script), err)
	}

	if err = txn.Commit(); err != nil {
		return errwrap.Wrapf("could not commit transaction: {{err}}", err)
	}

	d.SetId(script)

	return nil
}
