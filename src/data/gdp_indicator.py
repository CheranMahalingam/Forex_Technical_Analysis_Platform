import pandas as pd
import numpy as np
from datetime import datetime


class GdpIndicator:

    def gdp_data_selection(self):
        gdp_data = pd.read_csv("../../data/external/gdp/GDP_2018-2020.csv")
        gdp_data = gdp_data.rename(columns={"TIME": "Time", "Value": "GDP"})
        gdp_data = gdp_data[{"Time", "GDP"}]
        gdp_data["Time"] = gdp_data["Time"].transform(
            lambda time: convert_quarter_to_date(time))
        gdp_data["Time"] = gdp_data["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))

        time_frame = pd.date_range(
            start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
        time_frame = pd.DataFrame(time_frame, columns=["Time"])
        time_frame["Time"] = time_frame["Time"].dt.strftime(
            "%Y-%m-%d %H:%M:%S")

        aud_gdp = gdp_data[16:31]
        gbp_gdp = gdp_data[46:61]
        jpy_gdp = gdp_data[107:122]
        cad_gdp = gdp_data[167:182]
        usd_gdp = gdp_data[258:274]
        chf_gdp = gdp_data[552:567]
        nzd_gdp = gdp_data[567:582]
        eur_gdp = gdp_data[720:736]

        create_gdp_csv(aud_gdp, time_frame, aud_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/aud_gdp_processed.csv", index=False)
        create_gdp_csv(gbp_gdp, time_frame, gbp_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/gbp_gdp_processed.csv", index=False)
        create_gdp_csv(jpy_gdp, time_frame, jpy_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/jpy_gdp_processed.csv", index=False)
        create_gdp_csv(cad_gdp, time_frame, cad_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/cad_gdp_processed.csv", index=False)
        create_gdp_csv(usd_gdp, time_frame, usd_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/usd_gdp_processed.csv", index=False)
        create_gdp_csv(chf_gdp, time_frame, chf_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/chf_gdp_processed.csv", index=False)
        create_gdp_csv(nzd_gdp, time_frame, nzd_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/nzd_gdp_processed.csv", index=False)
        create_gdp_csv(eur_gdp, time_frame, eur_gdp["GDP"].iloc[0]).to_csv(
            "../../data/interim/gdp/eur_gdp_processed.csv", index=False)

    def convert_quarter_to_date(self, date):
        if "Q1" in date:
            date = str(date[:5]) + "01-01"
        elif "Q2" in date:
            date = str(date[:5]) + "04-01"
        elif "Q3" in date:
            date = str(date[:5]) + "07-01"
        else:
            date = str(date[:5]) + "10-01"
        return date

    def create_gdp_csv(self, pair, time, initial):
        pair = time.merge(pair, how="left", on="Time")
        pair["GDP"].iloc[0] = initial
        pair = pair.fillna(method="ffill")
        return pair
