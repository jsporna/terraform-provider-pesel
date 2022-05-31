package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func peselIdData() *schema.Resource {
	return &schema.Resource{
		ReadContext: peselIdDataRead,
		Description: `The ` + "`pesel_id`" + ` data source implements the standard data source lifecycle but does not interact with any external APIs.

Get information about provided PESEL ID like date of birthday, gender. Additionally data source validates checksum of ID.
More information about PESEL ID on [wikipedia](https://en.wikipedia.org/wiki/PESEL).
`,
		Schema: map[string]*schema.Schema{
			"year": {
				Description: "Year of birthday of person described by PESEL ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"month": {
				Description: "Month of birthday of person described by PESEL ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"day": {
				Description: "Day of birthday of person described by PESEL ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"gender": {
				Description: "Gender of person described by PESEL ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"date": {
				Description: "Date of birthday of person described by PESEL ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"male": {
				Description: "Does PESEL ID belong to male person",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"female": {
				Description: "Does PESEL ID belong to female person",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"id": {
				Description: "PESEL ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func peselIdDataRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	pesel := d.Get("id").(string)
	var diags diag.Diagnostics
	for _, c := range pesel {
		if c < '0' || c > '9' {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Provided PESEL ID contains not only digits",
			})
		}
	}
	calculatedChecksum := checksum(pesel)
	peselChecksum := int(pesel[10]) - '0'
	if calculatedChecksum != peselChecksum {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Provided PESEL ID has not valid checksum; pesel: %s, cs: %d, calc: %d", pesel, peselChecksum, calculatedChecksum),
		})
	}
	year, _ := strconv.Atoi(pesel[0:2])
	month, _ := strconv.Atoi(pesel[2:4])
	offset := month - month%20
	calculatedYear := 5*offset + year + YEAR_BASE
	if calculatedYear > YEAR_MAX {
		calculatedYear = calculatedYear - CYCLE_SIZE
	}
	month = month % 20
	day, _ := strconv.Atoi(pesel[4:6])
	genderValue, _ := strconv.Atoi(pesel[9:10])
	_gender := genderValue % 2

	if err := d.Set("year", calculatedYear); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("month", month); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("day", day); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date", fmt.Sprintf("%d-%02d-%02d", calculatedYear, month, day)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("male", _gender == 1); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("female", _gender == 0); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gender", INT2GENDER[_gender]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pesel)
	return nil
}
