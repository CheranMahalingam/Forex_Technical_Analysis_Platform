import pandas as pd
import numpy as np
from datetime import datetime


class UnemploymentRateIndicator:

    def unemployment_rate_data_selection(self):
        ue_data = pd.read_csv(
            "../../data/external/unemployment/Unemployment_2018-2020.csv")
        ue_data = ue_data.rename(
            columns={"TIME": "Time", "Value": "Unemployment Rate"})
        ue_data = ue_data[{"Time", "Unemployment Rate"}]
        ue_data["Time"] = ue_data["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m").strftime("%Y-%m-%d %H:%M:%S"))

        time_frame = pd.date_range(
            start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
        time_frame = pd.DataFrame(time_frame, columns=["Time"])
        time_frame["Time"] = time_frame["Time"].dt.strftime(
            "%Y-%m-%d %H:%M:%S")

        aud_ue = ue_data[12:47]
        cad_ue = ue_data[153:189]
        jpy_ue = ue_data[668:703]
        gbp_ue = ue_data[1228:1261]
        usd_ue = ue_data[1273:1309]
        eur_ue = ue_data[1638:1673]

        nzd_data = pd.read_csv(
            "../../data/external/unemployment/NZD_Unemployment_2018-2020.csv")
        chf_data = pd.read_csv(
            "../../data/external/unemployment/CHF_Unemployment_2018-2020.csv")
        nzd_data = nzd_data.rename(
            columns={"DATE": "Time", "LRUNTTTTNZQ156S": "Unemployment Rate"})
        chf_data = chf_data.rename(
            columns={"DATE": "Time", "LRUNTTTTCHQ156S": "Unemployment Rate"})
        nzd_data["Time"] = nzd_data["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))
        chf_data["Time"] = chf_data["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))

        nzd_ue = nzd_data[128:139]
        chf_ue = chf_data[75:86]

        create_unemployment_csv(aud_ue, time_frame, aud_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/aud_ue_processed.csv", index=False)
        create_unemployment_csv(cad_ue, time_frame, cad_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/cad_ue_processed.csv", index=False)
        create_unemployment_csv(jpy_ue, time_frame, jpy_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/jpy_ue_processed.csv", index=False)
        create_unemployment_csv(gbp_ue, time_frame, gbp_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/gbp_ue_processed.csv", index=False)
        create_unemployment_csv(usd_ue, time_frame, usd_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/usd_ue_processed.csv", index=False)
        create_unemployment_csv(eur_ue, time_frame, eur_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/eur_ue_processed.csv", index=False)
        create_unemployment_csv(nzd_ue, time_frame, nzd_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/nzd_ue_processed.csv", index=False)
        create_unemployment_csv(chf_ue, time_frame, chf_ue["Unemployment Rate"].iloc[0]).to_csv(
            "../../data/interim/unemployment_rate/chf_ue_processed.csv", index=False)

    def create_unemployment_csv(self, pair, time, initial):
        pair = time.merge(pair, how="left", on="Time")
        pair["Unemployment Rate"].iloc[0] = initial
        pair = pair.fillna(method="ffill")
        return pair
