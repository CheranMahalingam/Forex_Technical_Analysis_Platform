import pandas as pd
import numpy as np
from datetime import datetime


class PpiIndicator:

    def ppi_data_selection(self):
        ppi_data = pd.read_csv("../../data/external/ppi/PPI_2018-2020.csv")
        ppi_data = ppi_data.rename(columns={"TIME": "Time", "Value": "PPI"})
        ppi_data = ppi_data[{"Time", "PPI"}]
        ppi_data["Time"] = ppi_data["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m").strftime("%Y-%m-%d %H:%M:%S"))

        jpy_ppi = ppi_data[576:611]
        chf_ppi = ppi_data[1093:1128]
        gbp_ppi = ppi_data[1187:1222]
        eur_ppi = ppi_data[1793:1828]

        time_frame = pd.date_range(
            start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
        time_frame = pd.DataFrame(time_frame, columns=["Time"])
        time_frame["Time"] = time_frame["Time"].dt.strftime(
            "%Y-%m-%d %H:%M:%S")

        cad_ppi = pd.read_csv("../../data/external/ppi/CAD_PPI_2018-2020.csv")
        usd_ppi = pd.read_csv("../../data/external/ppi/USD_PPI_2018-2020.csv")
        aud_ppi = pd.read_csv("../../data/external/ppi/AUD_PPI_2018-2020.csv")
        nzd_ppi = pd.read_csv("../../data/external/ppi/NZD_PPI_2018-2020.csv")

        cad_ppi = cad_ppi[232:243]
        usd_ppi = usd_ppi[232:243]
        aud_ppi = aud_ppi[198:209]
        nzd_ppi = nzd_ppi[232:243]

        cad_ppi = cad_ppi.rename(
            columns={"DATE": "Time", "PIEAMP01CAQ661N": "PPI"})
        usd_ppi = usd_ppi.rename(
            columns={"DATE": "Time", "PIEAMP01USQ661N": "PPI"})
        aud_ppi = aud_ppi.rename(
            columns={"DATE": "Time", "PIEAMP01AUQ661N": "PPI"})
        nzd_ppi = nzd_ppi.rename(
            columns={"DATE": "Time", "PIEAMP01NZQ661N": "PPI"})
        cad_ppi["Time"] = cad_ppi["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))
        usd_ppi["Time"] = usd_ppi["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))
        aud_ppi["Time"] = aud_ppi["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))
        nzd_ppi["Time"] = nzd_ppi["Time"].transform(
            lambda time: datetime.strptime(time, "%Y-%m-%d").strftime("%Y-%m-%d %H:%M:%S"))

        create_ppi_csv(jpy_ppi, time_frame, jpy_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/jpy_ppi_processed.csv", index=False)
        create_ppi_csv(chf_ppi, time_frame, chf_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/chf_ppi_processed.csv", index=False)
        create_ppi_csv(gbp_ppi, time_frame, gbp_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/gbp_ppi_processed.csv", index=False)
        create_ppi_csv(eur_ppi, time_frame, eur_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/eur_ppi_processed.csv", index=False)
        create_ppi_csv(cad_ppi, time_frame, cad_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/cad_ppi_processed.csv", index=False)
        create_ppi_csv(usd_ppi, time_frame, usd_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/usd_ppi_processed.csv", index=False)
        create_ppi_csv(aud_ppi, time_frame, aud_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/aud_ppi_processed.csv", index=False)
        create_ppi_csv(nzd_ppi, time_frame, nzd_ppi["PPI"].iloc[0]).to_csv(
            "../../data/interim/ppi/nzd_ppi_processed.csv", index=False)

    def create_ppi_csv(self, pair, time, initial):
        pair = time.merge(pair, how="left", on="Time")
        pair["PPI"].iloc[0] = initial
        pair = pair.fillna(method="ffill")
        return pair
