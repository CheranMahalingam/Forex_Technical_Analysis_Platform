"""Preprocessing for cpi indicator."""

from datetime import datetime
import pandas as pd


def cpi_preprocess():
    """
    Reads cpi indicator data and creates a csv file containing processed cpi data.
    Reformats datetime and strips information outside of 2018-2020.
    The resulting dataframe is stored in the interim folder.
    """
    cpi_data = pd.read_csv("lstm_model/data/external/cpi/CPI_2018-2020.csv")
    cpi_data = cpi_data.rename(columns={"TIME": "Time", "Value": "CPI"})
    cpi_data = cpi_data[{"Time", "CPI"}]
    cpi_data["Time"] = cpi_data["Time"].transform(
        lambda time: datetime.strptime(time, "%Y-%m").strftime("%Y-%m-%d %H:%M:%S"))

    # Creates a dataframe with dates from 2018-2020 on a 1 minute interval
    time_frame = pd.date_range(start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    # Selects data relevant to time interval 2018-2020
    cad_cpi = cpi_data[1:37]
    jpy_cpi = cpi_data[38:74]
    chf_cpi = cpi_data[75:111]
    gbp_cpi = cpi_data[112:148]
    usd_cpi = cpi_data[149:185]
    eur_cpi = cpi_data[223:259]

    # Similar steps must be applied to AUD and NZD
    # AUD, NZD comes from different data source with different csv structure
    australia_cpi_data = pd.read_csv("lstm_model/data/external/cpi/AUD_CPI_2018-2020.csv")
    new_zealand_cpi_data = pd.read_csv("lstm_model/data/external/cpi/NZD_CPI_2018-2020.csv")
    aud_cpi = australia_cpi_data[25:37]
    nzd_cpi = new_zealand_cpi_data[4:]

    nzd_cpi = nzd_cpi.rename(columns={"Year ended": "Time", "Percentage change": "CPI"})
    aud_cpi = aud_cpi.rename(columns={"Quarter": "Time", "Change from previous quarter (%)": "CPI"})
    aud_cpi["Time"] = aud_cpi["Time"].transform(
        lambda time: datetime.strptime(
            time[:4] + "20" + time[4:], "%b-%Y").strftime("%Y-%m-%d %H:%M:%S"))
    nzd_cpi["Time"] = nzd_cpi["Time"].transform(
        lambda time: datetime.strptime(
            time[:4] + "20" + time[4:], "%b-%Y").strftime("%Y-%m-%d %H:%M:%S"))

    # Convert string cpi values to floats
    aud_cpi = create_cpi_csv(aud_cpi, time_frame, aud_cpi["CPI"].iloc[0])
    aud_cpi['CPI'] = pd.to_numeric(aud_cpi['CPI'])

    aud_cpi.to_csv("lstm_model/data/interim/cpi/aud_cpi_processed.csv", index=False)
    create_cpi_csv(
        cad_cpi, time_frame, cad_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/cad_cpi_processed.csv", index=False)
    create_cpi_csv(
        jpy_cpi, time_frame, jpy_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/jpy_cpi_processed.csv", index=False)
    create_cpi_csv(
        chf_cpi, time_frame, chf_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/chf_cpi_processed.csv", index=False)
    create_cpi_csv(
        gbp_cpi, time_frame, gbp_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/gbp_cpi_processed.csv", index=False)
    create_cpi_csv(
        usd_cpi, time_frame, usd_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/usd_cpi_processed.csv", index=False)
    create_cpi_csv(
        eur_cpi, time_frame, eur_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/eur_cpi_processed.csv", index=False)
    create_cpi_csv(
        nzd_cpi, time_frame, nzd_cpi["CPI"].iloc[0]).to_csv(
            "lstm_model/data/interim/cpi/nzd_cpi_processed.csv", index=False)

def create_cpi_csv(pair_df, time_df, initial):
    """
    Fills cpi dataframe missing values with previous values. Originally cpi data changes once
    a month, but dataframe must contain entries for every minute.

    Args:
        pair_df: Currency dataframe containing columns for time and cpi
        time_df: Dataframe specifying time range used for lstm training
        initial: Float representing the cpi of the month or quarter before 2018

    Returns:
        Dataframe with cpi column and time column containing data from 2018-2020 for every minute
    """
    pair_df = time_df.merge(pair_df, how="left", on="Time")
    pair_df.iloc[0, pair_df.columns.get_loc("CPI")] = initial
    pair_df = pair_df.fillna(method="ffill")
    return pair_df
