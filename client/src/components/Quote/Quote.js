import React, { useState, useEffect } from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import { Chart, ChartCanvas } from "@react-financial-charts/core";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import { discontinuousTimeScaleProviderBuilder } from "@react-financial-charts/scales";
import { CandlestickSeries } from "@react-financial-charts/series";
import { withDeviceRatio, withSize } from "@react-financial-charts/utils";
import { fitWidth } from "react-stockcharts/lib/helper";
import {
  CrossHairCursor,
  MouseCoordinateY,
  MouseCoordinateX,
  mouseBasedZoomAnchor,
  ZoomButtons,
  EdgeIndicator,
  Label,
} from "react-financial-charts";
import "./Quote.css";

const candlesAppearance = {
  wickStroke: function fill(d) {
    return d.close > d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  //fill: function fill(d) {
  //return d.close > d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  //},
  fill: "#000000",
  stroke: function outline(d) {
    return d.close > d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  candleStrokeWidth: 0.5,
  width: 3,
};

function Quote(props) {
  const { type, width, ratio, data: initialData } = props;
  const xScaleProvider = discontinuousTimeScaleProviderBuilder().inputDateAccessor(
    (d) => d.timestamp
  );
  const { data, xScale, xAccessor, displayXAccessor } = xScaleProvider(
    initialData
  );
  const max = xAccessor(data[data.length - 1]);
  const min = xAccessor(data[0]);
  const xExtents = [min, max + 2];
  const pricesDisplayFormat = format(".5f");
  const margin = { left: 70, right: 70, top: 20, bottom: 30 };

  return (
    <React.Fragment>
      <div className="dark">
        <ChartCanvas
          height={400}
          ratio={ratio}
          width={width}
          margin={margin}
          type={type}
          data={data}
          seriesName="MSFT"
          xAccessor={xAccessor}
          xScale={xScale}
          xExtents={xExtents}
          displayXAccessor={displayXAccessor}
          zoomAnchor={mouseBasedZoomAnchor}
        >
          {/*<Label
            x={(width - margin.left - margin.right) / 2}
            y={-5}
            fontSize="25"
            text={props.pair}
          />*/}
          <Chart id={1} yExtents={(d) => [d.high, d.low]}>
            <XAxis
              strokeStyle="#FFFFFF"
              strokeDasharray="Solid"
              axisAt="bottom"
              orient="bottom"
              ticks={15}
              tickLabelFill="#FFFFFF"
              gridLinesStrokeDasharray="Solid"
            />
            <YAxis
              showGridLines
              strokeStyle="#FFFFFF"
              axisAt="left"
              orient="left"
              ticks={10}
              tickFormat={pricesDisplayFormat}
              tickLabelFill="#FFFFFF"
              gridLinesStrokeDasharray="Solid"
            />
            <CandlestickSeries {...candlesAppearance} />
            <EdgeIndicator
              itemType="last"
              orient="right"
              edgeAt="right"
              yAccessor={(d) => d.close}
              lineStroke="#FFFFFF"
              displayFormat={pricesDisplayFormat}
              fill={(d) => (d.close > d.open ? "#F72119" : "#39FF14")}
            />
            <MouseCoordinateX
              at="bottom"
              orient="bottom"
              displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
              fill="#303030"
            />
            <MouseCoordinateY
              displayFormat={pricesDisplayFormat}
              fill="#303030"
            />
            <ZoomButtons />
          </Chart>
          <CrossHairCursor strokeDasharray="DashDot" strokeStyle="#4D4DFF" />
        </ChartCanvas>
      </div>
    </React.Fragment>
  );
}

Quote.defaultProps = {
  type: "svg",
};

Quote = fitWidth(Quote);

export default Quote;
