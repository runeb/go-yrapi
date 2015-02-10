package yrapi

import(
  "time"
  "errors"
  "fmt"
  "encoding/xml"
  "net/http"
  "github.com/jonaz/astrotime"
)

type WeatherModel struct {
  Name string `xml:"name,attr"`
  Termin string `xml:"termin,attr"`
  RunEnded string `xml:"runended,attr"`
  NextRun string `xml:"nextrun,attr"`
  From string `xml:"from,attr"`
  To string `xml:"to,attr"`
}

type WeatherMetadata struct {
  Model WeatherModel `xml:"model"`
}

type WeatherProduct struct {
  ProductClass string `xml:"class,attr"`
  WeatherTimes []WeatherTime `xml:"time"`
}

type WeatherData struct {
  Created string `xml:"created,attr"`
  Metadata WeatherMetadata `xml:"meta"`
  Product WeatherProduct `xml:"product"`
}

type WeatherTime struct {
  Type string `xml:"datatype,attr"`
  From string `xml:"from,attr"`
  To string `xml:"to,attr"`

  Location WeatherLocation `xml:"location"`
}

type WeatherLocation struct {
  Altitude int `xml:"altitude,attr"`
  Latitude float64 `xml:"latitude,attr"`
  Longitude float64 `xml:"longitude,attr"`

  Precipitation *WeatherPrecipitation `xml:"precipitation"`
  Symbol *WeatherSymbol `xml:"symbol"`

  Temperature *WeatherTemperature `xml:"temperature,omitempty"`
  WindDirection *WeatherWindDirection `xml:"windDirection,omitempty"`
  WindSpeed *WeatherWindSpeed `xml:"windSpeed,omitempty"`
  Humidity *WeatherHumidity `xml:"humidity,omitempty"`
  Pressure *WeatherPressure `xml:"pressure,omitempty"`
  Cloudiness *WeatherCloudiness `xml:"cloudiness,omitempty"`
  Fog *WeatherFog `xml:"fog,omitempty"`
  LowClouds *WeatherLowClouds `xml:"lowClouds,omitempty"`
  HighClouds *WeatherHighClouds `xml:"highClouds,omitempty"`
  DewpointTemperature *WeatherDewpointTemperature `xml:"dewpointTemperature,omitempty"`
}

type WeatherPrecipitation struct {
  Unit string `xml:"unit,attr"`
  Value float32 `xml:"value,attr"`
  MinValue float32 `xml:"minvalue,attr"`
  MaxValue float32 `xml:"maxvalue,attr"`
}

type WeatherSymbol struct {
  Id string `xml:"id,attr"`
  Number int `xml:"number,attr"`
}

type WeatherDewpointTemperature struct {
  Id string `xml:"id,attr"`
  Unit string `xml:"unit,attr"`
  Value float32 `xml:"value,attr"`
}

type WeatherHighClouds struct {
  Id string `xml:"id,attr"`
  Percent float32 `xml:"percent,attr"`
}

type WeatherLowClouds struct {
  Id string `xml:"id,attr"`
  Percent float32 `xml:"percent,attr"`
}

type WeatherFog struct {
  Id string `xml:"id,attr"`
  Percent float32 `xml:"percent,attr"`
}

type WeatherCloudiness struct {
  Id string `xml:"id,attr"`
  Percent float32 `xml:"percent,attr"`
}

type WeatherPressure struct {
  Id string `xml:"id,attr"`
  Unit string `xml:"unit,attr"`
  Value float32 `xml:"value,attr"`
}

type WeatherTemperature struct {
  Id string `xml:"id,attr"`
  Unit string `xml:"unit,attr"`
  Value float32 `xml:"value,attr"`
}

type WeatherWindDirection struct {
  Id string `xml:"id,attr"`
  Degrees float32 `xml:"deg,attr"`
  Name string `xml:"name,attr"`
}

type WeatherWindSpeed struct {
  Id string `xml:"id,attr"`
  MilesPerSecond float32 `xml:"mps,attr"`
  Beaufort int `xml:"beaufort,attr"`
  Name string `xml:"name,attr"`
}

type WeatherHumidity struct {
  Value float32 `xml:"value,attr"`
  Unit string `xml:"unit,attr"`
}

// Adds a method to WeatherTime returning URL for weather symbol
// corresponding to the forecast at the time and place
func (wt *WeatherTime) SymbolURL() (string, error) {

  // Have to test since Symbol is a struct pointer
  if(wt.Location.Symbol == nil) {
    return "", errors.New("No Location.Symbol")
  }

  // Figure out if its night time at the Location
  sunset := astrotime.CalcSunset(time.Now(), wt.Location.Latitude, wt.Location.Longitude)
  var isNight bool = time.Since(sunset) >= 0

  // TODO: Figure out if the Location is in polar night time
  isPolarnight := false
  var contentType string = "image/svg"

  return WeatherIcon(wt.Location.Symbol.Number, contentType, isNight, isPolarnight), nil
}

// API methods

func WeatherIcon(symbol int, contentType string, isNight bool, isPolarNight bool) string {
  night := 0
  if(isNight) {
    night = 1
  }
  polarNight := 0
  if(isPolarNight) {
    polarNight = 1
  }
  url := fmt.Sprintf(
    "http://api.yr.no/weatherapi/weathericon/1.1/?symbol=%d;is_night=%d;is_polarnight=%d;content_type=%s",
    symbol, night, polarNight, contentType)
  return url
}

func LocationforecastLTS(lat float64, lng float64) (WeatherData, error) {
  // Query the external API
  url := fmt.Sprintf("http://api.yr.no/weatherapi/locationforecastlts/1.2/?lat=%3.5f;lon=%3.5f", lat, lng)
  fmt.Print(url)

  resp, err := http.Get(url)
  if err != nil {
    return WeatherData{}, err
  }

  defer resp.Body.Close()

  // Parse the XML
  var wd WeatherData
  err = xml.NewDecoder(resp.Body).Decode(&wd)
  if err != nil {
    return WeatherData{}, err
  }

  return wd, nil
}

