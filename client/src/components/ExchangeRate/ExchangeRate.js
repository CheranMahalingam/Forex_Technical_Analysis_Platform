import React from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import { CandlestickSeries } from "@react-financial-charts/series";
import {
  MouseCoordinateY,
  MouseCoordinateX,
  ZoomButtons,
  EdgeIndicator,
  OHLCTooltip,
} from "react-financial-charts";
import BollingerBand from "../Technical_Indicators/Bollinger_Band";
import Ema from "../Technical_Indicators/Ema";

function ExchangeRate(props) {
  const pricesDisplayFormat = format(".5f");
  const changeDisplayFormat = format("+.5f");

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
        ticks={10}
        tickFormat={pricesDisplayFormat}
        tickLabelFill="#FFFFFF"
        gridLinesStrokeDasharray="Solid"
      />
      {props.selectedIndicators.includes("Bollinger Bands") ? (
        <BollingerBand bollingerCalculator={props.bollingerCalculator} />
      ) : null}
      {props.selectedIndicators.includes("Exponential Moving Average") ? (
        <Ema ema={props.ema} />
      ) : null}
      <CandlestickSeries {...candlesAppearance} />
      <OHLCTooltip
        origin={[props.moveRight, 0]}
        labelFill="#FFFFFF"
        textFill={(d) => props.handleColour(d)}
        fontSize={12.5}
        ohlcFormat={pricesDisplayFormat}
        changeFormat={changeDisplayFormat}
      />
      <EdgeIndicator
        itemType="last"
        orient="right"
        edgeAt="right"
        yAccessor={(d) => d.close}
        lineStroke="#FFFFFF"
        displayFormat={pricesDisplayFormat}
        fill={(d) => (d.close < d.open ? "#F70D1A" : "#52D017")}
        fontSize={13.5}
        fitToText={true}
      />
      <MouseCoordinateX
        at="bottom"
        orient="bottom"
        displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
        fill="#303030"
      />
      <MouseCoordinateY displayFormat={pricesDisplayFormat} fill="#303030" />
      {props.isBottom ? <ZoomButtons /> : null}
    </React.Fragment>
  );
}

const candlesAppearance = {
  wickStroke: function fill(d) {
    return d.close < d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  fill: "#000000",
  stroke: function outline(d) {
    return d.close < d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  candleStrokeWidth: 1,
  width: 3,
};

export default ExchangeRate;
