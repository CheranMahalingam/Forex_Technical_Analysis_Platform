import React from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import { MACDSeries } from "@react-financial-charts/series";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import {
  MouseCoordinateY,
  MouseCoordinateX,
  ZoomButtons,
  MACDTooltip,
} from "react-financial-charts";

function Macd(props) {
  const macdDisplayFormat = format(".5f");

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
        gridLinesStrokeWidth={0.5}
        strokeStyle="#FFFFFF"
        axisAt="right"
        orient="right"
        ticks={5}
        tickFormat={macdDisplayFormat}
        tickLabelFill="#FFFFFF"
        gridLinesStrokeDasharray="Solid"
      />
      <MACDSeries
        yAccessor={props.macdCalculator.accessor()}
        options={props.macdCalculator.options()}
        {...props.macdAppearance}
      />
      <MACDTooltip
        origin={[10, 10]}
        yAccessor={props.macdCalculator.accessor()}
        options={props.macdCalculator.options()}
        appearance={props.macdAppearance}
        textFill="#FFFFFF"
        labelFill="#FFFFFF"
        fontSize={12.5}
        displayFormat={macdDisplayFormat}
      />
      <MouseCoordinateX
        at="bottom"
        orient="bottom"
        displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
        fill="#303030"
      />
      <MouseCoordinateY displayFormat={macdDisplayFormat} fill="#303030" />
      {props.isBottom ? <ZoomButtons /> : null}
    </React.Fragment>
  );
}

export default Macd;
