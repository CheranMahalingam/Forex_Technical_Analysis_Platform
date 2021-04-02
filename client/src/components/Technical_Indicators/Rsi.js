import React from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import { RSISeries } from "@react-financial-charts/series";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import {
  MouseCoordinateY,
  MouseCoordinateX,
  ZoomButtons,
  RSITooltip,
} from "react-financial-charts";

function Rsi(props) {
  const rsiDisplayFormat = format(".5f");

  return (
    <React.Fragment>
      {props.isBottom ? (
        <XAxis
          strokeStyle="#FFFFFF"
          strokeDasharray="Solid"
          axisAt="bottom"
          orient="bottom"
          ticks={15}
          tickLabelFill="#FFFFFF"
          gridLinesStrokeDasharray="Solid"
        />
      ) : null}
      <YAxis
        showGridLines
        strokeStyle="#FFFFFF"
        axisAt="right"
        orient="right"
        ticks={5}
        tickLabelFill="#FFFFFF"
        gridLinesStrokeDasharray="Solid"
      />
      <RSISeries
        yAccessor={props.rsiCalculator.accessor()}
        strokeStyle={props.rsiCalculator.stroke()}
      />
      <RSITooltip
        origin={[10, -5]}
        yAccessor={props.rsiCalculator.accessor()}
        options={props.rsiCalculator.options()}
        textFill="#FFFFFF"
        labelFill="#FFFFFF"
        fontSize={12.5}
        displayFormat={rsiDisplayFormat}
      />
      <MouseCoordinateX
        at="bottom"
        orient="bottom"
        displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
        fill="#303030"
      />
      <MouseCoordinateY displayFormat={rsiDisplayFormat} fill="#303030" />
      {props.isBottom ? <ZoomButtons /> : null}
    </React.Fragment>
  );
}

export default Rsi;
