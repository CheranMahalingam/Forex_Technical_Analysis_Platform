import React from "react";
import { format } from "d3-format";
import { LineSeries } from "@react-financial-charts/series";
import { MovingAverageTooltip } from "react-financial-charts";

function Ema(props) {
  const emaDisplayFormat = format(".5f");
  return (
    <React.Fragment>
      {props.ema.map((ema) => {
        return (
          <LineSeries yAccessor={ema.accessor()} strokeStyle={ema.stroke()} />
        );
      })}
      <MovingAverageTooltip
        origin={[10, -45]}
        textFill="#FFFFFF"
        labelFill="#FFFFFF"
        fontSize={12.5}
        displayFormat={emaDisplayFormat}
        options={props.ema.map((ema) => {
          return {
            stroke: ema.stroke(),
            type: "EMA",
            windowSize: ema.options().windowSize,
            yAccessor: ema.accessor(),
          };
        })}
      ></MovingAverageTooltip>
    </React.Fragment>
  );
}

export default Ema;
