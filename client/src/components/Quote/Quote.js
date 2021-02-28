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

function Quote(props) {
  const { type, width, ratio, data: initialData } = props;
  const xScaleProvider = discontinuousTimeScaleProviderBuilder().inputDateAccessor(
    (d) => d.date
  );
  const { data, xScale, xAccessor, displayXAccessor } = xScaleProvider(
    initialData
  );
  const max = xAccessor(data[data.length - 1]);
  const min = xAccessor(data[0]);
  const xExtents = [min, max + 2];
  const pricesDisplayFormat = format(".4f");
  const margin = { left: 60, right: 50, top: 30, bottom: 30 };
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
          <Label
            x={(width - margin.left - margin.right) / 2}
            y={-5}
            fontSize="25"
            text={props.pair}
          />
          <Chart id={1} yExtents={(d) => [d.high, d.low]}>
            <XAxis
              showGridLines
              axisAt="bottom"
              orient="bottom"
              ticks={10}
              tickLabelFill="#FFFFFF"
            />
            <YAxis
              showGridLines
              axisAt="left"
              orient="left"
              ticks={10}
              tickFormat={pricesDisplayFormat}
              tickLabelFill="#FFFFFF"
            />
            <CandlestickSeries />
            <EdgeIndicator
              itemType="last"
              orient="right"
              edgeAt="right"
              yAccessor={(d) => d.close}
              fill={(d) => (d.close > d.open ? "#6BA583" : "#FF0000")}
            />
            <MouseCoordinateX
              at="bottom"
              orient="bottom"
              displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
            />
            <MouseCoordinateY displayFormat={pricesDisplayFormat} />
            <ZoomButtons />
          </Chart>
          <CrossHairCursor strokeDasharray="LongDashDot" />
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
