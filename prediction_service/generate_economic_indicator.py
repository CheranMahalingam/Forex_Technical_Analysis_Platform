from fred import Fred
from dotenv import load_dotenv
import os
import series_id_defs
load_dotenv()

fr = Fred(api_key=os.getenv("FRED_API_KEY"), response_type="json")


def generate_interest_rate(start_date, country):
    params = {
        "observation_start": start_date,
    }
    indicator_data = {}
    if country == "Switzerland":
        indicator_data = fr.series.observations(series_id_defs.SWITZERLAND["interest_rate_id"], params=params)
    elif country == "Euro":
        indicator_data = fr.series.observations(series_id_defs.EURO["interest_rate_id"], params=params)
    elif country == "Canada":
        indicator_data = fr.series.observations(series_id_defs.CANADA["interest_rate_id"], params=params)
    elif country == "USA":
        indicator_data = fr.series.observations(series_id_defs.USA["interest_rate_id"], params=params)
    elif country == "Australia":
        indicator_data = fr.series.observations(series_id_defs.AUSTRALIA["interest_rate_id"], params=params)
    elif country == "Japan":
        indicator_data = fr.series.observations(series_id_defs.JAPAN["interest_rate_id"], params=params)
    elif country == "UK":
        indicator_data = fr.series.observations(series_id_defs.UK["interest_rate_id"], params=params)
    elif country == "New Zealand":
        indicator_data = fr.series.observations(series_id_defs.NEW_ZEALAND["interest_rate_id"], params=params)
    else:
        raise Exception("No interest rate data for", country)

    interest_rate = []
    new_observations = indicator_data["observations"]
    for i in range(len(new_observations)):
        interest_rate.append({"date": new_observations[i]["date"], "value": new_observations[i]["value"]})
    return interest_rate


def generate_cpi(start_date, country):
    params = {
        "observation_start": start_date,
    }
    indicator_data = {}
    if country == "Switzerland":
        indicator_data = fr.series.observations(series_id_defs.SWITZERLAND["cpi_id"], params=params)
    elif country == "Euro":
        indicator_data = fr.series.observations(series_id_defs.EURO["cpi_id"], params=params)
    elif country == "Canada":
        indicator_data = fr.series.observations(series_id_defs.CANADA["cpi_id"], params=params)
    elif country == "USA":
        indicator_data = fr.series.observations(series_id_defs.USA["cpi_id"], params=params)
    elif country == "Australia":
        indicator_data = fr.series.observations(series_id_defs.AUSTRALIA["cpi_id"], params=params)
    elif country == "Japan":
        indicator_data = fr.series.observations(series_id_defs.JAPAN["cpi_id"], params=params)
    elif country == "UK":
        indicator_data = fr.series.observations(series_id_defs.UK["cpi_id"], params=params)
    elif country == "New Zealand":
        indicator_data = fr.series.observations(series_id_defs.NEW_ZEALAND["cpi_id"], params=params)
    else:
        raise Exception("No cpi data for", country)

    cpi = []
    new_observations = indicator_data["observations"]
    for i in range(len(new_observations)):
        cpi.append({"date": new_observations[i]["date"], "value": new_observations[i]["value"]})
    return cpi


def generate_gdp(start_date, country):
    params = {
        "observation_start": start_date,
    }
    indicator_data = {}
    if country == "Switzerland":
        indicator_data = fr.series.observations(series_id_defs.SWITZERLAND["gdp_id"], params=params)
    elif country == "Euro":
        indicator_data = fr.series.observations(series_id_defs.EURO["gdp_id"], params=params)
    elif country == "Canada":
        indicator_data = fr.series.observations(series_id_defs.CANADA["gdp_id"], params=params)
    elif country == "USA":
        indicator_data = fr.series.observations(series_id_defs.USA["gdp_id"], params=params)
    elif country == "Australia":
        indicator_data = fr.series.observations(series_id_defs.AUSTRALIA["gdp_id"], params=params)
    elif country == "Japan":
        indicator_data = fr.series.observations(series_id_defs.JAPAN["gdp_id"], params=params)
    elif country == "UK":
        indicator_data = fr.series.observations(series_id_defs.UK["gdp_id"], params=params)
    elif country == "New Zealand":
        indicator_data = fr.series.observations(series_id_defs.NEW_ZEALAND["gdp_id"], params=params)
    else:
        raise Exception("No gdp data for", country)

    gdp = []
    new_observations = indicator_data["observations"]
    for i in range(len(new_observations)):
        gdp.append({"date": new_observations[i]["date"], "value": new_observations[i]["value"]})
    return gdp
