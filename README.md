# IBGE Rio de Janeiro County Code Scraper

This script scrapes the [IBGE](https://ibge.gov.br/explica/codigos-dos-municipios.php) to get the country code of Rio de Janeiro State

## CSS Selector

To get the Rio de Janeiro table we used the selector below

```
body > section > article > div.container-codigos > table:nth-child(21) > tbody
```


