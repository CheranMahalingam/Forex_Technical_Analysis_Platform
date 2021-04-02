import React from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import { StochasticSeries } from "@react-financial-charts/series";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import {
  MouseCoordinateY,
  MouseCoordinateX,
  ZoomButtons,
  StochasticTooltip,
} from "react-financial-charts";

function Stochastic(props) {
  const stochasticDisplayFormat = format(".5f");

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
      <StochasticSeries yAccessor={props.stoCalculator.accessor()} />
      <StochasticTooltip
        origin={[10, -5]}
        yAccessor={props.stoCalculator.accessor()}
        options={props.stoCalculator.options()}
        label="Fast STO"
        textFill="#FFFFFF"
        labelFill="#FFFFFF"
        fontSize={12.5}
        displayFormat={stochasticDisplayFormat}
        appearance={{ stroke: StochasticSeries.defaultProps.strokeStyle }}
      />
      <MouseCoordinateX
        at="bottom"
        orient="bottom"
        displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
        fill="#303030"
      />
      <MouseCoordinateY
        displayFormat={stochasticDisplayFormat}
        fill="#303030"
      />
      {props.isBottom ? <ZoomButtons /> : null}
    </React.Fragment>
  );
}

export default Stochastic;
