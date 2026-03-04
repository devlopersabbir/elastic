import re

with open("table-example.md", "r") as f:
    sql = f.read()

# find fields
matches = re.finditer(r'^\s*"?([a-zA-Z0-9_]+)"?\s+([a-zA-Z0-9_\s]+?)(?:,|\))', sql, re.MULTILINE)

print("type Vehicle struct {")
for m in matches:
    col_name = m.group(1)
    if col_name.upper() in ["PRIMARY", "idx_vehicle_data_vehicle_id", "idx_vehicle_data_data_source", "idx_vehicle_data_created_at", "idx_vehicle_data_updated_at", "idx_vehicle_data_scraped_at", "idx_vehicle_data_is_vehicle_available", "idx_vehicle_data_listing_url", "INDEX"]:
        continue
    col_type = m.group(2).strip().lower()

    # Determine Go type
    if 'json' in col_type:
        go_type = 'json.RawMessage'
    elif 'timestamp' in col_type or 'date' in col_type:
        go_type = '*time.Time'
    elif 'boolean' in col_type:
        if col_name == 'is_vehicle_available':
            go_type = 'bool'
        else:
            go_type = '*bool'
    else:
        go_type = 'string'
        
    # Convert col_name to CamelCase for Go field name
    parts = col_name.split('_')
    go_field = ""
    for p in parts:
        if p.lower() == "id": go_field += "ID"
        elif p.lower() == "url": go_field += "URL"
        elif p.lower() == "sku": go_field += "SKU"
        elif p.lower() == "hsn": go_field += "HSN"
        elif p.lower() == "tsn": go_field += "TSN"
        elif p.lower() == "km": go_field += "KM"
        elif p.lower() == "kw": go_field += "KW"
        elif p.lower() == "hp": go_field += "HP"
        elif p.lower() == "ccm": go_field += "CCM"
        elif p.lower() == "co2": go_field += "CO2"
        elif p.lower() == "eu": go_field += "EU"
        elif p.lower() == "wltp": go_field += "WLTP"
        elif p.lower() == "abs": go_field += "ABS"
        elif p.lower() == "esp": go_field += "ESP"
        elif p.lower() == "led": go_field += "LED"
        elif p.lower() == "usb": go_field += "USB"
        elif p.lower() == "mp3": go_field += "MP3"
        elif p.lower() == "dab": go_field += "DAB"
        elif p.lower() == "tv": go_field += "TV"
        elif p.lower() == "cd": go_field += "CD"
        else: go_field += p.capitalize()

    # json and db tags
    tag_json = f'"{col_name},omitempty"'
    tag_gorm = f'"column:{col_name}"'
    if col_name == 'unique_id':
        tag_json = f'"{col_name}"'
        tag_gorm = f'"column:{col_name};primaryKey"'
    elif col_name == 'vehicle_id' or col_name == 'data_source' or col_name == 'listing_url':
        tag_json = f'"{col_name}"'
    elif col_name == 'is_vehicle_available':
        tag_json = f'"{col_name}"'
    elif col_name == 'images':
        tag_gorm = f'"column:{col_name};type:json"'

    print(f'\t{go_field} {go_type} `json:{tag_json} gorm:{tag_gorm}`')

print("}")
