package models

import (
	"encoding/json"
	"time"
)

// Vehicle maps 1-to-1 with the vehicle_data table in PostgreSQL.
// We use `db` tags for GORM (raw column names), `json` for ES & HTTP responses.
type Vehicle struct {
	UniqueID                            string          `json:"unique_id" gorm:"column:unique_id;primaryKey"`
	VehicleID                           string          `json:"vehicle_id" gorm:"column:vehicle_id"`
	DataSource                          string          `json:"data_source" gorm:"column:data_source"`
	ListingURL                          string          `json:"listing_url" gorm:"column:listing_url"`
	Images                              json.RawMessage `json:"images,omitempty" gorm:"column:images;type:json"`
	CreatedAt                           *time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	ScrapedAt                           *time.Time      `json:"scraped_at,omitempty" gorm:"column:scraped_at"`
	UpdatedAt                           *time.Time      `json:"updated_at,omitempty" gorm:"column:updated_at"`
	IsVehicleAvailable                  bool            `json:"is_vehicle_available" gorm:"column:is_vehicle_available"`
	Title                               string          `json:"title,omitempty" gorm:"column:title"`
	Price                               string          `json:"price,omitempty" gorm:"column:price"`
	Make                                string          `json:"make,omitempty" gorm:"column:make"`
	Model                               string          `json:"model,omitempty" gorm:"column:model"`
	ModelVersion                        string          `json:"model_version,omitempty" gorm:"column:model_version"`
	ModelRange                          string          `json:"model_range,omitempty" gorm:"column:model_range"`
	PriceInfo                           string          `json:"price_info,omitempty" gorm:"column:price_info"`
	Subtitle                            string          `json:"subtitle,omitempty" gorm:"column:subtitle"`
	PriceLabel                          string          `json:"price_label,omitempty" gorm:"column:price_label"`
	TrimLine                            string          `json:"trim_line,omitempty" gorm:"column:trim_line"`
	VehicleType                         string          `json:"vehicle_type,omitempty" gorm:"column:vehicle_type"`
	Category                            string          `json:"category,omitempty" gorm:"column:category"`
	BodyType                            string          `json:"body_type,omitempty" gorm:"column:body_type"`
	MakeID                              string          `json:"make_id,omitempty" gorm:"column:make_id"`
	ModelID                             string          `json:"model_id,omitempty" gorm:"column:model_id"`
	ModelGenerationID                   string          `json:"model_generation_id,omitempty" gorm:"column:model_generation_id"`
	ModelVariantID                      string          `json:"model_variant_id,omitempty" gorm:"column:model_variant_id"`
	MotorTypeID                         string          `json:"motor_type_id,omitempty" gorm:"column:motor_type_id"`
	TrimLineID                          string          `json:"trim_line_id,omitempty" gorm:"column:trim_line_id"`
	AirConditioningType                 string          `json:"air_conditioning_type,omitempty" gorm:"column:air_conditioning_type"`
	SKU                                 string          `json:"sku,omitempty" gorm:"column:sku"`
	HSNTSN                              string          `json:"hsn_tsn,omitempty" gorm:"column:hsn_tsn"`
	Identifier                          string          `json:"identifier,omitempty" gorm:"column:identifier"`
	PowerKW                             string          `json:"power_kw,omitempty" gorm:"column:power_kw"`
	PowerHP                             string          `json:"power_hp,omitempty" gorm:"column:power_hp"`
	PowerDisplay                        string          `json:"power_display,omitempty" gorm:"column:power_display"`
	PowerKWDisplay                      string          `json:"power_kw_display,omitempty" gorm:"column:power_kw_display"`
	PowerHPDisplay                      string          `json:"power_hp_display,omitempty" gorm:"column:power_hp_display"`
	DisplacementCCM                     string          `json:"displacement_ccm,omitempty" gorm:"column:displacement_ccm"`
	DisplacementDisplay                 string          `json:"displacement_display,omitempty" gorm:"column:displacement_display"`
	Cylinders                           string          `json:"cylinders,omitempty" gorm:"column:cylinders"`
	Gears                               string          `json:"gears,omitempty" gorm:"column:gears"`
	Weight                              string          `json:"weight,omitempty" gorm:"column:weight"`
	NetWeight                           string          `json:"net_weight,omitempty" gorm:"column:net_weight"`
	FuelType                            string          `json:"fuel_type,omitempty" gorm:"column:fuel_type"`
	FuelCategory                        string          `json:"fuel_category,omitempty" gorm:"column:fuel_category"`
	PrimaryFuel                         string          `json:"primary_fuel,omitempty" gorm:"column:primary_fuel"`
	Transmission                        string          `json:"transmission,omitempty" gorm:"column:transmission"`
	TransmissionType                    string          `json:"transmission_type,omitempty" gorm:"column:transmission_type"`
	DriveTrain                          string          `json:"drive_train,omitempty" gorm:"column:drive_train"`
	HasParticleFilter                   string          `json:"has_particle_filter,omitempty" gorm:"column:has_particle_filter"`
	FuelConsumptionCombined             string          `json:"fuel_consumption_combined,omitempty" gorm:"column:fuel_consumption_combined"`
	FuelConsumptionUrban                string          `json:"fuel_consumption_urban,omitempty" gorm:"column:fuel_consumption_urban"`
	FuelConsumptionExtraUrban           string          `json:"fuel_consumption_extra_urban,omitempty" gorm:"column:fuel_consumption_extra_urban"`
	CO2Emission                         string          `json:"co2_emission,omitempty" gorm:"column:co2_emission"`
	CO2EmissionCombined                 string          `json:"co2_emission_combined,omitempty" gorm:"column:co2_emission_combined"`
	CO2EmissionsCombinedFallback        string          `json:"co2_emissions_combined_fallback,omitempty" gorm:"column:co2_emissions_combined_fallback"`
	CO2EmissionsCombinedWeighted        string          `json:"co2_emissions_combined_weighted,omitempty" gorm:"column:co2_emissions_combined_weighted"`
	CO2EmissionsDischarged              string          `json:"co2_emissions_discharged,omitempty" gorm:"column:co2_emissions_discharged"`
	CO2Class                            string          `json:"co2_class,omitempty" gorm:"column:co2_class"`
	CO2ClassDischarged                  string          `json:"co2_class_discharged,omitempty" gorm:"column:co2_class_discharged"`
	EmissionSticker                     string          `json:"emission_sticker,omitempty" gorm:"column:emission_sticker"`
	EmissionStandard                    string          `json:"emission_standard,omitempty" gorm:"column:emission_standard"`
	ConsumptionCombined                 string          `json:"consumption_combined,omitempty" gorm:"column:consumption_combined"`
	ConsumptionCombinedWeighted         string          `json:"consumption_combined_weighted,omitempty" gorm:"column:consumption_combined_weighted"`
	ConsumptionCombinedDischarged       string          `json:"consumption_combined_discharged,omitempty" gorm:"column:consumption_combined_discharged"`
	ConsumptionElectricCombined         string          `json:"consumption_electric_combined,omitempty" gorm:"column:consumption_electric_combined"`
	ConsumptionElectricCombinedWeighted string          `json:"consumption_electric_combined_weighted,omitempty" gorm:"column:consumption_electric_combined_weighted"`
	ConsumptionCity                     string          `json:"consumption_city,omitempty" gorm:"column:consumption_city"`
	ConsumptionCityDischarged           string          `json:"consumption_city_discharged,omitempty" gorm:"column:consumption_city_discharged"`
	ConsumptionSuburban                 string          `json:"consumption_suburban,omitempty" gorm:"column:consumption_suburban"`
	ConsumptionSuburbanDischarged       string          `json:"consumption_suburban_discharged,omitempty" gorm:"column:consumption_suburban_discharged"`
	ConsumptionRural                    string          `json:"consumption_rural,omitempty" gorm:"column:consumption_rural"`
	ConsumptionRuralDischarged          string          `json:"consumption_rural_discharged,omitempty" gorm:"column:consumption_rural_discharged"`
	ConsumptionHighway                  string          `json:"consumption_highway,omitempty" gorm:"column:consumption_highway"`
	ConsumptionHighwayDischarged        string          `json:"consumption_highway_discharged,omitempty" gorm:"column:consumption_highway_discharged"`
	ConsumptionElectricCity             string          `json:"consumption_electric_city,omitempty" gorm:"column:consumption_electric_city"`
	ConsumptionElectricSuburban         string          `json:"consumption_electric_suburban,omitempty" gorm:"column:consumption_electric_suburban"`
	ConsumptionElectricRural            string          `json:"consumption_electric_rural,omitempty" gorm:"column:consumption_electric_rural"`
	ConsumptionElectricHighway          string          `json:"consumption_electric_highway,omitempty" gorm:"column:consumption_electric_highway"`
	EnvironmentPkwEnvkv                 string          `json:"environment_pkw_envkv,omitempty" gorm:"column:environment_pkw_envkv"`
	EnvironmentEUDirective              string          `json:"environment_eu_directive,omitempty" gorm:"column:environment_eu_directive"`
	EnvironmentBimschv35                string          `json:"environment_bimschv35,omitempty" gorm:"column:environment_bimschv35"`
	WLTP                                string          `json:"wltp,omitempty" gorm:"column:wltp"`
	EnergyConsumption                   string          `json:"energy_consumption,omitempty" gorm:"column:energy_consumption"`
	ConsumptionCosts                    string          `json:"consumption_costs,omitempty" gorm:"column:consumption_costs"`
	ConsumptionCostsYear                string          `json:"consumption_costs_year,omitempty" gorm:"column:consumption_costs_year"`
	FuelPrice                           string          `json:"fuel_price,omitempty" gorm:"column:fuel_price"`
	CO2Costs                            string          `json:"co2_costs,omitempty" gorm:"column:co2_costs"`
	CO2CostsAverage                     string          `json:"co2_costs_average,omitempty" gorm:"column:co2_costs_average"`
	CO2CostsHigh                        string          `json:"co2_costs_high,omitempty" gorm:"column:co2_costs_high"`
	CO2CostsLow                         string          `json:"co2_costs_low,omitempty" gorm:"column:co2_costs_low"`
	VehicleTax                          string          `json:"vehicle_tax,omitempty" gorm:"column:vehicle_tax"`
	EngineType                          string          `json:"engine_type,omitempty" gorm:"column:engine_type"`
	OtherEnergySource                   string          `json:"other_energy_source,omitempty" gorm:"column:other_energy_source"`
	MileageKM                           string          `json:"mileage_km,omitempty" gorm:"column:mileage_km"`
	MileageDisplay                      string          `json:"mileage_display,omitempty" gorm:"column:mileage_display"`
	MileageDetail                       string          `json:"mileage_detail,omitempty" gorm:"column:mileage_detail"`
	MileageInKM                         string          `json:"mileage_in_km,omitempty" gorm:"column:mileage_in_km"`
	FirstRegistration                   string          `json:"first_registration,omitempty" gorm:"column:first_registration"`
	FirstRegistrationRaw                string          `json:"first_registration_raw,omitempty" gorm:"column:first_registration_raw"`
	FirstRegistrationDate               string          `json:"first_registration_date,omitempty" gorm:"column:first_registration_date"`
	ProductionYear                      string          `json:"production_year,omitempty" gorm:"column:production_year"`
	ConstructionYear                    string          `json:"construction_year,omitempty" gorm:"column:construction_year"`
	LastInspection                      string          `json:"last_inspection,omitempty" gorm:"column:last_inspection"`
	NextInspection                      string          `json:"next_inspection,omitempty" gorm:"column:next_inspection"`
	LastServiceDate                     string          `json:"last_service_date,omitempty" gorm:"column:last_service_date"`
	LastServiceMileage                  string          `json:"last_service_mileage,omitempty" gorm:"column:last_service_mileage"`
	LastBeltService                     string          `json:"last_belt_service,omitempty" gorm:"column:last_belt_service"`
	FullServiceHistory                  string          `json:"full_service_history,omitempty" gorm:"column:full_service_history"`
	OfferType                           string          `json:"offer_type,omitempty" gorm:"column:offer_type"`
	Condition                           string          `json:"condition,omitempty" gorm:"column:condition"`
	DamageCondition                     string          `json:"damage_condition,omitempty" gorm:"column:damage_condition"`
	HadAccident                         string          `json:"had_accident,omitempty" gorm:"column:had_accident"`
	PreviousOwners                      string          `json:"previous_owners,omitempty" gorm:"column:previous_owners"`
	IsRental                            string          `json:"is_rental,omitempty" gorm:"column:is_rental"`
	NonSmoking                          string          `json:"non_smoking,omitempty" gorm:"column:non_smoking"`
	HasRegistration                     string          `json:"has_registration,omitempty" gorm:"column:has_registration"`
	NewDriverSuitable                   string          `json:"new_driver_suitable,omitempty" gorm:"column:new_driver_suitable"`
	CountryVersion                      string          `json:"country_version,omitempty" gorm:"column:country_version"`
	OriginalMarket                      string          `json:"original_market,omitempty" gorm:"column:original_market"`
	Type                                string          `json:"type,omitempty" gorm:"column:type"`
	Seats                               string          `json:"seats,omitempty" gorm:"column:seats"`
	Doors                               string          `json:"doors,omitempty" gorm:"column:doors"`
	Color                               string          `json:"color,omitempty" gorm:"column:color"`
	ColorOriginal                       string          `json:"color_original,omitempty" gorm:"column:color_original"`
	ManufacturerColor                   string          `json:"manufacturer_color,omitempty" gorm:"column:manufacturer_color"`
	PaintType                           string          `json:"paint_type,omitempty" gorm:"column:paint_type"`
	Upholstery                          string          `json:"upholstery,omitempty" gorm:"column:upholstery"`
	UpholsteryColor                     string          `json:"upholstery_color,omitempty" gorm:"column:upholstery_color"`
	Interior                            string          `json:"interior,omitempty" gorm:"column:interior"`
	InteriorColor                       string          `json:"interior_color,omitempty" gorm:"column:interior_color"`
	InteriorType                        string          `json:"interior_type,omitempty" gorm:"column:interior_type"`
	VehicleModelID                      string          `json:"vehicle_model_id,omitempty" gorm:"column:vehicle_model_id"`
	VehicleTransmission                 string          `json:"vehicle_transmission,omitempty" gorm:"column:vehicle_transmission"`
	VehicleFuelType                     string          `json:"vehicle_fuel_type,omitempty" gorm:"column:vehicle_fuel_type"`
	VehicleFuelConsumption              string          `json:"vehicle_fuel_consumption,omitempty" gorm:"column:vehicle_fuel_consumption"`
	ModelOrLineID                       string          `json:"model_or_line_id,omitempty" gorm:"column:model_or_line_id"`
	WheelBase                           string          `json:"wheel_base,omitempty" gorm:"column:wheel_base"`
	TotalHeight                         string          `json:"total_height,omitempty" gorm:"column:total_height"`
	TotalWidth                          string          `json:"total_width,omitempty" gorm:"column:total_width"`
	TotalLength                         string          `json:"total_length,omitempty" gorm:"column:total_length"`
	GrossVehicleWeight                  string          `json:"gross_vehicle_weight,omitempty" gorm:"column:gross_vehicle_weight"`
	GrossVehicleWeightDetail            string          `json:"gross_vehicle_weight_detail,omitempty" gorm:"column:gross_vehicle_weight_detail"`
	Payload                             string          `json:"payload,omitempty" gorm:"column:payload"`
	LoadWidth                           string          `json:"load_width,omitempty" gorm:"column:load_width"`
	LoadHeight                          string          `json:"load_height,omitempty" gorm:"column:load_height"`
	LoadLength                          string          `json:"load_length,omitempty" gorm:"column:load_length"`
	LoadVolume                          string          `json:"load_volume,omitempty" gorm:"column:load_volume"`
	TrailerLoadBraked                   string          `json:"trailer_load_braked,omitempty" gorm:"column:trailer_load_braked"`
	TrailerLoadUnbraked                 string          `json:"trailer_load_unbraked,omitempty" gorm:"column:trailer_load_unbraked"`
	MaxTowingWeight                     string          `json:"max_towing_weight,omitempty" gorm:"column:max_towing_weight"`
	MaxNoseWeight                       string          `json:"max_nose_weight,omitempty" gorm:"column:max_nose_weight"`
	FuelTankVolume                      string          `json:"fuel_tank_volume,omitempty" gorm:"column:fuel_tank_volume"`
	BatteryOwnership                    string          `json:"battery_ownership,omitempty" gorm:"column:battery_ownership"`
	BatteryChargingTime                 string          `json:"battery_charging_time,omitempty" gorm:"column:battery_charging_time"`
	BatteryCapacity                     string          `json:"battery_capacity,omitempty" gorm:"column:battery_capacity"`
	Battery                             string          `json:"battery,omitempty" gorm:"column:battery"`
	ElectricRange                       string          `json:"electric_range,omitempty" gorm:"column:electric_range"`
	ElectricRangeCity                   string          `json:"electric_range_city,omitempty" gorm:"column:electric_range_city"`
	NumberOfBeds                        string          `json:"number_of_beds,omitempty" gorm:"column:number_of_beds"`
	NumberOfAxles                       string          `json:"number_of_axles,omitempty" gorm:"column:number_of_axles"`
	VehicleArt                          string          `json:"vehicle_art,omitempty" gorm:"column:vehicle_art"`
	CountryCode                         string          `json:"country_code,omitempty" gorm:"column:country_code"`
	PostalCode                          string          `json:"postal_code,omitempty" gorm:"column:postal_code"`
	City                                string          `json:"city,omitempty" gorm:"column:city"`
	Street                              string          `json:"street,omitempty" gorm:"column:street"`
	SellerName                          string          `json:"seller_name,omitempty" gorm:"column:seller_name"`
	LicensePlate                        string          `json:"license_plate,omitempty" gorm:"column:license_plate"`
	Description                         string          `json:"description,omitempty" gorm:"column:description"`

	// Boolean feature flags (subset — most relevant for similarity matching)
	ABS                   *bool `json:"abs,omitempty" gorm:"column:abs"`
	ESP                   *bool `json:"esp,omitempty" gorm:"column:esp"`
	NavigationSystem      *bool `json:"navigation_system,omitempty" gorm:"column:navigation_system"`
	Bluetooth             *bool `json:"bluetooth,omitempty" gorm:"column:bluetooth"`
	Sunroof               *bool `json:"sunroof,omitempty" gorm:"column:sunroof"`
	PanoramicRoof         *bool `json:"panoramic_roof,omitempty" gorm:"column:panoramic_roof"`
	HeatedSeats           *bool `json:"heated_seats,omitempty" gorm:"column:heated_seats"`
	LeatherInterior       *bool `json:"leather_interior,omitempty" gorm:"column:leather_interior"`
	AllWheelDrive         *bool `json:"all_wheel_drive,omitempty" gorm:"column:all_wheel_drive"`
	CruiseControl         *bool `json:"cruise_control,omitempty" gorm:"column:cruise_control"`
	AdaptiveCruiseControl *bool `json:"adaptive_cruise_control,omitempty" gorm:"column:adaptive_cruise_control"`
	ParkingCamera         *bool `json:"parking_camera,omitempty" gorm:"column:parking_camera"`
	LedHeadlights         *bool `json:"led_headlights,omitempty" gorm:"column:led_headlights"`
	TouchScreen           *bool `json:"touchscreen,omitempty" gorm:"column:touchscreen"`
	AppleCarPlay          *bool `json:"apple_carplay,omitempty" gorm:"column:apple_carplay"`
	AndroidAuto           *bool `json:"android_auto,omitempty" gorm:"column:android_auto"`
	StartStopSystem       *bool `json:"start_stop_system,omitempty" gorm:"column:start_stop_system"`
	Warranty              *bool `json:"warranty,omitempty" gorm:"column:warranty"`
}

// TableName tells GORM which table to use
func (Vehicle) TableName() string {
	return "vehicle_marketplace.vehicle_data"
}

// ---------- Search Request / Response ----------

// VehicleSearchRequest is the HTTP request body for POST /api/v1/vehicles/search
type VehicleSearchRequest struct {
	ListingURL string `json:"listing_url" binding:"required"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
}

// VehicleSearchResponse wraps the paginated result
type VehicleSearchResponse struct {
	Total        int64     `json:"total"`
	Page         int       `json:"page"`
	PageSize     int       `json:"page_size"`
	QueryVehicle *Vehicle  `json:"query_vehicle"`
	Results      []Vehicle `json:"results"`
}

// ---------- Indexing Trigger ----------

// IndexTriggerResponse is returned by /api/v1/vehicles/index
type IndexTriggerResponse struct {
	Message   string `json:"message"`
	Scheduled bool   `json:"scheduled"`
}
