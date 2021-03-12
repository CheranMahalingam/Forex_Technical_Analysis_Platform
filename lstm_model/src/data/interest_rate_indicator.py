"""Preprocessing for interest rate indicator."""

from datetime import datetime
import pandas as pd


def interest_rate_preprocess():
    """
    Reads interest rate indicator data and creates a csv file with interim interest rate
    data. Reformats datetime and strips information outside of 2018-2020. The
    resulting dataframe is stored in the interim folder.
    """
    ir_data = pd.read_csv("../data/external/interest_rates/Interest_rates_2018-2020.csv")
    ir_data = ir_data.rename(columns={"TIME": "Time", "Value": "Interest Rate"})
    ir_data = ir_data[{"Time", "Interest Rate"}]
    ir_data["Time"] = ir_data["Time"].transform(
        lambda time: datetime.strptime(time, "%Y-%m").strftime("%Y-%m-%d %H:%M:%S"))

    # Selects data relevant to time interval 2018-2020
    aud_ir = ir_data[13:49]
    eur_ir = ir_data[356:392]
    jpy_ir = ir_data[307:343]
    gbp_ir = ir_data[258:294]
    nzd_ir = ir_data[209:245]
    usd_ir = ir_data[160:196]
    chf_ir = ir_data[62:98]
    cad_ir = ir_data[111:147]

    # Creates a dataframe with dates from 2018-2020 on a 1 minute interval
    time_frame = pd.date_range(start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    create_interest_rate_csv(
        aud_ir, time_frame, aud_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/aud_ir_processed.csv", index=False)
    create_interest_rate_csv(
        cad_ir, time_frame, cad_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/cad_ir_processed.csv", index=False)
    create_interest_rate_csv(
        eur_ir, time_frame, eur_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/eur_ir_processed.csv", index=False)
    create_interest_rate_csv(
        jpy_ir, time_frame, jpy_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/jpy_ir_processed.csv", index=False)
    create_interest_rate_csv(
        gbp_ir, time_frame, gbp_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/gbp_ir_processed.csv", index=False)
    create_interest_rate_csv(
        nzd_ir, time_frame, nzd_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/nzd_ir_processed.csv", index=False)
    create_interest_rate_csv(
        usd_ir, time_frame, usd_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/usd_ir_processed.csv", index=False)
    create_interest_rate_csv(
        chf_ir, time_frame, chf_ir["Interest Rate"].iloc[0]).to_csv(
            "../data/interim/interest_rate/chf_ir_processed.csv", index=False)

def create_interest_rate_csv(pair, time, initial):
    """
    Fills interest rate dataframe missing values with previous values. Originally interest rate
    data changes once a month, but dataframe must contain entries for every minute.

    Args:
        pair_df: Currency dataframe containing columns for time and interest rate
        time_df: Dataframe specifying time range used for lstm training
        initial: Float representing the interest rate of the month or quarter before 2018

    Returns:
        Dataframe with interest rate column and time column containing data from 2018-2020 for
        every minute.
    """
    pair = time.merge(pair, how="left", on="Time")
    pair.iloc[0, pair.columns.get_loc("Interest Rate")] = initial
    pair = pair.fillna(method="ffill")
    return pair
