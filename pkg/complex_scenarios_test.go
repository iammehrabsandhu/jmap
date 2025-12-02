package jmap

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestComplexScenarios(t *testing.T) {
	// Scenario 1: FHIR -> EMR.
	t.Run("Healthcare_FHIR_Patient", func(t *testing.T) {
		input := `{
			"resourceType": "Patient",
			"id": "example",
			"identifier": [
				{
					"use": "usual",
					"type": {
						"coding": [
							{
								"system": "http://hl7.org/fhir/v2/0203",
								"code": "MR"
							}
						]
					},
					"system": "urn:oid:1.2.36.146.595.217.0.1",
					"value": "12345"
				}
			],
			"active": true,
			"name": [
				{
					"use": "official",
					"family": "Chalmers",
					"given": [
						"Peter",
						"James"
					]
				}
			],
			"telecom": [
				{
					"system": "phone",
					"value": "(03) 5555 6473",
					"use": "work"
				}
			],
			"gender": "male",
			"birthDate": "1974-12-25",
			"address": [
				{
					"use": "home",
					"line": [
						"534 Erewhon St"
					],
					"city": "PleasantVille",
					"state": "Vic",
					"postalCode": "3999"
				}
			]
		}`
		output := `{
			"patientId": "12345",
			"fullName": "Peter James Chalmers",
			"contact": {
				"phone": "(03) 5555 6473",
				"email": ""
			},
			"demographics": {
				"gender": "male",
				"dob": "1974-12-25",
				"isActive": true
			},
			"primaryAddress": {
				"street": "534 Erewhon St",
				"city": "PleasantVille",
				"state": "Vic",
				"zip": "3999"
			}
		}`
		runScenario(t, input, output)
	})

	// Scenario 2: Shopify -> ERP.
	t.Run("Ecommerce_Shopify_Order", func(t *testing.T) {
		input := `{
			"id": 450789469,
			"email": "bob.norman@hostmail.com",
			"closed_at": null,
			"created_at": "2008-01-10T11:00:00-05:00",
			"updated_at": "2008-01-10T11:00:00-05:00",
			"number": 1,
			"note": null,
			"token": "123456abcd",
			"gateway": "authorize_net",
			"test": false,
			"total_price": "409.94",
			"subtotal_price": "398.00",
			"total_weight": 0,
			"total_tax": "11.94",
			"taxes_included": false,
			"currency": "USD",
			"financial_status": "authorized",
			"confirmed": false,
			"total_discounts": "0.00",
			"total_line_items_price": "398.00",
			"cart_token": "68778783ad298f1c80c3bafcddeea02f",
			"buyer_accepts_marketing": false,
			"name": "#1001",
			"referring_site": "http://www.otherexample.com",
			"landing_site": "http://www.example.com?source=abc",
			"cancelled_at": null,
			"cancel_reason": null,
			"total_price_usd": "409.94",
			"checkout_token": "bd5a8aa1ecd019dd3520ff791ee3a97c",
			"reference": "fhwdgads",
			"user_id": null,
			"location_id": null,
			"source_identifier": null,
			"source_url": null,
			"processed_at": "2008-01-10T11:00:00-05:00",
			"device_id": null,
			"phone": null,
			"customer_locale": "en",
			"app_id": null,
			"browser_ip": null,
			"landing_site_ref": null,
			"order_number": 1001,
			"discount_applications": [],
			"discount_codes": [],
			"note_attributes": [],
			"payment_gateway_names": [
				"authorize_net"
			],
			"processing_method": "direct",
			"checkout_id": 450789469,
			"source_name": "web",
			"fulfillment_status": null,
			"tax_lines": [
				{
					"price": "11.94",
					"rate": 0.06,
					"title": "State Tax"
				}
			],
			"tags": "",
			"contact_email": "bob.norman@hostmail.com",
			"order_status_url": "https://checkout.local/690933842/orders/123456abcd/authenticate?key=690933842",
			"presentment_currency": "USD",
			"total_line_items_price_set": {
				"shop_money": {
					"amount": "398.00",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "398.00",
					"currency_code": "USD"
				}
			},
			"total_discounts_set": {
				"shop_money": {
					"amount": "0.00",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "0.00",
					"currency_code": "USD"
				}
			},
			"total_shipping_price_set": {
				"shop_money": {
					"amount": "0.00",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "0.00",
					"currency_code": "USD"
				}
			},
			"subtotal_price_set": {
				"shop_money": {
					"amount": "398.00",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "398.00",
					"currency_code": "USD"
				}
			},
			"total_price_set": {
				"shop_money": {
					"amount": "409.94",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "409.94",
					"currency_code": "USD"
				}
			},
			"total_tax_set": {
				"shop_money": {
					"amount": "11.94",
					"currency_code": "USD"
				},
				"presentment_money": {
					"amount": "11.94",
					"currency_code": "USD"
				}
			},
			"line_items": [
				{
					"id": 466157049,
					"variant_id": 39072856,
					"title": "IPod Nano - 8gb",
					"quantity": 1,
					"sku": "IPOD2008GREEN",
					"variant_title": "green",
					"vendor": null,
					"fulfillment_service": "manual",
					"product_id": 632910392,
					"requires_shipping": true,
					"taxable": true,
					"gift_card": false,
					"name": "IPod Nano - 8gb - green",
					"variant_inventory_management": "shopify",
					"properties": [],
					"product_exists": true,
					"fulfillable_quantity": 1,
					"grams": 200,
					"price": "199.00",
					"total_discount": "0.00",
					"fulfillment_status": null,
					"price_set": {
						"shop_money": {
							"amount": "199.00",
							"currency_code": "USD"
						},
						"presentment_money": {
							"amount": "199.00",
							"currency_code": "USD"
						}
					},
					"total_discount_set": {
						"shop_money": {
							"amount": "0.00",
							"currency_code": "USD"
						},
						"presentment_money": {
							"amount": "0.00",
							"currency_code": "USD"
						}
					},
					"discount_allocations": [],
					"duties": [],
					"admin_graphql_api_id": "gid://shopify/LineItem/466157049",
					"tax_lines": [
						{
							"title": "State Tax",
							"price": "3.98",
							"rate": 0.06
						}
					]
				}
			],
			"shipping_lines": [
				{
					"id": 369256396,
					"title": "Free Shipping",
					"price": "0.00",
					"code": "Free Shipping",
					"source": "shopify",
					"phone": null,
					"requested_fulfillment_service_id": null,
					"delivery_category": null,
					"carrier_identifier": null,
					"discounted_price": "0.00",
					"price_set": {
						"shop_money": {
							"amount": "0.00",
							"currency_code": "USD"
						},
						"presentment_money": {
							"amount": "0.00",
							"currency_code": "USD"
						}
					},
					"discounted_price_set": {
						"shop_money": {
							"amount": "0.00",
							"currency_code": "USD"
						},
						"presentment_money": {
							"amount": "0.00",
							"currency_code": "USD"
						}
					},
					"discount_allocations": [],
					"tax_lines": []
				}
			]
		}`
		output := `{
			"erpOrderId": "1001",
			"customerEmail": "bob.norman@hostmail.com",
			"orderDate": "2008-01-10T11:00:00-05:00",
			"totals": {
				"grandTotal": 409.94,
				"tax": 11.94,
				"currency": "USD"
			},
			"items": [
				{
					"sku": "IPOD2008GREEN",
					"description": "IPod Nano - 8gb - green",
					"qty": 1,
					"unitPrice": 199.00
				}
			],
			"shipping": {
				"method": "Free Shipping",
				"cost": 0.00
			}
		}`
		runScenario(t, input, output)
	})

	// Scenario 3: SWIFT -> ISO 20022.
	t.Run("Finance_SWIFT_ISO20022", func(t *testing.T) {
		input := `{
			"block1": "{1:F01BANKBEBBAXXX2039063581}",
			"block2": "{2:I103BANKDEFFXXXXN}",
			"block4": {
				"20": "REFERENCE123",
				"23B": "CRED",
				"32A": "090901EUR100000,",
				"50K": "/12345678\nJOHN DOE\nSTREET 1\nCITY",
				"59": "/87654321\nJANE SMITH\nAVENUE 2\nTOWN",
				"70": "INVOICE 999"
			}
		}`
		output := `{
			"AppHdr": {
				"Fr": { "BICFI": "BANKBEBBAXXX" },
				"To": { "BICFI": "BANKDEFFXXXX" },
				"MsgDefIdr": "pacs.008.001.08"
			},
			"Document": {
				"FIToFICstmrCdtTrf": {
					"GrpHdr": {
						"MsgId": "REFERENCE123",
						"CreDtTm": "2009-09-01T12:00:00Z",
						"NbOfTx": "1"
					},
					"CdtTrfTxInf": [
						{
							"PmtId": { "EndToEndId": "REFERENCE123" },
							"IntrBkSttlmAmt": { "Ccy": "EUR", "Value": 100000 },
							"Dbtr": { "Nm": "JOHN DOE", "PstlAdr": { "StrtNm": "STREET 1", "TwnNm": "CITY" } },
							"Cdtr": { "Nm": "JANE SMITH", "PstlAdr": { "StrtNm": "AVENUE 2", "TwnNm": "TOWN" } },
							"RmtInf": { "Ustrd": "INVOICE 999" }
						}
					]
				}
			}
		}`
		runScenario(t, input, output)
	})

	// Scenario 4: Logistics.
	t.Run("Logistics_Shipment", func(t *testing.T) {
		input := `{
			"tracking_number": "1Z999AA10123456784",
			"carrier": "UPS",
			"status": "In Transit",
			"scan_history": [
				{
					"location": "Louisville, KY",
					"timestamp": "2023-10-25T08:00:00Z",
					"activity": "Departure Scan"
				},
				{
					"location": "Nashville, TN",
					"timestamp": "2023-10-25T14:00:00Z",
					"activity": "Arrival Scan"
				}
			],
			"estimated_delivery": "2023-10-27"
		}`
		output := `{
			"shipmentId": "1Z999AA10123456784",
			"provider": "UPS",
			"currentStatus": "In Transit",
			"lastLocation": "Nashville, TN",
			"lastUpdate": "2023-10-25T14:00:00Z",
			"eta": "2023-10-27",
			"events": [
				{
					"city": "Louisville, KY",
					"time": "2023-10-25T08:00:00Z",
					"type": "Departure Scan"
				},
				{
					"city": "Nashville, TN",
					"time": "2023-10-25T14:00:00Z",
					"type": "Arrival Scan"
				}
			]
		}`
		runScenario(t, input, output)
	})

	// Scenario 5: Workday -> AD.
	t.Run("HR_Workday_AD", func(t *testing.T) {
		input := `{
			"Worker_Data": {
				"Worker_ID": "1001",
				"Personal_Data": {
					"Name_Data": {
						"Legal_Name_Data": {
							"Name_Detail_Data": {
								"First_Name": "John",
								"Last_Name": "Doe"
							}
						}
					},
					"Contact_Data": {
						"Email_Address_Data": [
							{
								"Email_Address": "john.doe@example.com",
								"Usage_Data": { "Type_Data": { "Type_Reference": { "ID": "WORK" } } }
							}
						]
					}
				},
				"Employment_Data": {
					"Job_Data": {
						"Job_Title": "Software Engineer",
						"Supervisory_Organization": "Engineering"
					}
				}
			}
		}`
		output := `{
			"sAMAccountName": "jdoe",
			"givenName": "John",
			"sn": "Doe",
			"displayName": "John Doe",
			"mail": "john.doe@example.com",
			"title": "Software Engineer",
			"department": "Engineering",
			"employeeID": "1001",
			"userPrincipalName": "john.doe@example.com"
		}`
		runScenario(t, input, output)
	})

	// Scenario 6: GDS -> Itinerary.
	t.Run("Travel_GDS_Itinerary", func(t *testing.T) {
		input := `{
			"PNR": "ABC123",
			"segments": [
				{
					"flightNumber": "UA100",
					"airline": "UA",
					"departure": {
						"airport": "SFO",
						"time": "2023-11-01T08:00:00"
					},
					"arrival": {
						"airport": "JFK",
						"time": "2023-11-01T16:30:00"
					},
					"class": "Economy"
				}
			],
			"passengers": [
				{
					"firstName": "Alice",
					"lastName": "Smith",
					"ticketNumber": "0161234567890"
				}
			]
		}`
		output := `{
			"bookingReference": "ABC123",
			"flights": [
				{
					"code": "UA100",
					"carrier": "United Airlines",
					"origin": "SFO",
					"destination": "JFK",
					"departs": "2023-11-01T08:00:00",
					"arrives": "2023-11-01T16:30:00"
				}
			],
			"travelers": [
				{
					"name": "Alice Smith",
					"ticket": "0161234567890"
				}
			]
		}`
		runScenario(t, input, output)
	})

	// Scenario 7: IoT Telemetry.
	t.Run("IoT_Telemetry", func(t *testing.T) {
		input := `{
			"deviceId": "sensor-001",
			"timestamp": 1698240000,
			"readings": {
				"temperature": {
					"value": 22.5,
					"unit": "C"
				},
				"humidity": {
					"value": 45,
					"unit": "%"
				},
				"battery": {
					"voltage": 3.7,
					"percentage": 85
				}
			},
			"metadata": {
				"firmware": "1.2.3",
				"location": {
					"lat": 37.7749,
					"lon": -122.4194
				}
			}
		}`
		output := `{
			"id": "sensor-001",
			"ts": 1698240000,
			"temp_c": 22.5,
			"humidity_pct": 45,
			"batt_volts": 3.7,
			"batt_pct": 85,
			"fw_version": "1.2.3",
			"latitude": 37.7749,
			"longitude": -122.4194
		}`
		runScenario(t, input, output)
	})

	// Scenario 8: FB -> Profile.
	t.Run("Social_FB_Profile", func(t *testing.T) {
		input := `{
			"id": "123456789",
			"name": "John Doe",
			"first_name": "John",
			"last_name": "Doe",
			"email": "john@example.com",
			"picture": {
				"data": {
					"height": 50,
					"is_silhouette": false,
					"url": "https://platform-lookaside.fbsbx.com/platform/profilepic/?asid=123456789&height=50&width=50&ext=1698240000&hash=AeQ...",
					"width": 50
				}
			},
			"friends": {
				"data": [
					{ "name": "Jane Smith", "id": "987654321" }
				],
				"summary": { "total_count": 500 }
			}
		}`
		output := `{
			"userId": "123456789",
			"displayName": "John Doe",
			"contactEmail": "john@example.com",
			"avatarUrl": "https://platform-lookaside.fbsbx.com/platform/profilepic/?asid=123456789&height=50&width=50&ext=1698240000&hash=AeQ...",
			"friendCount": 500,
			"topFriends": [
				{ "name": "Jane Smith", "uid": "987654321" }
			]
		}`
		runScenario(t, input, output)
	})

	// Scenario 9: Canvas -> SIS.
	t.Run("Education_Canvas_SIS", func(t *testing.T) {
		input := `{
			"id": 101,
			"name": "Introduction to Computer Science",
			"course_code": "CS101",
			"start_at": "2023-09-01T00:00:00Z",
			"end_at": "2023-12-15T00:00:00Z",
			"enrollments": [
				{
					"type": "student",
					"user": {
						"id": 1,
						"name": "Alice",
						"login_id": "alice123"
					},
					"grades": {
						"current_score": 95.5,
						"final_score": 95.5,
						"current_grade": "A"
					}
				}
			],
			"teachers": [
				{
					"id": 5,
					"display_name": "Dr. Smith",
					"email": "smith@university.edu"
				}
			]
		}`
		output := `{
			"sis_course_id": "CS101_FALL_2023",
			"title": "Introduction to Computer Science",
			"term": "Fall 2023",
			"instructor": {
				"name": "Dr. Smith",
				"email": "smith@university.edu"
			},
			"roster": [
				{
					"student_id": "alice123",
					"name": "Alice",
					"grade": "A",
					"score": 95.5
				}
			]
		}`
		runScenario(t, input, output)
	})

	// Scenario 10: Salesforce -> HubSpot.
	t.Run("CRM_SFDC_HubSpot", func(t *testing.T) {
		input := `{
			"Id": "0012E00001Z8q9xQAB",
			"Name": "Acme Corp",
			"BillingStreet": "123 Main St",
			"BillingCity": "San Francisco",
			"BillingState": "CA",
			"BillingPostalCode": "94105",
			"BillingCountry": "USA",
			"Phone": "(415) 555-1234",
			"Website": "https://www.acme.com",
			"Industry": "Technology",
			"AnnualRevenue": 10000000,
			"NumberOfEmployees": 500,
			"Owner": {
				"Name": "Sales Rep",
				"Email": "sales@acme.com"
			},
			"Contacts": {
				"records": [
					{
						"FirstName": "John",
						"LastName": "Doe",
						"Email": "john@acme.com",
						"Title": "CEO"
					}
				]
			}
		}`
		output := `{
			"properties": {
				"name": "Acme Corp",
				"address": "123 Main St",
				"city": "San Francisco",
				"state": "CA",
				"zip": "94105",
				"country": "USA",
				"phone": "(415) 555-1234",
				"website": "https://www.acme.com",
				"industry": "Technology",
				"annual_revenue": 10000000,
				"number_of_employees": 500,
				"hubspot_owner_id": "sales@acme.com"
			},
			"associations": {
				"contacts": [
					{
						"email": "john@acme.com",
						"firstname": "John",
						"lastname": "Doe",
						"jobtitle": "CEO"
					}
				]
			}
		}`
		runScenario(t, input, output)
	})
}

func runScenario(t *testing.T, inputJSON, outputJSON string) {
	spec, err := SuggestSpec(inputJSON, outputJSON)
	if err != nil {
		t.Fatalf("SuggestSpec failed: %v", err)
	}

	// 1. Parse Output.
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(outputJSON), &output); err != nil {
		t.Fatalf("Invalid output JSON: %v", err)
	}
	outputFields := flatten(output, "")

	// 2. Collect Targets.
	shiftTargets := make(map[string]bool)
	defaultKeys := make(map[string]bool)

	for _, op := range spec.Operations {
		switch op.Type {
		case "shift":
			collectShiftTargets(op.Spec, shiftTargets)
		case "default":
			collectDefaultKeys(op.Spec, "", defaultKeys)
		}
	}

	shiftCount := 0
	defaultCount := 0
	unaccountedCount := 0

	for _, field := range outputFields {
		schemaPath := normalizePath(field)

		// Check coverage.
		// We check exact match or wildcard match
		if isCovered(schemaPath, shiftTargets) {
			shiftCount++
		} else if isCovered(schemaPath, defaultKeys) {
			defaultCount++
		} else {
			unaccountedCount++
		}
	}

	total := shiftCount + defaultCount + unaccountedCount
	t.Logf("Field Stats - Total: %d | Shift: %d | Default: %d | Unaccounted: %d",
		total, shiftCount, defaultCount, unaccountedCount)

	if total > 0 {
		ratio := float64(shiftCount) / float64(total)
		t.Logf("Mapping Success Rate (by field): %.1f%%", ratio*100)
	}
}

// Helpers

func flatten(data interface{}, prefix string) []string {
	var fields []string
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			fields = append(fields, flatten(val, p)...)
		}
	case []interface{}:
		for i, val := range v {
			p := fmt.Sprintf("%s[%d]", prefix, i)
			fields = append(fields, flatten(val, p)...)
		}
	default:
		if prefix != "" {
			fields = append(fields, prefix)
		}
	}
	return fields
}

func collectShiftTargets(spec interface{}, targets map[string]bool) {
	switch v := spec.(type) {
	case map[string]interface{}:
		for _, val := range v {
			collectShiftTargets(val, targets)
		}
	case string:
		targets[v] = true
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				targets[s] = true
			}
		}
	}
}

func collectDefaultKeys(spec interface{}, prefix string, keys map[string]bool) {
	if m, ok := spec.(map[string]interface{}); ok {
		for k, v := range m {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			// Recurse or leaf.
			if nested, ok := v.(map[string]interface{}); ok {
				collectDefaultKeys(nested, p, keys)
			} else {
				keys[p] = true
			}
		}
	}
}

func normalizePath(path string) string {
	// Remove array indices: items[0].id -> items.id
	// We use a regex-like approach or simple string building
	var sb strings.Builder
	inBracket := false
	for _, r := range path {
		if r == '[' {
			inBracket = true
			continue
		}
		if r == ']' {
			inBracket = false
			continue
		}
		if !inBracket {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func isCovered(path string, targets map[string]bool) bool {
	// 1. Exact match
	if targets[path] {
		return true
	}
	// 2. Wildcard match (simple version: if target has *, check if path matches pattern)
	// For now, SuggestSpec generates explicit paths or * wildcards.
	// If the spec has "items.*.id", and path is "items.id" (normalized), they match.
	// Wait, normalizePath removes indices, so "items[0].id" becomes "items.id".
	// The spec target might be "items.id" (if array was flattened) or "items[&].id".
	// JMap's SuggestSpec currently generates "items" -> "*" -> "id" which results in target "items.&.id" or similar?
	// Actually SuggestSpec generates "items" -> "*" -> "id" : "sourcePath".
	// The TARGET path in shift spec is the value.
	// In SuggestSpec, the value is the SOURCE path.
	// Wait, shift spec is: "source": "target".
	// SuggestSpec generates: "inputPath": "outputPath".
	// So we need to collect VALUES from shift spec.

	// If SuggestSpec generates "orders": { "*": { "id": "orderId" } },
	// then for input "orders[0].id", it maps to "orderId".
	// The target is "orderId".
	// Our output field is "orderId". Normalized: "orderId". Match!

	// If output is "items[0].id", normalized "items.id".
	// Spec might say "items": { "*": { "id": "items[&].id" } } ?
	// No, SuggestSpec usually tries to map to specific fields.
	// If the output has an array, SuggestSpec might generate "items[&].id".
	// Let's assume exact match on normalized path for now.
	// If the target contains "&", we might need to be smarter.
	// For this test, let's strip "&" from targets too?

	for t := range targets {
		// Normalize target path too (remove indices)
		normTarget := normalizePath(t)
		normTarget = strings.ReplaceAll(normTarget, "&", "")
		normTarget = strings.ReplaceAll(normTarget, "[]", "") // handle [] syntax if any
		if normTarget == path {
			return true
		}
	}
	return false
}
