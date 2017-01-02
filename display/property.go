package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/wallix/awless/cloud/aws"
)

const truncateSize = 25

var (
	// PropertiesDisplayer contains all the display properties of resources
	PropertiesDisplayer = AwlessResourcesDisplayer{
		Services: map[string]*ServiceDisplayer{
			aws.InfraServiceName: &ServiceDisplayer{
				Resources: map[string]*ResourceDisplayer{
					"instance": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "Tags[].Name", Label: "Name"},
							{Property: "State.Name", Label: "State", ColoredValues: map[string]string{"running": "green", "stopped": "red"}},
							{Property: "Type"},
							{Property: "PublicIp"},
							{Property: "PrivateIp"},
						},
					},
					"vpc": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "IsDefault", Label: "Default", ColoredValues: map[string]string{"true": "green"}},
							{Property: "State"},
							{Property: "CidrBlock"},
						},
					},
					"subnet": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "MapPublicIpOnLaunch", Label: "Public VMs", ColoredValues: map[string]string{"true": "red"}},
							{Property: "State", ColoredValues: map[string]string{"available": "green"}},
							{Property: "CidrBlock"},
						},
					},
				},
			},
			aws.AccessServiceName: &ServiceDisplayer{
				Resources: map[string]*ResourceDisplayer{
					"user": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "Name"},
							{Property: "Arn"},
							{Property: "Path"},
							{Property: "PasswordLastUsed"},
						},
					},
					"role": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "Name"},
							{Property: "Arn"},
							{Property: "CreateDate"},
							{Property: "Path"},
						},
					},
					"policy": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "Name"},
							{Property: "Arn"},
							{Property: "Description"},
							{Property: "isAttachable"},
							{Property: "CreateDate"},
							{Property: "UpdateDate"},
							{Property: "Path"},
						},
					},
					"group": &ResourceDisplayer{
						Properties: []*PropertyDisplayer{
							{Property: "Id"},
							{Property: "Name"},
							{Property: "Arn"},
							{Property: "CreateDate"},
							{Property: "Path"},
						},
					},
				},
			},
		},
	}
)

// AwlessResourcesDisplayer contains how to display awless cloud services
type AwlessResourcesDisplayer struct {
	Services map[string]*ServiceDisplayer
}

// ServiceDisplayer contains how to display the resources of a cloud service
type ServiceDisplayer struct {
	Resources map[string]*ResourceDisplayer
}

// ResourceDisplayer contains how to display the properties of a cloud resource
type ResourceDisplayer struct {
	Properties []*PropertyDisplayer
}

// PropertyDisplayer describe how to display a property in a table
type PropertyDisplayer struct {
	Property      string
	Label         string
	ColoredValues map[string]string
	DontTruncate  bool
	TruncateRight bool
}

func (p *PropertyDisplayer) displayName() string {
	if p.Label == "" {
		return p.Property
	}
	return p.Label
}

func (p *PropertyDisplayer) display(value string) string {
	if !p.DontTruncate {
		if p.TruncateRight {
			value = truncateRight(value, truncateSize)
		} else {
			value = truncateLeft(value, truncateSize)
		}
	}
	if p.ColoredValues != nil {
		return colorDisplay(value, p.ColoredValues)
	}
	return value
}

func (p *PropertyDisplayer) displayForceColor(value string, c color.Attribute) string {
	if !p.DontTruncate {
		if p.TruncateRight {
			value = truncateRight(value, truncateSize)
		} else {
			value = truncateLeft(value, truncateSize)
		}
	}
	return color.New(c).SprintFunc()(value)
}

func (p *PropertyDisplayer) propertyValue(properties aws.Properties) string {
	var res string
	if s := strings.SplitN(p.Property, "[].", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = propertyValueSlice(i, s[1])
		}
	} else if s := strings.SplitN(p.Property, "[]length", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = propertyValueSliceLength(i)
		}
	} else if s := strings.SplitN(p.Property, ".", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].(map[string]interface{}); ok {
			res = propertyValueAttribute(i, s[1])
		}
	} else {
		res = propertyValueString(properties[p.Property])
	}
	return res
}

func propertyValueString(prop interface{}) string {
	switch pp := prop.(type) {
	case string:
		return pp
	case bool:
		return fmt.Sprint(pp)
	default:
		return ""
	}
}

func propertyValueSlice(prop []interface{}, key string) string {
	for _, p := range prop {
		//map [key: result]
		if m, ok := p.(map[string]interface{}); ok && m[key] != nil {
			return fmt.Sprint(m[key])
		}

		//map["Key": key, "Value": result]
		if m, ok := p.(map[string]interface{}); ok && m["Key"] == key {
			return fmt.Sprint(m["Value"])
		}
	}
	return ""
}

func propertyValueSliceLength(prop []interface{}) string {
	return strconv.Itoa(len(prop))
}

func propertyValueAttribute(attr map[string]interface{}, key string) string {
	return fmt.Sprint(attr[key])
}

func (p *PropertyDisplayer) firstLevelProperty() string {
	if s := strings.SplitN(p.Property, "[].", 2); len(s) >= 2 {
		return s[0]
	} else if s := strings.SplitN(p.Property, "[]length", 2); len(s) >= 2 {
		return s[0]
	} else if s := strings.SplitN(p.Property, ".", 2); len(s) >= 2 {
		return s[0]
	}
	return p.Property
}