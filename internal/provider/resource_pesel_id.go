package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
	"time"
)

func peselIdResource() *schema.Resource {
	return &schema.Resource{
		Description:   ``,
		CreateContext: peselIdResourceWrite,
		DeleteContext: peselIdResourceDelete,
		ReadContext:   peselIdResourceRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"year": {
				Description: "Year of birthday of person described by generated PESEL ID",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"month": {
				Description: "Month of birthday of person described by generated PESEL ID",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"day": {
				Description: "Day of birthday of person described by generated PESEL ID",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"gender": {
				Description: "Gender of person described by generated PESEL ID",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "male",
				ForceNew:    true,
			},
			"date": {
				Description: "Date of birthday of person described by generated PESEL ID",
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
				Description: "Generate PESEL ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func peselIdResourceWrite(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	year := d.Get("year").(int)
	month := d.Get("month").(int)
	day := d.Get("day").(int)
	gender := d.Get("gender").(string)

	if year != 0 && day == 29 && month == 2 && !isLeapYear(year) {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Year %d is not leap so February has only 28 days", year),
		})
	}

	rand.Seed(time.Now().UnixNano())

	var _gender int
	if gender == "" {
		_gender = rand.Intn(2)
	} else {
		_gender = GENDER2INT[gender]
	}

	_year := 0
	if year != 0 {
		_year = year
	} else {
		_year = rand.Intn(YEAR_MAX-YEAR_MIN) + YEAR_MIN
	}
	if _year < YEAR_MIN || _year > YEAR_MAX {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Year should have value betweeen %d & %d", YEAR_MIN, YEAR_MAX),
		})
	}

	_month := 0
	_day := 0
	if day == 0 {
		if month != 0 {
			_month = month
		} else {
			_month = rand.Intn(12) + 1
		}
		if _month < 1 || _month > 12 {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Month should have value between 1 and 12",
			})
		}
		maxDay := lastDayOfMonth(_year, _month)
		_day = rand.Intn(maxDay) + 1
	} else {
		if month == 0 {
			_month = randomMonth(_year, day)
		} else {
			_month = month
		}

		if _month < 1 || _month > 12 {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Month should have value between 1 and 12",
			})
		}

		_day = day
	}
	maxDay := lastDayOfMonth(_year, _month)
	if _day > maxDay || _day < 1 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Day should have value between 1 and %d for month %d in year %d", maxDay, _month, _year),
		})
	}

	_subPesel := fmt.Sprintf("%02d%02d%02d%03d%d",
		_year%100,
		_month+monthOffset(_year),
		_day,
		rand.Intn(1000),
		randomGenderValue(_gender))

	_checksum := checksum(_subPesel)
	_pesel := fmt.Sprintf("%s%d", _subPesel, _checksum)

	if err := d.Set("date", fmt.Sprintf("%d-%02d-%02d", _year, _month, _day)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("male", _gender == 1); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("female", _gender == 0); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(_pesel)

	return diags
}

func peselIdResourceRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func peselIdResourceDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
