"""Preprocessing for gdp indicator."""

from datetime import datetime
import pandas as pd


def gdp_preprocess():
    """
    Reads gdp indicator data and creates a csv file with interim gdp data.
    Reformats datetime and strips information outside of 2018-2020.
    The resulting dataframe is stored in the interim folder.
    """
    gdp_data = pd.read_csv("lstm_model/data/external/gdp/GDP_2018-2020.csv")
    gdp_data = gdp_data.rename(columns={"TIME": "Time", "Value": "GDP"})
    gdp_data = gdp_data[{"Time", "GDP"}]

    # Dates come in format year-quarter and are converted to y-m-d H:M:S format
    gdp_data["Time"] = gdp_data["Time"].transform(
        lambda time: datetime.strptime(
            convert_quarter_to_date(time), "%Y-%m-%d"
        ).strftime("%Y-%m-%d %H:%M:%S")
    )

    # Creates a dataframe with dates from 2018-2020 on a 1 minute interval
    time_frame = pd.date_range(
        start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00"
    )
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    # Selects data relevant to time interval 2018-2020
    aud_gdp = gdp_data[16:31]
    gbp_gdp = gdp_data[46:61]
    jpy_gdp = gdp_data[107:122]
    cad_gdp = gdp_data[167:182]
    usd_gdp = gdp_data[258:274]
    chf_gdp = gdp_data[552:567]
    nzd_gdp = gdp_data[567:582]
    eur_gdp = gdp_data[720:736]

    create_gdp_csv(aud_gdp, time_frame, aud_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/aud_gdp_processed.csv", index=False
    )
    create_gdp_csv(gbp_gdp, time_frame, gbp_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/gbp_gdp_processed.csv", index=False
    )
    create_gdp_csv(jpy_gdp, time_frame, jpy_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/jpy_gdp_processed.csv", index=False
    )
    create_gdp_csv(cad_gdp, time_frame, cad_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/cad_gdp_processed.csv", index=False
    )
    create_gdp_csv(usd_gdp, time_frame, usd_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/usd_gdp_processed.csv", index=False
    )
    create_gdp_csv(chf_gdp, time_frame, chf_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/chf_gdp_processed.csv", index=False
    )
    create_gdp_csv(nzd_gdp, time_frame, nzd_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/nzd_gdp_processed.csv", index=False
    )
    create_gdp_csv(eur_gdp, time_frame, eur_gdp["GDP"].iloc[0]).to_csv(
        "lstm_model/data/interim/gdp/eur_gdp_processed.csv", index=False
    )


def create_gdp_csv(pair_df, time_df, initial):
    """
    Fills gdp dataframe missing values with previous values. Originally gdp data changes once
    a month, but dataframe must contain entries for every minute.

    Args:
        pair_df: Currency dataframe containing columns for time and gdp
        time_df: Dataframe specifying time range used for lstm training
        initial: Float representing the gdp of the month or quarter before 2018

    Returns:
        Dataframe with gdp column and time column containing data from 2018-2020 for every minute.
    """
    pair_df = time_df.merge(pair_df, how="left", on="Time")
    pair_df.iloc[0, pair_df.columns.get_loc("GDP")] = initial
    pair_df = pair_df.fillna(method="ffill")
    return pair_df


def convert_quarter_to_date(time):
    """
    Reformats dates given with quarters to month-day format.

    Args:
        time: String representing time in the format year-quarter(e.g. Q1, Q2, ...)

    Returns:
        String representing time in the format year-month-day.
    """
    date = time
    if "Q1" in date:
        date = str(date[:5]) + "01-01"
    elif "Q2" in date:
        date = str(date[:5]) + "04-01"
    elif "Q3" in date:
        date = str(date[:5]) + "07-01"
    else:
        date = str(date[:5]) + "10-01"
    return date
