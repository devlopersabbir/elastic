import re

with open("table-example.md", "r") as f:
    sql = f.read()

lines = sql.split("\n")

out = """package models

import (
"encoding/json"
"time"
)

// Vehicle maps 1-to-1 with the vehicle_data table in PostgreSQL.
type Vehicle struct {
"""

for line in lines:
    line = line.strip()
    if not line or line.startswith('--') or line.startswith('CREATE') or line.startswith('PRIMARY') or line.startswith(')') or line.startswith('```') or line.startswith('INDEX'):
        continue
    
    parts = line.split(" ", 1)
    if len(parts) < 2: continue
    col_name = parts[0].strip('"')
    if col_name == "PRIMARY" or col_name == ");": continue

    col_type = parts[1].lower()

    if 'json' in col_type:
        go_type = "json.RawMessage"
    elif 'timestamp' in col_type or 'date' in col_type:
        go_type = "*time.Time"
    elif 'boolean' in col_type:
        if col_name == 'is_vehicle_available':
            go_type = "bool"
        else:
            go_type = "*bool"
    else:
        go_type = "string"

    camel = ""
    for s in col_name.split("_"):
        up = s.upper()
        if up in ["ID", "URL", "SKU", "HSN", "TSN", "KM", "KW", "HP", "CCM", "CO2", "EU", "WLTP", "ABS", "ESP", "LED", "USB", "MP3", "DAB", "TV", "CD"]:
            camel += up
        else:
            camel += s.capitalize()

    # Apply special casing adjustments to match existing API correctly without extra fixes.
    if camel == "AppleCarplay":
        camel = "AppleCarPlay"
    if camel == "LEDHeadlights":
        camel = "LedHeadlights"

    tag_json = f'"{col_name},omitempty"'
    tag_gorm = f'"column:{col_name}"'

    if col_name == "unique_id":
        tag_json = '"unique_id"'
        tag_gorm = '"column:unique_id;primaryKey"'
    elif col_name in ["vehicle_id", "data_source", "listing_url", "is_vehicle_available"]:
        tag_json = f'"{col_name}"'
    
    if col_name == "images":
        tag_gorm = '"column:images;type:json"'

    out += f'\t{camel} {go_type} `json:{tag_json} gorm:{tag_gorm}`\n'

out += """}

// TableName tells GORM which table to use
func (Vehicle) TableName() string {
return "vehicle_marketplace.vehicle_data"
}

// ---------- Search Request / Response ----------

type VehicleSearchRequest struct {
ListingURL string `json:"listing_url" binding:"required"`
Page       int    `json:"page,omitempty"`
PageSize   int    `json:"page_size,omitempty"`
}

type VehicleSearchResponse struct {
Total        int64     `json:"total"`
Page         int       `json:"page"`
PageSize     int       `json:"page_size"`
QueryVehicle *Vehicle  `json:"query_vehicle"`
Results      []Vehicle `json:"results"`
}

// ---------- Indexing Trigger ----------

type IndexTriggerResponse struct {
Message   string `json:"message"`
Scheduled bool   `json:"scheduled"`
}
"""

with open("models/vehicle.go", "w") as f:
    f.write(out)

