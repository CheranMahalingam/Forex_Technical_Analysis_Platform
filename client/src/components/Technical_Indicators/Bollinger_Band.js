import React from "react";
import { format } from "d3-format";
import { BollingerSeries } from "@react-financial-charts/series";
import { BollingerBandTooltip } from "react-financial-charts";

function BollingerBand(props) {
  const bbDisplayFormat = format(".5f");
  return (
    <React.Fragment>
      <BollingerSeries
        options={props.bollingerCalculator.options()}
        strokeStyle={{
          top: "#4EE2EC",
          middle: "#4EE2EC",
          bottom: "#4EE2EC",
        }}
        fillStyle="#151B54"
      />
      <BollingerBandTooltip
        options={props.bollingerCalculator.options()}
        origin={[10, 40]}
        textFill="#FFFFFF"
        labelFill="#FFFFFF"
        fontSize={12.5}
        displayFormat={bbDisplayFormat}
      />
    </React.Fragment>
  );
}

export default BollingerBand;
