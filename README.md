# Command line weather reports

## Introduction

I love aviation and one of the things I especially enjoy about it is the concise format of the weather reports used by the airmen. Being the UAV operator, I frequently find myself in need of an accurate weather info and I believe there is a benefit in having a light alternative to cumbersome weather services. Therefore I have created an app which downloads the weather reports issued by user-specified airports. It uses [Aviation Weather Center](https://aviationweather.gov)'s website as the source of data. 

### METAR

**Met**eorological **A**erodrome **R**eport is the data about current weather compressed into standardised format. It is issued by airports every set interval of time and typically contains information about wind, temperature, visibility, clouds, pressure and special conditions (fog, precipitation etc.). It is worth noting, that north american reports differ from more commonly used international format. See the [Resources](#resources) section for more information.

Example METAR depeche for Warsaw Airport, issued on 6th day of the month 18:30 Zulu (UTC):

```
METAR EPWA 061830Z 23007KT 200V260 9999 SCT034 BKN039 09/05 Q1023 NOSIG
```

Interpretation

Wind direction 230&deg; (varying between 200&deg; and 260&deg;), speed 7kt. Visibility at or above 10000m. Scattered clouds (3/8 - 4/8 of the sky covered) at 3400ft, broken ceiling (5/8 - 7/8 of the sky covered) at 3900ft. Temperature 9&deg;C, dew point of 5&deg;C. Pressure at 1023hPa. No significant change expected within 2h.

### TAF

**T**erminal **A**erodrome **F**orecast is the coded format of the weather prediction, typically for a period of 24 hours. It uses the same condition codes as METAR along with some specific keywords that express temporal aspects of the predicted weather changes. It is usually issued 4 times a day.

Example TAF depeche for Warsaw Airport, issued on 6th day of the month 17:30 Zulu (UTC):

```
TAF EPWA 061730Z 0618/0718 23010KT 9999 BKN040
  TEMPO 0711/0718 21015G25KT -SHRA BKN012CB
  BECMG 0712/0714 -RA BKN014 BKN030
  BECMG 0714/0716 SCT004 BKN007
```

Interpretation

Forecast for 6th day of the month 18:00 to 7th 18:00. Wind direction 230&deg;, speed 10kt. Visibility at or above 10km. Broken ceiling (5/8 - 7/8 coverage) at 4000ft. <br>
Temporary conditions on 7th 11:00 to 18:00: Wind speed of 15kt gusting up to 25kt at 210&deg;, weak rain showers. Broken ceiling at 1200ft with cumulonimbus cloud in airport's vicinity. <br>
Conditions expected to develop on 7th at 14:00: weak rain. Broken ceiling at 1400ft and breaking clouds deck at 3000ft. <br>
Conditions expected to develop on 7th at 16:00: scattered clouds (3/8 - 4/8 coverage) at 400ft and broken ceiling at 700ft.

## Application

### Installation

You can download and unpack the zipped release into the directory of your choice. I would, however, encourage you to clone this repo and compile the code yourself. I have provided build scripts for Windows and Linux. 

Once you have your program directory set up, the next step is to add it to PATH. At this point, calling `metar` from command line should result in response `No code specified.`.

### Command syntax

```
metar [-notaf] <icao_code1> [<icao_code2>...]
```

By default, the app requests METAR and TAF. When `-notaf` is passed, only METAR data is requested.

## Resources

* [ICAO airport code](https://en.wikipedia.org/wiki/ICAO_airport_code)
* [More on METAR](https://en.wikipedia.org/wiki/METAR)
* [More on TAF](https://aviationweather.gov/taf/decoder)
* [Aviation Weather Center](https://aviationweather.gov)

## License

This software is available under MIT License.
